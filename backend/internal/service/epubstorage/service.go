package epubstorage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

const (
	DefaultMaxSize      = 50 << 20
	DefaultMaxCoverSize = 5 << 20
)

type Service struct {
	root    string
	maxSize int64
}

type StoredFile struct {
	SHA256 [sha256.Size]byte
	Temp   string
}

func New(root string, maxSize int64) *Service {
	if maxSize <= 0 {
		maxSize = DefaultMaxSize
	}
	return &Service{root: root, maxSize: maxSize}
}

func (s *Service) StoreTemporary(name string, source io.Reader) (StoredFile, error) {
	if !strings.EqualFold(filepath.Ext(name), ".epub") {
		return StoredFile{}, domain.ErrInvalidUpload
	}

	if err := os.MkdirAll(filepath.Join(s.root, "tmp"), 0o750); err != nil {
		return StoredFile{}, fmt.Errorf("create temporary directory: %w", err)
	}

	temp, err := os.CreateTemp(filepath.Join(s.root, "tmp"), "upload-*.epub")
	if err != nil {
		return StoredFile{}, fmt.Errorf("create temporary file: %w", err)
	}

	keep := false

	defer func() {
		_ = temp.Close()
		if !keep {
			_ = os.Remove(temp.Name())
		}
	}()

	hash := sha256.New()
	limited := io.LimitReader(source, s.maxSize+1)

	written, err := io.Copy(io.MultiWriter(temp, hash), limited)
	if err != nil {
		return StoredFile{}, fmt.Errorf("write upload: %w", err)
	}

	if written > s.maxSize {
		return StoredFile{}, domain.ErrTooLarge
	}

	if written < 4 {
		return StoredFile{}, domain.ErrInvalidUpload
	}

	if _, err = temp.Seek(0, io.SeekStart); err != nil {
		return StoredFile{}, fmt.Errorf("rewind upload: %w", err)
	}

	magic := make([]byte, 4)
	if _, err = io.ReadFull(temp, magic); err != nil || string(magic[:2]) != "PK" {
		return StoredFile{}, domain.ErrInvalidUpload
	}

	var sum [sha256.Size]byte
	copy(sum[:], hash.Sum(nil))
	keep = true

	return StoredFile{SHA256: sum, Temp: temp.Name()}, nil
}

func (s *Service) MoveToBook(tempPath, bookID string) (string, error) {
	directory := filepath.Join(s.root, "books", bookID)
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return "", fmt.Errorf("create book directory: %w", err)
	}

	path := filepath.Join(directory, "original.epub")
	if err := os.Rename(tempPath, path); err != nil {
		return "", fmt.Errorf("move EPUB: %w", err)
	}

	return path, nil
}

func (s *Service) StoreCover(bookID, contentType string, data []byte) (string, error) {
	extension := map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/gif": ".gif", "image/webp": ".webp"}[contentType]
	if extension == "" || len(data) == 0 || int64(len(data)) > DefaultMaxCoverSize {
		return "", domain.ErrInvalidUpload
	}
	directory := filepath.Join(s.root, "books", bookID)
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return "", fmt.Errorf("create book directory: %w", err)
	}
	path := filepath.Join(directory, "cover"+extension)
	if err := os.WriteFile(path, data, 0o640); err != nil {
		return "", fmt.Errorf("write cover: %w", err)
	}
	for _, otherExtension := range []string{".jpg", ".png", ".gif", ".webp"} {
		if otherExtension != extension {
			_ = os.Remove(filepath.Join(directory, "cover"+otherExtension))
		}
	}
	return path, nil
}

// RemoveCover removes every supported cover variant for a book. It is used when
// a reprocessed EPUB has no cover, so an old cover is never shown by mistake.
func (s *Service) RemoveCover(bookID string) {
	directory := filepath.Join(s.root, "books", bookID)
	for _, extension := range []string{".jpg", ".png", ".gif", ".webp"} {
		_ = os.Remove(filepath.Join(directory, "cover"+extension))
	}
}

func (s *Service) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (s *Service) Remove(path string) { _ = os.Remove(path) }

func (s *Service) Delete(path string) error { return os.Remove(path) }
