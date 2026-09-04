package epubingestion

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"github.com/deniskrylov/english-reader/backend/internal/service/epubparser"
)

type JobRepository interface {
	ClaimNextIngestion(context.Context) (domain.IngestionJob, error)
	CompleteIngestion(context.Context, string, string, string, string, []domain.Chapter) error
	FailIngestion(context.Context, string, string) error
}

type CoverStorage interface {
	StoreCover(string, string, []byte) (string, error)
	RemoveCover(string)
	Remove(string)
}

type EPUBParser interface {
	Parse(string, string) (epubparser.Result, error)
}
