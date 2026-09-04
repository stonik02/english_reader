package lookupword

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
)

type Request struct {
	UserID       string
	BookID       string
	ChapterID    string
	SelectedText string
	SentenceText string
}
type UseCase struct {
	normalizer WordNormalizer
	morphology Morphology
	dictionary DictionaryRepository
	vocabulary VocabularyReader
}

func New(n WordNormalizer, m Morphology, d DictionaryRepository, v VocabularyReader) *UseCase {
	return &UseCase{normalizer: n, morphology: m, dictionary: d, vocabulary: v}
}
func (u *UseCase) Execute(ctx context.Context, request Request) (domain.LookupResponse, error) {
	word, err := u.normalizer.Normalize(request.SelectedText)
	if err != nil {
		return domain.LookupResponse{}, err
	}
	lemma := word
	result := domain.LookupResult{}
	for _, candidate := range u.morphology.Candidates(word) {
		value, lookupErr := u.dictionary.Lookup(ctx, candidate)
		if lookupErr != nil {
			return domain.LookupResponse{}, lookupErr
		}
		if value.LemmaID != 0 {
			lemma = candidate
			result = value
			break
		}
	}
	response := domain.LookupResponse{LemmaID: result.LemmaID, NormalizedLemma: lemma, Senses: result.Senses, Source: result.Source, SourceVersion: result.SourceVersion}
	if result.LemmaID != 0 {
		alreadySaved, savedErr := u.vocabulary.IsSaved(ctx, request.UserID, result.LemmaID)
		if savedErr != nil {
			return domain.LookupResponse{}, savedErr
		}
		response.AlreadySaved = alreadySaved
	}
	return response, nil
}
