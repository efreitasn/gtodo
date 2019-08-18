package handlers

import (
	"context"
	"net/http"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := t.c.Find(
		ctx,
		bson.D{},
	)

	if err != nil {
		utils.WriteTemplates(w, "Error while fetching the list of todos.", "error")

		return
	}

	var todos []todo.Todo

	cur.All(
		context.Background(),
		&todos,
	)

	utils.WriteTemplates(w, todos, "todos")
}
