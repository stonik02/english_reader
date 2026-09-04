package vocabulary

import (
	"errors"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid vocabulary input")
	ErrNotFound     = errors.New("vocabulary resource not found")
	ErrInvalidSense = errors.New("selected sense does not belong to lemma")
	ErrNotReady     = errors.New("book is not ready")
)

type Sense struct {
	ID           int64
	PartOfSpeech string
	Translations []string
	ExampleEN    string
	ExampleRU    string
}

type Entry struct {
	ID          string
	LemmaID     int64
	Lemma       string
	SourceForm  string
	ChosenSense *Sense
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SaveRequest struct {
	LemmaID       int64
	ChosenSenseID *int64
	SourceForm    string
}

type Cursor struct {
	CreatedAt time.Time
	ID        string
}

type ListRequest struct {
	Cursor *Cursor
	Limit  int
	Query  string
}

type HighlightMatchKind string

const (
	HighlightMatchKindLemma         HighlightMatchKind = "lemma"
	HighlightMatchKindExactFallback HighlightMatchKind = "exact_fallback"
)

type HighlightLemma struct {
	ID    int64
	Lemma string
}

type HighlightToken struct {
	LemmaID   int64
	Lemma     string
	Positions []int
	Texts     []string
	MatchKind HighlightMatchKind
}
