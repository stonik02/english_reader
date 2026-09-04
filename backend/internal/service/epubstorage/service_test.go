package epubstorage

import (
	"bytes"
	"os"
	"testing"
)

func TestServiceStoreAndMove(t *testing.T) {
	service := New(t.TempDir(), 1024)
	stored, err := service.StoreTemporary("book.epub", bytes.NewReader([]byte("PK\x03\x04content")))
	if err != nil {
		t.Fatalf("StoreTemporary() error = %v", err)
	}
	path, err := service.MoveToBook(stored.Temp, "book-1")
	if err != nil {
		t.Fatalf("MoveToBook() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stored EPUB is absent: %v", err)
	}
}

func TestServiceRejectsInvalidExtension(t *testing.T) {
	_, err := New(t.TempDir(), 1024).StoreTemporary("book.pdf", bytes.NewReader([]byte("PK\x03\x04")))
	if err == nil {
		t.Fatal("StoreTemporary() accepted a non-EPUB extension")
	}
}

func TestServiceStoresCoverWithSafeExtension(t *testing.T) {
	service := New(t.TempDir(), 1024)
	path, err := service.StoreCover("book-1", "image/png", []byte("\x89PNG\r\n\x1a\n"))
	if err != nil {
		t.Fatalf("StoreCover() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("stored cover is absent: %v", err)
	}
	if _, err := service.StoreCover("book-1", "image/svg+xml", []byte("<svg/>")); err == nil {
		t.Fatal("StoreCover() accepted SVG")
	}
}
