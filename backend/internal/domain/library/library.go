package library

import (
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("book not found")
	ErrNotReady      = errors.New("book is not ready")
	ErrInvalidUpload = errors.New("invalid EPUB upload")
	ErrTooLarge      = errors.New("EPUB is too large")
)

type Book struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Author           string    `json:"author"`
	Status           string    `json:"status"`
	CoverURL         string    `json:"cover_url,omitempty"`
	CoverPath        string    `json:"-"`
	UploadedByUserID string    `json:"uploaded_by_user_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type UserBook struct {
	Book            Book      `json:"book"`
	AddedAt         time.Time `json:"added_at"`
	AddedVia        string    `json:"added_via"`
	ProgressPercent float64   `json:"progress_percent"`
}

type Page[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type Chapter struct {
	ID            string
	BookID        string
	Sequence      int
	Href          string
	StartCFI      string
	EndCFI        string
	SanitizedHTML string
	PlainText     string
}

type IngestionJob struct {
	BookID         string
	SourceFilePath string
}

type StoredBookFiles struct {
	SourcePath string
	CoverPath  string
}
