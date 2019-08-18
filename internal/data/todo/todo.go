package todo

import "time"

// Todo is an item of the list.
type Todo struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt time.Time
}
