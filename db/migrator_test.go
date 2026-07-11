package db_test

import (
	"context"
	"testing"

	"github.com/theusualdeveloper/notes-api/db"
)

func TestRun(t *testing.T) {
	dsn := "postgres://admin:admin@localhost:5332/notes"
	ctx := context.Background()
	pgxpool, err := db.NewDB(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Run(ctx, pgxpool, "../migrations")
	if err != nil {
		t.Fatal(err)
	}
	var exists bool
	var bmc, amc int
	q := `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables WHERE table_name = 'notes'
	)`
	err = pgxpool.QueryRow(ctx, q).Scan(&exists)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("notes table does not exist")
	}
	q = `SELECT COUNT(*) AS migrationsCount FROM schema_migrations`
	err = pgxpool.QueryRow(ctx, q).Scan(&bmc)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Run(ctx, pgxpool, "../migrations")
	if err != nil {
		t.Fatal(err)
	}
	err = pgxpool.QueryRow(ctx, q).Scan(&amc)
	if err != nil {
		t.Fatal(err)
	}
	if amc != bmc {
		t.Fatal("schema migrations count must be equal")
	}
}
