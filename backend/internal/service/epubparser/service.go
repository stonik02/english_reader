package epubparser

import (
	"archive/zip"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"path"
	"regexp"
	"strings"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

const (
	maxFiles        = 2_000
	maxUnpackedSize = 200 << 20
	maxChapterSize  = 5 << 20
	maxCoverSize    = 5 << 20
	minimumEPUBMIME = "application/epub+zip"
)

var (
	titleTag  = regexp.MustCompile(`(?is)<dc:title[^>]*>(.*?)</dc:title>`)
	authorTag = regexp.MustCompile(`(?is)<dc:creator[^>]*>(.*?)</dc:creator>`)
	scriptTag = regexp.MustCompile(`(?is)<script\b.*?</script\s*>|<iframe\b.*?</iframe\s*>|<form\b.*?</form\s*>`)
	eventAttr = regexp.MustCompile(`(?is)\son[a-z]+\s*=\s*("[^"]*"|'[^']*')`)
	imgTag    = regexp.MustCompile(`(?is)<img\b[^>]*>`)
	srcAttr   = regexp.MustCompile(`(?is)\bsrc\s*=\s*["']([^"']+)["']`)
	tags      = regexp.MustCompile(`(?is)<[^>]+>`)
	space     = regexp.MustCompile(`\s+`)
)

type Service struct {
}

type Result struct {
	Title    string
	Author   string
	Chapters []domain.Chapter
	Cover    *Cover
}

type Cover struct {
	Data        []byte
	ContentType string
}

func New() *Service {
	return &Service{}
}

func (s *Service) Parse(sourcePath, bookID string) (Result, error) {
	archive, err := zip.OpenReader(sourcePath)
	if err != nil {
		return Result{}, fmt.Errorf("open EPUB: %w", err)
	}
	defer archive.Close()
	if len(archive.File) == 0 || len(archive.File) > maxFiles {
		return Result{}, domain.ErrInvalidUpload
	}

	files := make(map[string]*zip.File, len(archive.File))
	var total uint64
	for _, file := range archive.File {
		if path.IsAbs(file.Name) || strings.Contains(path.Clean(file.Name), "../") {
			return Result{}, domain.ErrInvalidUpload
		}
		total += file.UncompressedSize64
		if total > maxUnpackedSize {
			return Result{}, domain.ErrInvalidUpload
		}
		files[file.Name] = file
	}
	mimetype, err := readFile(files, "mimetype", 128)
	if err != nil || strings.TrimSpace(string(mimetype)) != minimumEPUBMIME {
		return Result{}, domain.ErrInvalidUpload
	}
	container, err := readFile(files, "META-INF/container.xml", 1<<20)
	if err != nil {
		return Result{}, domain.ErrInvalidUpload
	}
	rootfile := regexp.MustCompile(`(?is)full-path=["']([^"']+)["']`).FindStringSubmatch(string(container))
	if len(rootfile) != 2 {
		return Result{}, domain.ErrInvalidUpload
	}
	opfPath := path.Clean(rootfile[1])
	opf, err := readFile(files, opfPath, 2<<20)
	if err != nil {
		return Result{}, domain.ErrInvalidUpload
	}
	result := Result{Title: textOf(titleTag, string(opf)), Author: textOf(authorTag, string(opf))}
	if result.Title == "" {
		result.Title = "Untitled EPUB"
	}
	manifest, spine, err := parsePackage(opf)
	if err != nil {
		return Result{}, domain.ErrInvalidUpload
	}
	for sequence, itemID := range spine {
		href, ok := manifest[itemID]
		if !ok {
			return Result{}, domain.ErrInvalidUpload
		}
		href = path.Join(path.Dir(opfPath), href)
		content, err := readFile(files, href, maxChapterSize)
		if err != nil {
			return Result{}, domain.ErrInvalidUpload
		}
		sanitized := embedImages(sanitize(string(content)), files, path.Dir(href))
		chapterText := plainText(sanitized)
		if chapterText == "" {
			continue
		}
		result.Chapters = append(result.Chapters, domain.Chapter{
			BookID:        bookID,
			Sequence:      len(result.Chapters),
			Href:          href,
			StartCFI:      fmt.Sprintf("epubcfi(/6/%d)", sequence+2),
			EndCFI:        fmt.Sprintf("epubcfi(/6/%d)", sequence+4),
			SanitizedHTML: sanitized,
			PlainText:     chapterText,
		})
	}
	if len(result.Chapters) == 0 {
		return Result{}, domain.ErrInvalidUpload
	}
	cover, err := parseCover(opf, files, opfPath)
	if err != nil {
		return Result{}, domain.ErrInvalidUpload
	}
	result.Cover = cover
	return result, nil
}

func parseCover(opf []byte, files map[string]*zip.File, opfPath string) (*Cover, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(opf)))
	type item struct{ href, mediaType, properties string }
	manifest := make(map[string]item)
	var coverID, coverProperty string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}
		switch start.Name.Local {
		case "item":
			id := attribute(start.Attr, "id")
			if id != "" {
				manifest[id] = item{attribute(start.Attr, "href"), attribute(start.Attr, "media-type"), attribute(start.Attr, "properties")}
			}
		case "meta":
			if attribute(start.Attr, "name") == "cover" {
				coverID = attribute(start.Attr, "content")
			}
			if attribute(start.Attr, "property") == "cover-image" {
				var value string
				if err := decoder.DecodeElement(&value, &start); err != nil {
					return nil, err
				}
				coverProperty = strings.TrimSpace(value)
			}
		}
	}
	if coverID == "" {
		coverID = coverProperty
	}
	if coverID == "" {
		for id, candidate := range manifest {
			for _, property := range strings.Fields(candidate.properties) {
				if property == "cover-image" {
					coverID = id
					break
				}
			}
			if coverID != "" {
				break
			}
		}
	}
	if coverID == "" {
		return nil, nil
	}
	candidate, ok := manifest[coverID]
	if !ok || !allowedCoverType(candidate.mediaType) {
		return nil, domain.ErrInvalidUpload
	}
	data, err := readFile(files, path.Join(path.Dir(opfPath), candidate.href), maxCoverSize)
	if err != nil || !validCoverMagic(candidate.mediaType, data) {
		return nil, domain.ErrInvalidUpload
	}
	return &Cover{Data: data, ContentType: candidate.mediaType}, nil
}

