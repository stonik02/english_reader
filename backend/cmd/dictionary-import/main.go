package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type record struct {
	Lemma        string   `json:"lemma"`
	Language     string   `json:"language"`
	PartOfSpeech string   `json:"part_of_speech"`
	Translations []string `json:"translations"`
	ExampleEN    string   `json:"example_en"`
	ExampleRU    string   `json:"example_ru"`
	SourceURL    string   `json:"source_url"`
	Attribution  string   `json:"attribution"`
	Position     int      `json:"position"`
}

func main() {
	var databaseURL, filePath, source, version string
	flag.StringVar(&databaseURL, "database-url", os.Getenv("DATABASE_URL"), "PostgreSQL URL")
	flag.StringVar(&filePath, "file", "", "JSONL dictionary dump")
	flag.StringVar(&source, "source", "wiktionary", "dictionary source")
	flag.StringVar(&version, "version", "", "source version")
	flag.Parse()
	if databaseURL == "" || filePath == "" || version == "" {
		fmt.Fprintln(os.Stderr, "database-url, file and version are required; set DATABASE_URL or pass -database-url")
		os.Exit(2)
	}
	if err := run(context.Background(), databaseURL, filePath, source, version); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, databaseURL, filePath, source, version string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read dump: %w", err)
	}
	checksum := sha256.Sum256(contents)
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	runID := uuid.NewString()
	if _, err = pool.Exec(ctx, `INSERT INTO dictionary_import_runs (id,source,source_version,checksum) VALUES ($1,$2,$3,$4)`, runID, source, version, checksum[:]); err != nil {
		return err
	}
	count, err := importRecords(ctx, pool, contents, source, version)
	if err != nil {
		_, _ = pool.Exec(ctx, `UPDATE dictionary_import_runs SET finished_at=NOW(),error_detail=$2 WHERE id=$1`, runID, err.Error())
		return err
	}
	_, err = pool.Exec(ctx, `UPDATE dictionary_import_runs SET finished_at=NOW(),imported_lemmas=$2 WHERE id=$1`, runID, count)
	return err
}

func importRecords(ctx context.Context, pool *pgxpool.Pool, contents []byte, source, version string) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(contents))
	scanner.Buffer(make([]byte, 1024), 2<<20)
	count := 0
	clearedLemmaIDs := make(map[int64]struct{})
	for scanner.Scan() {
		var value record
		if err := json.Unmarshal(scanner.Bytes(), &value); err != nil {
			return count, fmt.Errorf("decode JSONL: %w", err)
		}
		if value.Lemma == "" || value.Language != "en" || value.PartOfSpeech == "" || len(value.Translations) == 0 || value.SourceURL == "" || value.Attribution == "" {
			return count, fmt.Errorf("invalid record for lemma %q", value.Lemma)
		}
		translations, _ := json.Marshal(value.Translations)
		tx, err := pool.Begin(ctx)
		if err != nil {
			return count, err
		}
		var lemmaID int64
		err = tx.QueryRow(ctx, `INSERT INTO dictionary_lemmas (language,lemma,source,source_version) VALUES ($1,$2,$3,$4) ON CONFLICT (language,lemma,source,source_version) DO UPDATE SET lemma=EXCLUDED.lemma RETURNING id`, value.Language, value.Lemma, source, version).Scan(&lemmaID)
		if _, cleared := clearedLemmaIDs[lemmaID]; err == nil && !cleared {
			_, err = tx.Exec(ctx, `DELETE FROM dictionary_senses WHERE lemma_id=$1`, lemmaID)
			if err == nil {
				clearedLemmaIDs[lemmaID] = struct{}{}
			}
		}
		if err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO dictionary_senses (lemma_id,part_of_speech,translations,example_en,example_ru,source_url,attribution,position) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, lemmaID, value.PartOfSpeech, translations, value.ExampleEN, value.ExampleRU, value.SourceURL, value.Attribution, value.Position)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return count, err
		}
		if err = tx.Commit(ctx); err != nil {
			return count, err
		}
		count++
	}
	return count, scanner.Err()
}
