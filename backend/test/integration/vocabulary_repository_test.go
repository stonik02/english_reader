package integration

import (
	"context"
	"testing"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	repoauth "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/auth"
	repovocabulary "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/vocabulary"
)

func TestVocabularyRepositorySaveIsIdempotentAndIsolated(t *testing.T) {
	resetDatabase(t)
	ctx := context.Background()
	firstUser, secondUser := createVocabularyUsers(t, ctx)
	firstLemmaID, firstSenseID, secondLemmaID := createVocabularyDictionary(t, ctx)
	repository := repovocabulary.New(pool)

	entry, alreadySaved, err := repository.Save(ctx, firstUser, domain.SaveRequest{
		LemmaID:       firstLemmaID,
		ChosenSenseID: &firstSenseID,
		SourceForm:    "went",
	})
	if err != nil || alreadySaved || entry.ID == "" {
		t.Fatalf("Save() = %#v, %t, %v", entry, alreadySaved, err)
	}
	duplicate, alreadySaved, err := repository.Save(ctx, firstUser, domain.SaveRequest{
		LemmaID:    firstLemmaID,
		SourceForm: "goes",
	})
	if err != nil || !alreadySaved || duplicate.ID != entry.ID || duplicate.SourceForm != "went" || duplicate.ChosenSense == nil {
		t.Fatalf("duplicate Save() = %#v, %t, %v", duplicate, alreadySaved, err)
	}
	if _, _, err := repository.Save(ctx, firstUser, domain.SaveRequest{LemmaID: firstLemmaID, ChosenSenseID: &secondLemmaID, SourceForm: "went"}); err != domain.ErrInvalidSense {
		t.Fatalf("Save() invalid sense error = %v, want %v", err, domain.ErrInvalidSense)
	}

	entries, err := repository.List(ctx, secondUser, domain.ListRequest{Limit: 10})
	if err != nil || len(entries) != 0 {
		t.Fatalf("second user List() = %#v, %v", entries, err)
	}
	if err := repository.DeleteByID(ctx, secondUser, entry.ID); err != nil {
		t.Fatalf("DeleteByID() error = %v", err)
	}
	entries, err = repository.List(ctx, firstUser, domain.ListRequest{Limit: 10})
	if err != nil || len(entries) != 1 {
		t.Fatalf("first user entry was deleted by another user: %#v, %v", entries, err)
	}
}

func TestVocabularyRepositoryListsWithCursorAndSearch(t *testing.T) {
	resetDatabase(t)
	ctx := context.Background()
	userID, _ := createVocabularyUsers(t, ctx)
	firstLemmaID, _, secondLemmaID := createVocabularyDictionary(t, ctx)
	var thirdLemmaID int64
	if err := pool.QueryRow(ctx, `INSERT INTO dictionary_lemmas (language, lemma, source, source_version) VALUES ('en', 'run', 'fixture', 'v1') RETURNING id`).Scan(&thirdLemmaID); err != nil {
		t.Fatalf("insert third lemma: %v", err)
	}
	repository := repovocabulary.New(pool)
	for _, item := range []struct {
		lemmaID int64
		form    string
	}{
		{firstLemmaID, "went"},
		{secondLemmaID, "books"},
		{thirdLemmaID, "running"},
	} {
		if _, _, err := repository.Save(ctx, userID, domain.SaveRequest{LemmaID: item.lemmaID, SourceForm: item.form}); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
		time.Sleep(time.Millisecond)
	}

	page, err := repository.List(ctx, userID, domain.ListRequest{Limit: 2, Query: ""})
	if err != nil || len(page) != 2 {
		t.Fatalf("first List() = %#v, %v", page, err)
	}
	next, err := repository.List(ctx, userID, domain.ListRequest{
		Cursor: &domain.Cursor{CreatedAt: page[1].CreatedAt, ID: page[1].ID},
		Limit:  2,
	})
	if err != nil || len(next) != 1 || next[0].ID == page[0].ID || next[0].ID == page[1].ID {
		t.Fatalf("cursor List() = %#v, %v", next, err)
	}
	search, err := repository.List(ctx, userID, domain.ListRequest{Limit: 10, Query: "boo"})
	if err != nil || len(search) != 1 || search[0].Lemma != "book" {
		t.Fatalf("search List() = %#v, %v", search, err)
	}
}

func createVocabularyUsers(t *testing.T, ctx context.Context) (string, string) {
	t.Helper()
	repository := repoauth.New(pool)
	first, err := repository.CreateUserWithSession(ctx, "vocabulary-one@example.com", "hash", []byte("vocabulary-one"), time.Now().Add(time.Hour), "test")
	if err != nil {
		t.Fatalf("create first user: %v", err)
	}
	second, err := repository.CreateUserWithSession(ctx, "vocabulary-two@example.com", "hash", []byte("vocabulary-two"), time.Now().Add(time.Hour), "test")
	if err != nil {
		t.Fatalf("create second user: %v", err)
	}
	return first.ID, second.ID
}

func createVocabularyDictionary(t *testing.T, ctx context.Context) (int64, int64, int64) {
	t.Helper()
	var firstLemmaID int64
	if err := pool.QueryRow(ctx, `INSERT INTO dictionary_lemmas (language, lemma, source, source_version) VALUES ('en', 'go', 'fixture', 'v1') RETURNING id`).Scan(&firstLemmaID); err != nil {
		t.Fatalf("insert first lemma: %v", err)
	}
	var firstSenseID int64
	if err := pool.QueryRow(ctx, `INSERT INTO dictionary_senses (lemma_id, part_of_speech, translations, source_url, attribution) VALUES ($1, 'verb', '["идти"]', 'https://example.test/go', 'fixture') RETURNING id`, firstLemmaID).Scan(&firstSenseID); err != nil {
		t.Fatalf("insert first sense: %v", err)
	}
	var secondLemmaID int64
	if err := pool.QueryRow(ctx, `INSERT INTO dictionary_lemmas (language, lemma, source, source_version) VALUES ('en', 'book', 'fixture', 'v1') RETURNING id`).Scan(&secondLemmaID); err != nil {
		t.Fatalf("insert second lemma: %v", err)
	}
	return firstLemmaID, firstSenseID, secondLemmaID
}