func allowedCoverType(contentType string) bool {
	return contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif" || contentType == "image/webp"
}

func validCoverMagic(contentType string, data []byte) bool {
	switch contentType {
	case "image/jpeg":
		return len(data) >= 3 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff
	case "image/png":
		return len(data) >= 8 && string(data[:8]) == "\x89PNG\r\n\x1a\n"
	case "image/gif":
		return len(data) >= 6 && (string(data[:6]) == "GIF87a" || string(data[:6]) == "GIF89a")
	case "image/webp":
		return len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP"
	default:
		return false
	}
}

func parsePackage(opf []byte) (map[string]string, []string, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(opf)))
	manifest := make(map[string]string)
	spine := make([]string, 0)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return manifest, spine, nil
		}
		if err != nil {
			return nil, nil, err
		}

		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		switch start.Name.Local {
		case "item":
			id, href := attribute(start.Attr, "id"), attribute(start.Attr, "href")
			if id != "" && href != "" {
				manifest[id] = href
			}
		case "itemref":
			if idref := attribute(start.Attr, "idref"); idref != "" {
				spine = append(spine, idref)
			}
		}
	}
}

func attribute(attributes []xml.Attr, name string) string {
	for _, attribute := range attributes {
		if attribute.Name.Local == name {
			return attribute.Value
		}
	}
	return ""
}

func readFile(files map[string]*zip.File, name string, limit int64) ([]byte, error) {
	file, ok := files[name]
	if !ok || int64(file.UncompressedSize64) > limit {
		return nil, io.ErrUnexpectedEOF
	}
	reader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(io.LimitReader(reader, limit+1))
}

func sanitize(value string) string {
	return eventAttr.ReplaceAllString(scriptTag.ReplaceAllString(value, ""), "")
}

func embedImages(value string, files map[string]*zip.File, chapterDir string) string {
	return imgTag.ReplaceAllStringFunc(value, func(tag string) string {
		match := srcAttr.FindStringSubmatch(tag)
		if len(match) != 2 || strings.HasPrefix(strings.ToLower(match[1]), "data:") {
			return tag
		}
		source := strings.SplitN(match[1], "#", 2)[0]
		assetPath := path.Clean(path.Join(chapterDir, source))
		if path.IsAbs(source) || strings.HasPrefix(assetPath, "../") {
			return tag
		}
		contentType := imageContentType(assetPath)
		if contentType == "" {
			return tag
		}
		data, err := readFile(files, assetPath, maxCoverSize)
		if err != nil || !validCoverMagic(contentType, data) {
			return tag
		}
		encoded := "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(data)
		return strings.Replace(tag, match[0], `src="`+encoded+`"`, 1)
	})
}

func imageContentType(name string) string {
	return map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
	}[strings.ToLower(path.Ext(name))]
}

func plainText(value string) string {
	return strings.TrimSpace(space.ReplaceAllString(html.UnescapeString(tags.ReplaceAllString(value, " ")), " "))
}

func textOf(pattern *regexp.Regexp, value string) string {
	match := pattern.FindStringSubmatch(value)
	if len(match) != 2 {
		return ""
	}
	return plainText(match[1])
}
