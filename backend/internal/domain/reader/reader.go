package reader

import (
	"errors"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid reader input")
	ErrNotFound     = errors.New("reader resource not found")
	ErrNotReady     = errors.New("book is not ready")
)

type Chapter struct {
	ID            string `json:"id"`
	Href          string `json:"href"`
	Sequence      int    `json:"sequence"`
	StartCFI      string `json:"start_cfi"`
	EndCFI        string `json:"end_cfi"`
	SanitizedHTML string `json:"sanitized_html"`
	TotalChapters int    `json:"total_chapters"`
}

type Progress struct {
	ChapterID       string    `json:"chapter_id"`
	EPUBCFI         string    `json:"epub_cfi"`
	ProgressPercent float64   `json:"progress_percent"`
	Revision        int64     `json:"revision"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Settings struct {
	FontScale      int     `json:"font_scale"`
	Theme          string  `json:"theme"`
	LineHeight     float64 `json:"line_height"`
	HighlightColor string  `json:"highlight_color"`
}
type State struct {
	Chapter  Chapter  `json:"chapter"`
	Progress Progress `json:"progress"`
	Settings Settings `json:"settings"`
}
