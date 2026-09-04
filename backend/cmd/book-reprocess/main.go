// Command book-reprocess schedules existing EPUBs for parsing from their
// original stored files. It never uploads, replaces, or deletes an EPUB.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/deniskrylov/english-reader/backend/internal/repository/postgres/library"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var databaseURL, bookID string
	var all bool
	flag.StringVar(&databaseURL, "database-url", os.Getenv("DATABASE_URL"), "PostgreSQL URL")
	flag.StringVar(&bookID, "book-id", "", "book UUID to reprocess")
	flag.BoolVar(&all, "all", false, "reprocess every stored EPUB")
	flag.Parse()

	if databaseURL == "" || (bookID == "" && !all) || (bookID != "" && all) {
		fmt.Fprintln(os.Stderr, "database-url and exactly one of -book-id or -all are required")
		os.Exit(2)
	}
	if err := run(context.Background(), databaseURL, bookID, all); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, databaseURL, bookID string, all bool) error {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	defer pool.Close()

	repository := library.New(pool)
	if all {
		count, err := repository.RequeueAllIngestions(ctx)
		if err != nil {
			return fmt.Errorf("queue all books: %w", err)
		}
		fmt.Printf("Queued %d book(s) for EPUB reprocessing.\n", count)
		return nil
	}
	if err := repository.RequeueIngestion(ctx, bookID); err != nil {
		return fmt.Errorf("queue book %s: %w", bookID, err)
	}
	fmt.Printf("Queued book %s for EPUB reprocessing.\n", bookID)
	return nil
}
