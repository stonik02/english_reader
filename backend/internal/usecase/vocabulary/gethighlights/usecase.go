package gethighlights

import (
	"context"
	"sort"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	"github.com/deniskrylov/english-reader/backend/internal/service/tokenizer"
)

type Request struct {
	UserID    string
	BookID    string
	ChapterID string
}

type UseCase struct {
	chapters      ChapterReader
	vocabulary    VocabularyReader
	tokenizer     Tokenizer
	maxLemmas     int
	maxHighlights int
}

func New(chapters ChapterReader, vocabulary VocabularyReader, tokenizer Tokenizer, maxLemmas, maxHighlights int) *UseCase {
	return &UseCase{
		chapters:      chapters,
		vocabulary:    vocabulary,
		tokenizer:     tokenizer,
		maxLemmas:     maxLemmas,
		maxHighlights: maxHighlights,
	}
}

func (u *UseCase) Execute(ctx context.Context, request Request) ([]domain.HighlightToken, error) {
	if request.UserID == "" || request.BookID == "" || request.ChapterID == "" || u.maxLemmas <= 0 || u.maxHighlights <= 0 {
		return nil, domain.ErrInvalidInput
	}

	plainText, err := u.chapters.ChapterPlainText(ctx, request.UserID, request.BookID, request.ChapterID)
	if err != nil {
		return nil, err
	}

	lemmas, err := u.vocabulary.HighlightLemmas(ctx, request.UserID, u.maxLemmas)
	if err != nil {
		return nil, err
	}
	if len(lemmas) == 0 {
		return []domain.HighlightToken{}, nil
	}

	tokens, err := u.tokenizer.Tokenize(plainText)
	if err != nil {
		return nil, err
	}

	byLemma := make(map[string]domain.HighlightLemma, len(lemmas))
	for _, lemma := range lemmas {
		byLemma[lemma.Lemma] = lemma
	}

	byID := make(map[int64]*domain.HighlightToken)
	highlightCount := 0
	for _, token := range tokens {
		lemma, ok := byLemma[token.Lemma]
		if !ok {
			continue
		}
		if highlightCount >= u.maxHighlights {
			break
		}

		value, exists := byID[lemma.ID]
		if !exists {
			value = &domain.HighlightToken{
				LemmaID:   lemma.ID,
				Lemma:     lemma.Lemma,
				Positions: make([]int, 0),
				Texts:     make([]string, 0),
				MatchKind: matchKind(token.MatchKind),
			}
			byID[lemma.ID] = value
		}
		value.Positions = append(value.Positions, token.Start)
		if !contains(value.Texts, token.Text) {
			value.Texts = append(value.Texts, token.Text)
		}
		highlightCount++
	}

	result := make([]domain.HighlightToken, 0, len(byID))
	for _, value := range byID {
		result = append(result, *value)
	}
	sort.Slice(result, func(left, right int) bool {
		return result[left].Positions[0] < result[right].Positions[0]
	})

	return result, nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func matchKind(value tokenizer.MatchKind) domain.HighlightMatchKind {
	if value == tokenizer.MatchKindLemma {
		return domain.HighlightMatchKindLemma
	}
	return domain.HighlightMatchKindExactFallback
}
