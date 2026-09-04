package gethighlights

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	"github.com/deniskrylov/english-reader/backend/internal/service/tokenizer"
)

type ChapterReader interface {
	ChapterPlainText(context.Context, string, string, string) (string, error)
}

type VocabularyReader interface {
	HighlightLemmas(context.Context, string, int) ([]domain.HighlightLemma, error)
}

type Tokenizer interface {
	Tokenize(string) ([]tokenizer.Token, error)
}
