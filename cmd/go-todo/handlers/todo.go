package handlers

import (
	"context"
	"net/http"

	"github.com/efreitasn/go-todo/internal/data/todo"
	"github.com/efreitasn/go-todo/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Todo is the representation of all the todo-related handlers.
type Todo struct {
	c *mongo.Collection
}

// List list all todos.
func (t *Todo) List(w http.ResponseWriter, r *http.Request) {
	cur, _ := t.c.Find(
		context.Background(),
		bson.D{},
	)

	var todos []todo.Todo

	cur.All(
		context.Background(),
		&todos,
	)

	utils.WriteTemplates(w, todos, "todos")
}
