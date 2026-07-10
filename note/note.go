package note

import "time"

type Note struct {
	ID        int
	Title     string
	Body      string
	CreatedAt time.Time
}
