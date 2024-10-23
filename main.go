package main

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run(ctx context.Context) error {
	db, err := sql.Open("sqlite3", "pubsub.db")
	if err != nil {
		return err
	}
	defer db.Close()

	s := NewService(db)

	if err := s.Init(); err != nil {
		return err
	}
	return nil
}
