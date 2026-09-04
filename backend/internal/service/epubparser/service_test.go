package epubparser

import (
	"archive/zip"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestPlainTextRemovesUnsafeMarkup(t *testing.T) {
	value := sanitize(`<p onclick="evil()">Hello <script>alert(1)</script>world</p>`)
	if plainText(value) != "Hello world" {
		t.Fatalf("plainText() = %q", plainText(value))
	}
}

func TestParsePackageAcceptsManifestAttributesInAnyOrder(t *testing.T) {
	opf := []byte(`<?xml version="1.0"?>
<package>
  <manifest>
    <item href="chapter-1.xhtml" id="chapter-1" media-type="application/xhtml+xml" />
  </manifest>
  <spine><itemref linear="yes" idref="chapter-1" /></spine>
</package>`)

	manifest, spine, err := parsePackage(opf)
	if err != nil {
		t.Fatalf("parsePackage() error = %v", err)
	}
	if want := map[string]string{"chapter-1": "chapter-1.xhtml"}; !reflect.DeepEqual(manifest, want) {
		t.Fatalf("manifest = %#v, want %#v", manifest, want)
	}
	if want := []string{"chapter-1"}; !reflect.DeepEqual(spine, want) {
		t.Fatalf("spine = %#v, want %#v", spine, want)
	}
}

func TestCoverValidationAcceptsOnlyDeclaredSafeRasterFormats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		contentType string
		data        []byte
		want        bool
	}{
		{"image/jpeg", []byte{0xff, 0xd8, 0xff, 0xe0}, true},
		{"image/png", []byte("\x89PNG\r\n\x1a\n"), true},
		{"image/svg+xml", []byte("<svg/>"), false},
		{"image/png", []byte("not an image"), false},
	}
	for _, test := range tests {
		if got := allowedCoverType(test.contentType) && validCoverMagic(test.contentType, test.data); got != test.want {
			t.Errorf("cover validation for %s = %t, want %t", test.contentType, got, test.want)
		}
	}
}

func TestParseExtractsEPUB3CoverAndLeavesMissingCoverEmpty(t *testing.T) {
	for _, test := range []struct {
		name      string
		withCover bool
	}{
		{name: "with cover", withCover: true},
		{name: "without cover", withCover: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := writeTestEPUB(t, test.withCover)
			result, err := New().Parse(path, "book-1")
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if (result.Cover != nil) != test.withCover {
				t.Fatalf("Parse() cover = %#v, want present=%t", result.Cover, test.withCover)
			}
			if result.Cover != nil && result.Cover.ContentType != "image/png" {
				t.Fatalf("cover content type = %q", result.Cover.ContentType)
			}
			if !strings.Contains(result.Chapters[0].SanitizedHTML, "data:image/png;base64,") {
				t.Fatalf("chapter image was not embedded: %q", result.Chapters[0].SanitizedHTML)
			}
			if len(result.Chapters) != 1 {
				t.Fatalf("chapter count = %d, want 1 readable chapter", len(result.Chapters))
			}
		})
	}
}

func TestParseExtractsEPUB2MetaCover(t *testing.T) {
	path := writeTestEPUBWithOPF(t, `<package><metadata><dc:title xmlns:dc="x">Book</dc:title><meta name="cover" content="cover"/></metadata><manifest><item id="chapter" href="chapter.xhtml" media-type="application/xhtml+xml"/><item id="cover" href="cover.png" media-type="image/png"/></manifest><spine><itemref idref="chapter"/></spine></package>`)

	result, err := New().Parse(path, "book-1")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if result.Cover == nil || result.Cover.ContentType != "image/png" {
		t.Fatalf("Parse() cover = %#v, want PNG cover", result.Cover)
	}
}

func writeTestEPUB(t *testing.T, withCover bool) string {
	return writeTestEPUBWithOPF(t, `<package><metadata><dc:title xmlns:dc="x">Book</dc:title></metadata><manifest><item id="empty" href="empty.xhtml" media-type="application/xhtml+xml"/><item id="chapter" href="chapter.xhtml" media-type="application/xhtml+xml"/>`+map[bool]string{true: `<item id="cover" href="cover.png" media-type="image/png" properties="cover-image"/>`}[withCover]+`</manifest><spine><itemref idref="empty"/><itemref idref="chapter"/></spine></package>`)
}

func writeTestEPUBWithOPF(t *testing.T, opf string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "book.epub")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	entries := map[string]string{
		"mimetype":               "application/epub+zip",
		"META-INF/container.xml": `<container><rootfiles><rootfile full-path="OEBPS/content.opf"/></rootfiles></container>`,
		"OEBPS/empty.xhtml":      "<html><body><img src=\"cover.png\" /></body></html>",
		"OEBPS/chapter.xhtml":    `<html><body><p>Chapter</p><img src="cover.png" /></body></html>`,
		"OEBPS/cover.png":        "\x89PNG\r\n\x1a\n",
		"OEBPS/content.opf":      opf,
	}
	for name, value := range entries {
		entry, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write([]byte(value)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}
