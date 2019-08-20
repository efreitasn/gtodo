package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/efreitasn/go-todo/internal/data/todo"
	"github.com/efreitasn/go-todo/internal/utils"
	"github.com/efreitasn/go-todo/pkg/flash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		&options.FindOptions{
			Sort: bson.D{
				{
					Key:   "createdAt",
					Value: 1,
				},
			},
		},
	)

	if err != nil {
		fmt.Println(err)
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

func updateFlashMessage(w http.ResponseWriter, success bool) {
	var msg *flash.Message

	if success {
		msg = &flash.Message{
			Content: "Todos updated!",
			Kind:    0,
		}
	} else {
		msg = &flash.Message{
			Content: "Error while updating todos.",
			Kind:    1,
		}
	}

	flash.Add(
		"/todos",
		w,
		msg,
	)
}

// Update updates the todos.
func (t *Todo) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r.ParseForm()

	doneTodosIDs, ok := r.Form["done"]

	if ok {
		// Generate filtering-related values
		filterValue := make(bson.A, len(doneTodosIDs))

		for i, doneTodosID := range doneTodosIDs {
			oID, err := primitive.ObjectIDFromHex(doneTodosID)

			if err != nil {
				updateFlashMessage(w, false)

				return
			}

			filterValue[i] = bson.D{
				{
					Key:   "_id",
					Value: oID,
				},
			}
		}

		filterOr := bson.D{
			{
				Key:   "$or",
				Value: filterValue,
			},
		}

		filterNor := bson.D{
			{
				Key:   "$nor",
				Value: filterValue,
			},
		}

		_, err := t.c.UpdateMany(
			ctx,
			filterOr,
			bson.D{
				{
					Key: "$set",
					Value: bson.D{
						{
							Key:   "done",
							Value: true,
						},
					},
				},
			},
		)

		if err != nil {
			updateFlashMessage(w, false)

			return
		}

		_, err = t.c.UpdateMany(
			ctx,
			filterNor,
			bson.D{
				{
					Key: "$set",
					Value: bson.D{
						{
							Key:   "done",
							Value: false,
						},
					},
				},
			},
		)

		if err != nil {
			updateFlashMessage(w, false)

			return
		}
	} else {
		_, err := t.c.UpdateMany(
			ctx,
			bson.D{},
			bson.D{
				{
					Key: "$set",
					Value: bson.D{
						{
							Key:   "done",
							Value: false,
						},
					},
				},
			},
		)

		if err != nil {
			updateFlashMessage(w, false)

			return
		}
	}

	updateFlashMessage(w, true)
}
