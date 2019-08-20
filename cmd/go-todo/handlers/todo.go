package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/go-todo/internal/data/todo"
	"github.com/efreitasn/go-todo/internal/utils"
	"github.com/efreitasn/go-todo/pkg/flash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Add adds a todo.
func (t *Todo) Add(w http.ResponseWriter, r *http.Request) {
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

// Update updates the todos.
func (t *Todo) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r.ParseForm()

	doneTodosIDs, ok := r.Form["done"]

	if ok {
		orO := make([]interface{}, len(doneTodosIDs))

		for i, doneTodosID := range doneTodosIDs {
			oID, _ := primitive.ObjectIDFromHex(doneTodosID)

			orO[i] = bson.D{{"_id", oID}}
		}

		filterOr := bson.D{
			{"$or", bson.A(orO)},
		}

		filterNor := bson.D{
			{"$nor", bson.A(orO)},
		}

		_, err := t.c.UpdateMany(
			ctx,
			filterOr,
			bson.D{
				{
					"$set",
					bson.D{
						{"done", true},
					},
				},
			},
		)

		if err != nil {
			flash.Add(
				"/todos",
				w,
				&flash.Message{
					Content: "Error while updating todos.",
					Kind:    1,
				},
			)

			return
		}

		_, err = t.c.UpdateMany(
			ctx,
			filterNor,
			bson.D{
				{
					"$set",
					bson.D{
						{"done", false},
					},
				},
			},
		)

		if err != nil {
			flash.Add(
				"/todos",
				w,
				&flash.Message{
					Content: "Error while updating todos.",
					Kind:    1,
				},
			)

			return
		}
	} else {
		_, err := t.c.UpdateMany(
			ctx,
			bson.D{},
			bson.D{
				{
					"$set",
					bson.D{
						{"done", false},
					},
				},
			},
		)

		if err != nil {
			flash.Add(
				"/todos",
				w,
				&flash.Message{
					Content: "Error while updating todos.",
					Kind:    1,
				},
			)

			return
		}
	}

	flash.Add(
		"/todos",
		w,
		&flash.Message{
			Content: "Todos updated!",
			Kind:    0,
		},
	)
}
