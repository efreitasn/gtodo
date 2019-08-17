package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/efreitasn/go-todo/internal/data/todo"
)

var todos []todo.Todo = []todo.Todo{
	todo.Todo{
		ID:        10,
		Title:     "First",
		CreatedAt: time.Date(2015, 3, 20, 20, 30, 0, 0, time.UTC),
	},
}

// ListTodos list all todos.
func ListTodos(w http.ResponseWriter, r *http.Request) {
	jsonEnc := json.NewEncoder(w)

	jsonEnc.Encode(todos)
}
