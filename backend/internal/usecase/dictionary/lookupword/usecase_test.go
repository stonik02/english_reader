package lookupword

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
	"go.uber.org/mock/gomock"
)

func TestUseCaseReturnsSensesWhenProviderFails(t *testing.T) {
	controller := gomock.NewController(t)
	normalizer := NewMockWordNormalizer(controller)
	morphology := NewMockMorphology(controller)
	repository := NewMockDictionaryRepository(controller)
	reader := NewMockReaderPort(controller)
	vocabulary := NewMockVocabularyReader(controller)
	provider := NewMockTranslationProvider(controller)
	normalizer.EXPECT().Normalize("Went").Return("went", nil)
	reader.EXPECT().ChapterPlainText(gomock.Any(), "", "").Return("We went home.", nil)
	morphology.EXPECT().Candidates("went").Return([]string{"went", "go"})
	repository.EXPECT().Lookup(gomock.Any(), "went").Return(domain.LookupResult{}, nil)
	repository.EXPECT().Lookup(gomock.Any(), "go").Return(domain.LookupResult{LemmaID: 1, Senses: []domain.Sense{{PartOfSpeech: "verb"}}}, nil)
	vocabulary.EXPECT().IsSaved(gomock.Any(), "", int64(1)).Return(false, nil)
	repository.EXPECT().CachedTranslation(gomock.Any(), "We went home.", "model").Return(domain.CachedTranslation{}, false, nil)
	provider.EXPECT().Translate(gomock.Any(), "We went home.").Return("", errors.New("timeout"))

	result, err := New(normalizer, morphology, repository, reader, vocabulary, provider, "model", time.Hour).Execute(context.Background(), Request{SelectedText: "Went", SentenceText: "We went home."})
	if err != nil || result.ProviderError == "" || len(result.Senses) != 1 {
		t.Fatalf("Execute() = %#v, %v", result, err)
	}
}
