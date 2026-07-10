package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("note not found")

type PostgresNoteRepository struct {
	pgxpool *pgxpool.Pool
}

func NewPostgresNoteRepository(pgxpool *pgxpool.Pool) PostgresNoteRepository {
	return PostgresNoteRepository{pgxpool: pgxpool}
}

func (pnr *PostgresNoteRepository) Create(ctx context.Context, title, body string) (Note, error) {
	var note Note
	q := `INSERT INTO notes (title, body) 
	VALUES ($1, $2) 
	RETURNING id,created_at`
	err := pnr.pgxpool.QueryRow(ctx, q, title, body).Scan(
		&note.ID,
		&note.CreatedAt,
	)
	if err != nil {
		return Note{}, fmt.Errorf("scan row failed: %w", err)
	}
	note.Title, note.Body = title, body
	return note, nil
}

func (pnr *PostgresNoteRepository) GetAll(ctx context.Context) ([]Note, error) {
	var notes []Note
	q := `SELECT * FROM notes`
	rows, err := pnr.pgxpool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("get notes list query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var note Note

		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Body,
			&note.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan rows failed: %w", err)
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows failed: %w", err)
	}
	return notes, nil
}

func (pnr *PostgresNoteRepository) GetByID(ctx context.Context, id int) (Note, error) {
	var note Note
	q := `SELECT * FROM notes WHERE id = $1`
	row := pnr.pgxpool.QueryRow(ctx, q, id)
	err := row.Scan(
		&note.ID,
		&note.Title,
		&note.Body,
		&note.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Note{}, ErrNotFound
		}
		return Note{}, fmt.Errorf("scan row failed: %w", err)
	}
	return note, nil
}

func (pnr *PostgresNoteRepository) Delete(ctx context.Context, id int) error {
	q := `DELETE FROM notes WHERE id = $1`
	ct, err := pnr.pgxpool.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("deleting query failed: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
