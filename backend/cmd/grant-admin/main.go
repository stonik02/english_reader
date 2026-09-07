package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/deniskrylov/english-reader/backend/internal/config"
	"github.com/deniskrylov/english-reader/backend/internal/database"
	repository "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/auth"
	grantadmin "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/grantadmin"
)

func main() {
	email := flag.String("email", "", "email of the existing user to make an administrator")
	confirm := flag.Bool("confirm", false, "confirm granting the administrator role")
	flag.Parse()

	if *email == "" || !*confirm {
		fmt.Fprintln(os.Stderr, "usage: grant-admin -email user@example.com -confirm")
		os.Exit(2)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	pool, err := database.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	user, err := grantadmin.New(repository.New(pool), repository.New(pool)).Execute(context.Background(), *email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Administrator role granted to %s (%s).\n", user.Email, user.ID)
}
