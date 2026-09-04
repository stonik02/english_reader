package gethighlights

import (
	"context"
	"errors"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	"github.com/deniskrylov/english-reader/backend/internal/service/tokenizer"
	"go.uber.org/mock/gomock"
)

func TestUseCaseExecute(t *testing.T) {
	tests := []struct {
		name      string
		chapterID string
		prepare   func(*MockChapterReader, *MockVocabularyReader, *MockTokenizer)
		want      []domain.HighlightToken
		wantError error
	}{
		{
			name: "several matches are grouped by lemma",
			prepare: func(chapters *MockChapterReader, vocabulary *MockVocabularyReader, tokens *MockTokenizer) {
				chapters.EXPECT().ChapterPlainText(gomock.Any(), "user-1", "book-1", "chapter-1").Return("Went and went.", nil)
				vocabulary.EXPECT().HighlightLemmas(gomock.Any(), "user-1", 10).Return([]domain.HighlightLemma{{ID: 7, Lemma: "go"}}, nil)
				tokens.EXPECT().Tokenize("Went and went.").Return([]tokenizer.Token{
					{Lemma: "go", Start: 0, MatchKind: tokenizer.MatchKindLemma},
					{Lemma: "and", Start: 5, MatchKind: tokenizer.MatchKindExactFallback},
					{Lemma: "go", Start: 9, MatchKind: tokenizer.MatchKindLemma},
				}, nil)
			},
			want: []domain.HighlightToken{{
				LemmaID:   7,
				Lemma:     "go",
				Positions: []int{0, 9},
				MatchKind: domain.HighlightMatchKindLemma,
			}},
		},
		{
			name: "exact fallback is returned when morphology has no lemma",
			prepare: func(chapters *MockChapterReader, vocabulary *MockVocabularyReader, tokens *MockTokenizer) {
				chapters.EXPECT().ChapterPlainText(gomock.Any(), "user-1", "book-1", "chapter-1").Return("Reader", nil)
				vocabulary.EXPECT().HighlightLemmas(gomock.Any(), "user-1", 10).Return([]domain.HighlightLemma{{ID: 8, Lemma: "reader"}}, nil)
				tokens.EXPECT().Tokenize("Reader").Return([]tokenizer.Token{{Lemma: "reader", Start: 0, MatchKind: tokenizer.MatchKindExactFallback}}, nil)
			},
			want: []domain.HighlightToken{{
				LemmaID:   8,
				Lemma:     "reader",
				Positions: []int{0},
				MatchKind: domain.HighlightMatchKindExactFallback,
			}},
		},
		{
			name: "processing book is rejected",
			prepare: func(chapters *MockChapterReader, vocabulary *MockVocabularyReader, tokens *MockTokenizer) {
				chapters.EXPECT().ChapterPlainText(gomock.Any(), "user-1", "book-1", "chapter-1").Return("", domain.ErrNotReady)
			},
			wantError: domain.ErrNotReady,
		},
		{
			name:      "chapter outside book is rejected",
			chapterID: "chapter-2",
			prepare: func(chapters *MockChapterReader, vocabulary *MockVocabularyReader, tokens *MockTokenizer) {
				chapters.EXPECT().ChapterPlainText(gomock.Any(), "user-1", "book-1", "chapter-2").Return("", domain.ErrNotFound)
			},
			wantError: domain.ErrNotFound,
		},
		{
			name: "empty vocabulary returns no highlights without tokenizing",
			prepare: func(chapters *MockChapterReader, vocabulary *MockVocabularyReader, tokens *MockTokenizer) {
				chapters.EXPECT().ChapterPlainText(gomock.Any(), "user-1", "book-1", "chapter-1").Return("No words", nil)
				vocabulary.EXPECT().HighlightLemmas(gomock.Any(), "user-1", 10).Return([]domain.HighlightLemma{}, nil)
			},
			want: []domain.HighlightToken{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			chapters := NewMockChapterReader(controller)
			vocabulary := NewMockVocabularyReader(controller)
			tokens := NewMockTokenizer(controller)
			test.prepare(chapters, vocabulary, tokens)

			got, err := New(chapters, vocabulary, tokens, 10, 10).Execute(context.Background(), Request{
				UserID:    "user-1",
				BookID:    "book-1",
				ChapterID: chapterID(test.chapterID),
			})
			if !errors.Is(err, test.wantError) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantError)
			}
			if len(got) != len(test.want) {
				t.Fatalf("Execute() result count = %d, want %d", len(got), len(test.want))
			}
			for index := range test.want {
				if got[index].LemmaID != test.want[index].LemmaID || got[index].Lemma != test.want[index].Lemma || got[index].MatchKind != test.want[index].MatchKind {
					t.Errorf("Execute()[%d] = %#v, want %#v", index, got[index], test.want[index])
				}
				if len(got[index].Positions) != len(test.want[index].Positions) {
					t.Errorf("Execute()[%d].Positions = %#v, want %#v", index, got[index].Positions, test.want[index].Positions)
				}
			}
		})
	}
}

func chapterID(value string) string {
	if value == "" {
		return "chapter-1"
	}
	return value
}

func TestUseCaseExecuteInvalidRequest(t *testing.T) {
	controller := gomock.NewController(t)
	chapters := NewMockChapterReader(controller)
	vocabulary := NewMockVocabularyReader(controller)
	tokens := NewMockTokenizer(controller)

	_, err := New(chapters, vocabulary, tokens, 10, 10).Execute(context.Background(), Request{})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}
