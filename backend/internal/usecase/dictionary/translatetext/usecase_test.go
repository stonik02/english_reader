package translatetext

import (
	"context"
	"testing"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
)

func TestUseCaseReturnsCachedTranslationWithoutCallingProvider(t *testing.T) {
	cache := &translationCache{cached: domain.CachedTranslation{Text: "привет"}, hit: true}
	provider := &translationProvider{}
	result, err := New(cache, chapterReader{plainText: "Hello, world."}, provider, "test", time.Hour).Execute(
		context.Background(),
		Request{BookID: "book", ChapterID: "chapter", Text: "Hello"},
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Text != "привет" || !result.ContextVerified {
		t.Fatalf("Execute() = %#v, want cached verified translation", result)
	}
	if provider.calls != 0 {
		t.Fatalf("provider calls = %d, want 0", provider.calls)
	}
}

func TestUseCaseKeepsDictionaryLookupIndependentWhenProviderFails(t *testing.T) {
	provider := &translationProvider{err: context.DeadlineExceeded}
	result, err := New(&translationCache{}, chapterReader{plainText: "A long sentence."}, provider, "test", time.Hour).Execute(
		context.Background(),
		Request{BookID: "book", ChapterID: "chapter", Text: "A long sentence"},
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.ProviderError == "" {
		t.Fatalf("Execute() = %#v, want provider error in response", result)
	}
}

type chapterReader struct {
	plainText string
}

func (r chapterReader) ChapterPlainText(context.Context, string, string) (string, error) {
	return r.plainText, nil
}

type translationCache struct {
	cached domain.CachedTranslation
	hit    bool
}

func (c *translationCache) CachedTranslation(context.Context, string, string) (domain.CachedTranslation, bool, error) {
	return c.cached, c.hit, nil
}

func (c *translationCache) PutTranslation(context.Context, string, string, string, time.Duration) error {
	return nil
}

type translationProvider struct {
	calls int
	err   error
}

func (p *translationProvider) Translate(context.Context, string) (string, error) {
	p.calls++
	if p.err != nil {
		return "", p.err
	}
	return "перевод", nil
}
