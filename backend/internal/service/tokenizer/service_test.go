package tokenizer

import (
	"testing"

	"github.com/deniskrylov/english-reader/backend/internal/service/morphology"
	"github.com/deniskrylov/english-reader/backend/internal/service/wordnormalizer"
)

func TestServiceTokenize(t *testing.T) {
	service := New(wordnormalizer.New(64), morphology.New(), 100, 10)

	tests := []struct {
		name      string
		plainText string
		want      []Token
		wantErr   error
	}{
		{
			name:      "punctuation and irregular form",
			plainText: "Went, children!",
			want: []Token{
				{Text: "Went", Normalized: "went", Lemma: "go", Start: 0, End: 4, MatchKind: MatchKindLemma},
				{Text: "children", Normalized: "children", Lemma: "child", Start: 6, End: 14, MatchKind: MatchKindLemma},
			},
		},
		{
			name:      "unicode apostrophe",
			plainText: "Don’t stop.",
			want: []Token{
				{Text: "Don’t", Normalized: "don't", Lemma: "don't", Start: 0, End: len("Don’t"), MatchKind: MatchKindExactFallback},
				{Text: "stop", Normalized: "stop", Lemma: "stop", Start: len("Don’t "), End: len("Don’t stop"), MatchKind: MatchKindExactFallback},
			},
		},
		{
			name:      "xhtml is rejected",
			plainText: "<p>Went</p>",
			wantErr:   ErrInvalidPlainText,
		},
		{
			name:      "too many tokens",
			plainText: "one two three",
			wantErr:   ErrTooManyTokens,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			current := service
			if test.name == "too many tokens" {
				current = New(wordnormalizer.New(64), morphology.New(), 100, 2)
			}

			got, err := current.Tokenize(test.plainText)
			if err != test.wantErr {
				t.Fatalf("Tokenize() error = %v, want %v", err, test.wantErr)
			}
			if len(got) != len(test.want) {
				t.Fatalf("Tokenize() token count = %d, want %d", len(got), len(test.want))
			}
			for index := range test.want {
				if got[index] != test.want[index] {
					t.Errorf("Tokenize()[%d] = %#v, want %#v", index, got[index], test.want[index])
				}
			}
		})
	}
}

func TestServiceTokenizeInputLimit(t *testing.T) {
	service := New(wordnormalizer.New(64), morphology.New(), 3, 10)

	_, err := service.Tokenize("four")
	if err != ErrInputTooLarge {
		t.Fatalf("Tokenize() error = %v, want %v", err, ErrInputTooLarge)
	}
}
