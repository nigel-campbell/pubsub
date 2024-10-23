package main

import (
	"context"
	"database/sql"
	"os"
	"testing"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "pubsub.db")
	ok(t, err)
	defer db.Close()

	s := NewService(db)
	err = s.Init()
	ok(t, err)

	err = s.CreateTopic(ctx, "topic1", []byte("metadata"))
	ok(t, err)

	topic, err := s.GetTopic(ctx, "topic1")
	ok(t, err)
	equals(t, "topic1", topic.Name)

	_ = os.Remove("pubsub.db")
	t.Log("Test successful. Database removed")
}

func equals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func ok(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
