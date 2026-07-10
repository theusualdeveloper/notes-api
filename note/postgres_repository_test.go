package note_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theusualdeveloper/notes-api/db"
	"github.com/theusualdeveloper/notes-api/note"
)

func TestPostgresNoteRepository(t *testing.T) {
	dsn := "postgres://admin:admin@localhost:5332/notes"
	ctx := context.Background()
	pgxpool, err := db.NewDB(ctx, dsn)
	if err != nil {
		t.Fatalf("connecting to db: %s", err.Error())
	}
	repo := note.NewPostgresNoteRepository(pgxpool)
	createdNote, err := repo.Create(ctx, "test", "test body")
	if err != nil || createdNote.ID == 0 {
		t.Fatalf("creating note failed: %s", err.Error())
	}
	foundNote, err := repo.GetByID(ctx, createdNote.ID)
	if err != nil {
		t.Fatalf("finding note failed: %s", err.Error())
	} else if createdNote.Title != foundNote.Title ||
		createdNote.Body != foundNote.Body {
		t.Fatalf("want note with title: %s and body: %s, got title: %s and body: %s",
			createdNote.Title,
			createdNote.Body,
			foundNote.Title,
			foundNote.Body,
		)
	}
	err = repo.Delete(ctx, createdNote.ID)
	if err != nil {
		t.Fatalf("deleting note failed: %s", err.Error())
	}
	_, err = repo.GetByID(ctx, createdNote.ID)
	if !errors.Is(err, note.ErrNotFound) {
		t.Fatal("error must be ErrNotFound type")
	}
}
