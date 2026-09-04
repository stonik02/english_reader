package deletebook

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type testBooks struct{ files domain.StoredBookFiles }

func (b testBooks) Delete(context.Context, string) (domain.StoredBookFiles, error) {
	return b.files, nil
}

type testStorage struct{ deleted []string }

func (s *testStorage) Delete(path string) error { s.deleted = append(s.deleted, path); return nil }

func TestUseCaseDeletesSourceAndCoverFiles(t *testing.T) {
	storage := &testStorage{}
	usecase := New(testBooks{files: domain.StoredBookFiles{SourcePath: "/storage/original.epub", CoverPath: "/storage/cover.jpg"}}, storage)
	if err := usecase.Execute(context.Background(), "book-1"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(storage.deleted) != 2 || storage.deleted[0] != "/storage/original.epub" || storage.deleted[1] != "/storage/cover.jpg" {
		t.Fatalf("deleted files = %#v", storage.deleted)
	}
}
