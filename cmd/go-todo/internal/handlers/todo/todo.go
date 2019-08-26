package todo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Todo is the representation of all the todo-related handlers.
type Todo struct {
	c *mongo.Collection
}

// New creates a Todo struct.
func New(c *mongo.Collection) *Todo {
	return &Todo{c}
}
