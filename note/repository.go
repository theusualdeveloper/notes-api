package note

import "context"

type NoteRepository interface {
	Create(ctx context.Context, title, body string) (Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	GetById(ctx context.Context) (Note, error)
	Delete(ctx context.Context) error
}
