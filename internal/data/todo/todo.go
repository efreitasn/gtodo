package todo

import "time"

// Todo is an item of the list.
type Todo struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}
