package lookupword

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
	"go.uber.org/mock/gomock"
)

func TestUseCaseReturnsSensesWithoutWaitingForTranslation(t *testing.T) {
	controller := gomock.NewController(t)
	normalizer := NewMockWordNormalizer(controller)
	morphology := NewMockMorphology(controller)
	repository := NewMockDictionaryRepository(controller)
	vocabulary := NewMockVocabularyReader(controller)
	normalizer.EXPECT().Normalize("Went").Return("went", nil)
	morphology.EXPECT().Candidates("went").Return([]string{"went", "go"})
	repository.EXPECT().Lookup(gomock.Any(), "went").Return(domain.LookupResult{}, nil)
	repository.EXPECT().Lookup(gomock.Any(), "go").Return(domain.LookupResult{LemmaID: 1, Senses: []domain.Sense{{PartOfSpeech: "verb"}}}, nil)
	vocabulary.EXPECT().IsSaved(gomock.Any(), "", int64(1)).Return(false, nil)

	result, err := New(normalizer, morphology, repository, vocabulary).Execute(context.Background(), Request{SelectedText: "Went"})
	if err != nil || len(result.Senses) != 1 {
		t.Fatalf("Execute() = %#v, %v", result, err)
	}
}
