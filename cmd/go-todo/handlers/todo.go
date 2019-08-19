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

	cursor, err := t.c.Find(
		ctx,
		bson.D{},
	)

	if err != nil {
		utils.WriteTemplates(w, "Error while fetching the list of todos.", "error")

		return
	}

	var todos []todo.Todo

	cursor.All(
		context.Background(),
		&todos,
	)

	templateData := struct {
		Todos        []todo.Todo
		FlashMessage string
	}{
		todos,
		utils.ReadFlashMessage(w, r),
	}

	utils.WriteTemplates(w, templateData, "todos")
}

// Insert adds a todo.
func (t *Todo) Insert(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.ParseForm()

	todoToBeInserted := todo.InsertTodo{
		Title:     r.Form.Get("title"),
		CreatedAt: time.Now(),
	}

	_, err := t.c.InsertOne(
		ctx,
		todoToBeInserted,
	)

	if err != nil {
		utils.WriteTemplates(w, "Error while adding the todo.", "error")

		return
	}

	utils.AddFlashMessage(w, "Todo added!")

	w.Header().Set("Location", "/todos")
	w.WriteHeader(303)
}
