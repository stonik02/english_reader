package tokenizer

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrInvalidPlainText = errors.New("invalid plain text")
	ErrInputTooLarge    = errors.New("plain text is too large")
	ErrTooManyTokens    = errors.New("too many tokens")
)

type MatchKind string

const (
	MatchKindLemma         MatchKind = "lemma"
	MatchKindExactFallback MatchKind = "exact_fallback"
)

type Token struct {
	Text       string
	Normalized string
	Lemma      string
	Start      int
	End        int
	MatchKind  MatchKind
}

type Service struct {
	normalizer WordNormalizer
	morphology Morphology
	maxRunes   int
	maxTokens  int
}

func New(normalizer WordNormalizer, morphology Morphology, maxRunes, maxTokens int) *Service {
	return &Service{
		normalizer: normalizer,
		morphology: morphology,
		maxRunes:   maxRunes,
		maxTokens:  maxTokens,
	}
}

func (s *Service) Tokenize(plainText string) ([]Token, error) {
	if plainText == "" || strings.ContainsAny(plainText, "<>") {
		return nil, ErrInvalidPlainText
	}
	if s.maxRunes <= 0 || s.maxTokens <= 0 || len([]rune(plainText)) > s.maxRunes {
		return nil, ErrInputTooLarge
	}

	var tokens []Token
	start := -1

	flush := func(end int) error {
		if start < 0 {
			return nil
		}

		value := plainText[start:end]
		start = -1

		normalized, err := s.normalizer.Normalize(strings.ReplaceAll(value, "’", "'"))
		if err != nil {
			return nil
		}

		if len(tokens) >= s.maxTokens {
			return ErrTooManyTokens
		}

		lemma, known := s.morphology.Lemma(normalized)
		kind := MatchKindExactFallback
		if known {
			kind = MatchKindLemma
		}

		tokens = append(tokens, Token{
			Text:       value,
			Normalized: normalized,
			Lemma:      lemma,
			Start:      end - len(value),
			End:        end,
			MatchKind:  kind,
		})

		return nil
	}

	for offset, runeValue := range plainText {
		if unicode.IsLetter(runeValue) {
			if start < 0 {
				start = offset
			}
			continue
		}
		if isApostrophe(runeValue) && start >= 0 {
			continue
		}
		if err := flush(offset); err != nil {
			return nil, err
		}
	}

	if err := flush(len(plainText)); err != nil {
		return nil, err
	}

	return tokens, nil
}

func isApostrophe(value rune) bool {
	return value == '\'' || value == '’'
}
