package todo

import "time"

type Todo struct {
	ID          TodoID
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
