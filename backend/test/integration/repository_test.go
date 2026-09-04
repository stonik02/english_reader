package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/deniskrylov/english-reader/backend/internal/database"
	repoauth "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/auth"
	repodictionary "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/dictionary"
	repolibrary "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/library"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		os.Exit(0)
	}

	ctx := context.Background()
	var err error
	pool, err = database.NewPool(ctx, databaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	if err := database.Up(ctx, pool); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestDictionaryRepositoryLookupAndCache(t *testing.T) {
	resetDatabase(t)
	ctx := context.Background()
	var lemmaID int64
	err := pool.QueryRow(ctx, `INSERT INTO dictionary_lemmas (language,lemma,source,source_version) VALUES ('en','go','fixture','v1') RETURNING id`).Scan(&lemmaID)
	if err != nil {
		t.Fatalf("insert lemma: %v", err)
	}
	_, err = pool.Exec(ctx, `INSERT INTO dictionary_senses (lemma_id,part_of_speech,translations,source_url,attribution,position) VALUES ($1,'verb','["идти"]','https://example.test/go','fixture',0),($1,'noun','["ход"]','https://example.test/go','fixture',1)`, lemmaID)
	if err != nil {
		t.Fatalf("insert senses: %v", err)
	}
	repository := repodictionary.New(pool)
	result, err := repository.Lookup(ctx, "go")
	if err != nil || len(result.Senses) != 2 || result.Senses[0].PartOfSpeech != "verb" {
		t.Fatalf("Lookup() = %#v, %v", result, err)
	}
	if err := repository.PutTranslation(ctx, "We go home.", "Мы идём домой.", "lt-test", time.Hour); err != nil {
		t.Fatalf("PutTranslation() error = %v", err)
	}
	value, hit, err := repository.CachedTranslation(ctx, "We go home.", "lt-test")
	if err != nil || !hit || value.Text != "Мы идём домой." {
		t.Fatalf("CachedTranslation() = %#v, %v, %v", value, hit, err)
	}
}

func TestDictionaryRepositoryExpiredCacheIsMiss(t *testing.T) {
	resetDatabase(t)
	ctx := context.Background()
	repository := repodictionary.New(pool)
	if err := repository.PutTranslation(ctx, "Expired sentence.", "Просрочено.", "lt-test", -time.Hour); err != nil {
		t.Fatalf("PutTranslation() error = %v", err)
	}
	_, hit, err := repository.CachedTranslation(ctx, "Expired sentence.", "lt-test")
	if err != nil || hit {
		t.Fatalf("CachedTranslation() hit = %v, error = %v; want expired cache miss", hit, err)
	}
}

func TestDictionaryRepositoryLookupUsesNewestImportVersion(t *testing.T) {
	resetDatabase(t)
	ctx := context.Background()
	var oldLemmaID, newestLemmaID int64
	if err := pool.QueryRow(ctx, `INSERT INTO dictionary_lemmas (language,lemma,source,source_version) VALUES ('en','fair','fixture','2026-08-05') RETURNING id`).Scan(&oldLemmaID); err != nil {
		t.Fatalf("insert old lemma: %v", err)
	}
	if err := pool.QueryRow(ctx, `INSERT INTO dictionary_lemmas (language,lemma,source,source_version) VALUES ('en','fair','fixture','2026-09-04') RETURNING id`).Scan(&newestLemmaID); err != nil {
		t.Fatalf("insert newest lemma: %v", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO dictionary_senses (lemma_id,part_of_speech,translations,source_url,attribution,position) VALUES ($1,'noun','["ярмарка"]','https://example.test/fair','fixture',0),($2,'adjective','["справедливый"]','https://example.test/fair','fixture',0),($2,'noun','["ярмарка"]','https://example.test/fair','fixture',1),($2,'noun','["ярмарка"]','https://example.test/fair','fixture',2)`, oldLemmaID, newestLemmaID); err != nil {
		t.Fatalf("insert senses: %v", err)
	}

	result, err := repodictionary.New(pool).Lookup(ctx, "fair")
	if err != nil || result.LemmaID != newestLemmaID || len(result.Senses) != 2 {
		t.Fatalf("Lookup() = %#v, %v", result, err)
	}
}

func resetDatabase(t *testing.T) {
	t.Helper()
	_, err := pool.Exec(context.Background(), `TRUNCATE users, dictionary_lemmas CASCADE`)
	if err != nil {
		t.Fatalf("truncate database: %v", err)
	}
}

func TestAuthRepositorySessionLifecycle(t *testing.T) {
	resetDatabase(t)
	ctx := context.Background()
	repository := repoauth.New(pool)
	hash := []byte("refresh-hash-1")
	user, err := repository.CreateUserWithSession(ctx, "reader@example.com", "hash", hash, time.Now().Add(time.Hour), "test")
	if err != nil {
		t.Fatalf("CreateUserWithSession() error = %v", err)
	}
	stored, _, err := repository.FindUserByEmail(ctx, user.Email)
	if err != nil || stored.ID != user.ID {
		t.Fatalf("FindUserByEmail() = %#v, %v", stored, err)
	}
	if err := repository.RevokeSession(ctx, hash); err != nil {
		t.Fatalf("RevokeSession() error = %v", err)
	}
	if _, err := repository.FindUserByRefreshHash(ctx, hash); err == nil {
		t.Fatal("revoked refresh session remained available")
	}
}

func TestLibraryRepositoryPersonalLibraryIsolation(t *testing.T) {
	resetDatabase(t)
	ctx := context.Background()
	users := repoauth.New(pool)
	first, err := users.CreateUserWithSession(ctx, "one@example.com", "hash", []byte("one"), time.Now().Add(time.Hour), "test")
	if err != nil {
		t.Fatalf("create first user: %v", err)
	}
	second, err := users.CreateUserWithSession(ctx, "two@example.com", "hash", []byte("two"), time.Now().Add(time.Hour), "test")
	if err != nil {
		t.Fatalf("create second user: %v", err)
	}
	books := repolibrary.New(pool)
	book, err := books.Create(ctx, uuid.NewString(), first.ID, "Book", "/private/book.epub", []byte("book-hash"))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := books.CompleteIngestion(ctx, book.ID, "Book", "", "", nil); err != nil {
		t.Fatalf("CompleteIngestion() error = %v", err)
	}
	if _, err := books.Add(ctx, first.ID, book.ID, "button"); err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	firstPage, err := books.ListMine(ctx, first.ID, "", 20)
	if err != nil || len(firstPage.Items) != 1 {
		t.Fatalf("first library = %#v, %v", firstPage, err)
	}
	secondPage, err := books.ListMine(ctx, second.ID, "", 20)
	if err != nil || len(secondPage.Items) != 0 {
		t.Fatalf("second library = %#v, %v", secondPage, err)
	}
}
