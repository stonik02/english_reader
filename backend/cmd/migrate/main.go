package main

import (
	"context"
	"flag"
	"log"

	"github.com/deniskrylov/english-reader/backend/internal/config"
	"github.com/deniskrylov/english-reader/backend/internal/database"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	switch *direction {
	case "up":
		err = database.Up(ctx, pool)
	case "down":
		err = database.Down(ctx, pool)
	default:
		log.Fatalf("unsupported migration direction %q", *direction)
	}
	if err != nil {
		log.Fatal(err)
	}
}
