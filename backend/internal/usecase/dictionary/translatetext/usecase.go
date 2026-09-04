package translatetext

import (
	"context"
	"errors"
	"strings"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
)

type Request struct {
	BookID    string
	ChapterID string
	Text      string
}

type UseCase struct {
	dictionary   TranslationCache
	reader       ReaderPort
	translator   TranslationProvider
	modelVersion string
	cacheTTL     time.Duration
}

func New(d TranslationCache, r ReaderPort, t TranslationProvider, version string, ttl time.Duration) *UseCase {
	return &UseCase{dictionary: d, reader: r, translator: t, modelVersion: version, cacheTTL: ttl}
}

func (u *UseCase) Execute(ctx context.Context, request Request) (domain.TextTranslationResponse, error) {
	text := strings.TrimSpace(request.Text)
	if text == "" {
		return domain.TextTranslationResponse{}, errors.New("translation text is required")
	}
	plainText, err := u.reader.ChapterPlainText(ctx, request.BookID, request.ChapterID)
	if err != nil {
		return domain.TextTranslationResponse{}, err
	}
	response := domain.TextTranslationResponse{ContextVerified: strings.Contains(plainText, text)}
	if cached, hit, err := u.dictionary.CachedTranslation(ctx, text, u.modelVersion); err == nil && hit {
		response.Text = cached.Text
		return response, nil
	}
	translated, err := u.translator.Translate(ctx, text)
	if err != nil {
		response.ProviderError = "translation provider unavailable"
		return response, nil
	}
	response.Text = translated
	_ = u.dictionary.PutTranslation(ctx, text, translated, u.modelVersion, u.cacheTTL)
	return response, nil
}
