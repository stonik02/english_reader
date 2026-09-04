package lookupword

import (
	"context"
	"strings"
	"time"

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
	normalizer   WordNormalizer
	morphology   Morphology
	dictionary   DictionaryRepository
	reader       ReaderPort
	vocabulary   VocabularyReader
	translator   TranslationProvider
	modelVersion string
	cacheTTL     time.Duration
}

func New(n WordNormalizer, m Morphology, d DictionaryRepository, r ReaderPort, v VocabularyReader, t TranslationProvider, version string, ttl time.Duration) *UseCase {
	return &UseCase{normalizer: n, morphology: m, dictionary: d, reader: r, vocabulary: v, translator: t, modelVersion: version, cacheTTL: ttl}
}
func (u *UseCase) Execute(ctx context.Context, request Request) (domain.LookupResponse, error) {
	plainText, err := u.reader.ChapterPlainText(ctx, request.BookID, request.ChapterID)
	if err != nil {
		return domain.LookupResponse{}, err
	}
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
	response := domain.LookupResponse{LemmaID: result.LemmaID, NormalizedLemma: lemma, Senses: result.Senses, ContextVerified: strings.Contains(plainText, request.SentenceText), Source: result.Source, SourceVersion: result.SourceVersion}
	if result.LemmaID != 0 {
		alreadySaved, savedErr := u.vocabulary.IsSaved(ctx, request.UserID, result.LemmaID)
		if savedErr != nil {
			return domain.LookupResponse{}, savedErr
		}
		response.AlreadySaved = alreadySaved
	}
	if cached, hit, err := u.dictionary.CachedTranslation(ctx, request.SentenceText, u.modelVersion); err == nil && hit {
		response.SentenceTranslation = cached.Text
		return response, nil
	}
	translated, err := u.translator.Translate(ctx, request.SentenceText)
	if err != nil {
		response.ProviderError = "translation provider unavailable"
		return response, nil
	}
	response.SentenceTranslation = translated
	_ = u.dictionary.PutTranslation(ctx, request.SentenceText, translated, u.modelVersion, u.cacheTTL)
	return response, nil
}
