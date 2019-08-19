package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/go-todo/internal/data/todo"
	"github.com/efreitasn/go-todo/internal/utils"
	"github.com/efreitasn/go-todo/pkg/flash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Todo is the representation of all the todo-related handlers.
type Todo struct {
	c *mongo.Collection
}

// List list all todos.
func (t *Todo) List(w http.ResponseWriter, r *http.Request) {
	templateData := struct {
		Todos        []todo.Todo
		FlashMessage *flash.Message
	}{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := t.c.Find(
		ctx,
		bson.D{},
	)

	if err != nil {
		templateData.FlashMessage = &flash.Message{
			Kind:    1,
			Content: "Error while fetching the list of todos.",
		}

		utils.WriteTemplates(w, templateData, "todos")

		return
	}

	var todos []todo.Todo

	err = cursor.All(
		context.Background(),
		&todos,
	)

	if err != nil {
		templateData.FlashMessage = &flash.Message{
			Kind:    1,
			Content: "Error while fetching the list of todos.",
		}

		utils.WriteTemplates(w, templateData, "todos")

		return
	}

	templateData.Todos = todos
	templateData.FlashMessage = flash.Read(w, r)

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
		flash.Add(
			"/todos",
			w,
			&flash.Message{
				Kind:    1,
				Content: "Error while adding the todo.",
			},
		)

		return
	}

	flash.Add(
		"/todos",
		w,
		&flash.Message{
			Kind:    0,
			Content: "Todo added!",
		},
	)
}
