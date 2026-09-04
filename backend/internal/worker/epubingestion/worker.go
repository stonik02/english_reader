package epubingestion

import (
	"context"
	"errors"
	"log/slog"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type Worker struct {
	jobs    JobRepository
	parser  EPUBParser
	storage CoverStorage
	logger  *slog.Logger
	period  time.Duration
}

func New(jobs JobRepository, parser EPUBParser, storage CoverStorage, logger *slog.Logger, period time.Duration) *Worker {
	if period <= 0 {
		period = time.Second
	}

	return &Worker{
		jobs:    jobs,
		parser:  parser,
		storage: storage,
		logger:  logger,
		period:  period,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.period)
	defer ticker.Stop()

	for {
		w.processOne(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *Worker) processOne(ctx context.Context) {
	job, err := w.jobs.ClaimNextIngestion(ctx)
	if errors.Is(err, domain.ErrNotFound) {
		return
	}
	if err != nil {
		w.logger.ErrorContext(ctx, "claim EPUB ingestion job", "error", err)
		return
	}

	result, err := w.parser.Parse(job.SourceFilePath, job.BookID)
	if err != nil {
		w.logger.WarnContext(ctx, "EPUB ingestion failed", "book_id", job.BookID, "error", err)
		if failErr := w.jobs.FailIngestion(ctx, job.BookID, "invalid_epub"); failErr != nil {
			w.logger.ErrorContext(ctx, "mark EPUB ingestion failed", "book_id", job.BookID, "error", failErr)
		}
		return
	}

	coverPath := ""
	if result.Cover != nil {
		coverPath, err = w.storage.StoreCover(job.BookID, result.Cover.ContentType, result.Cover.Data)
		if err != nil {
			w.logger.ErrorContext(ctx, "store EPUB cover", "book_id", job.BookID, "error", err)
			_ = w.jobs.FailIngestion(ctx, job.BookID, "invalid_cover")
			return
		}
	}
	// CompleteIngestion atomically writes chapters and sets books.status = 'ready'.
	if err := w.jobs.CompleteIngestion(ctx, job.BookID, result.Title, result.Author, coverPath, result.Chapters); err != nil {
		if coverPath != "" {
			w.storage.Remove(coverPath)
		}
		w.logger.ErrorContext(ctx, "complete EPUB ingestion", "book_id", job.BookID, "error", err)
		return
	}
	if result.Cover == nil {
		w.storage.RemoveCover(job.BookID)
	}
}
