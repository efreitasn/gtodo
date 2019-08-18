package handlers

import (
	"net/http"
	"time"

	"github.com/efreitasn/go-todo/internal/data/todo"
	"github.com/efreitasn/go-todo/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

var todos []todo.Todo = []todo.Todo{
	todo.Todo{
		ID:        "10",
		Title:     "First",
		CreatedAt: time.Date(2015, 3, 20, 20, 30, 0, 0, time.UTC),
	},
}

// Todo is the representation of all the todo-related handlers.
type Todo struct {
	db *mongo.Database
}

// List list all todos.
func (t *Todo) List(w http.ResponseWriter, r *http.Request) {
	utils.WriteTemplates(w, todos, "todos")
}
