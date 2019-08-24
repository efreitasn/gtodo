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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Todo is the representation of all the todo-related handlers.
type Todo struct {
	c *mongo.Collection
}

func (t *Todo) fetch(w http.ResponseWriter, r *http.Request, mode string) {
	templateData := struct {
		TodosDone    []todo.Todo
		TodosNotDone []todo.Todo
		HasTodos     bool
		FlashMessage *flash.Message
		Mode         string
	}{
		Mode: mode,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Done
	cursorDone, err := t.c.Find(
		ctx,
		bson.D{
			{
				Key:   "done",
				Value: true,
			},
		},
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
		flashMsgCookie := flash.Read(w, r)

		if flashMsgCookie != nil {
			templateData.FlashMessage = flashMsgCookie
		} else {
			templateData.FlashMessage = &flash.Message{
				Kind:    1,
				Content: "Error while fetching the list of todos.",
			}
		}

		utils.WriteTemplates(w, templateData, "todos")

		return
	}

	var todosDone []todo.Todo

	err = cursorDone.All(
		context.Background(),
		&todosDone,
	)

	if err != nil {
		flashMsgCookie := flash.Read(w, r)

		if flashMsgCookie != nil {
			templateData.FlashMessage = flashMsgCookie
		} else {
			templateData.FlashMessage = &flash.Message{
				Kind:    1,
				Content: "Error while fetching the list of todos.",
			}
		}

		utils.WriteTemplates(w, templateData, "todos")

		return
	}

	// Not done
	cursorNotDone, err := t.c.Find(
		ctx,
		bson.D{
			{
				Key:   "done",
				Value: false,
			},
		},
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
		flashMsgCookie := flash.Read(w, r)

		if flashMsgCookie != nil {
			templateData.FlashMessage = flashMsgCookie
		} else {
			templateData.FlashMessage = &flash.Message{
				Kind:    1,
				Content: "Error while fetching the list of todos.",
			}
		}

		utils.WriteTemplates(w, templateData, "todos")

		return
	}

	var todosNotDone []todo.Todo

	err = cursorNotDone.All(
		context.Background(),
		&todosNotDone,
	)

	if err != nil {
		flashMsgCookie := flash.Read(w, r)

		if flashMsgCookie != nil {
			templateData.FlashMessage = flashMsgCookie
		} else {
			templateData.FlashMessage = &flash.Message{
				Kind:    1,
				Content: "Error while fetching the list of todos.",
			}
		}

		utils.WriteTemplates(w, templateData, "todos")

		return
	}

	templateData.TodosDone = todosDone
	templateData.TodosNotDone = todosNotDone
	templateData.HasTodos = len(todosDone) != 0 || len(todosNotDone) != 0
	templateData.FlashMessage = flash.Read(w, r)

	utils.WriteTemplates(w, templateData, "todos")
}

// List lists all todos.
func (t *Todo) List(w http.ResponseWriter, r *http.Request) {
	t.fetch(w, r, "add/update")
}

// DeleteList lists the todos to delete.
func (t *Todo) DeleteList(w http.ResponseWriter, r *http.Request) {
	t.fetch(w, r, "delete")
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
			"/",
			w,
			r,
			&flash.Message{
				Kind:    1,
				Content: "Error while adding the todo.",
			},
			true,
		)

		return
	}

	flash.Add(
		"/",
		w,
		r,
		&flash.Message{
			Kind:    0,
			Content: "Todo added!",
		},
		true,
	)
}

func updateFlashMessage(w http.ResponseWriter, r *http.Request, success bool) {
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
		"/",
		w,
		r,
		msg,
		true,
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
				updateFlashMessage(w, r, false)

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
			updateFlashMessage(w, r, false)

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
			updateFlashMessage(w, r, false)

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
			updateFlashMessage(w, r, false)

			return
		}
	}

	updateFlashMessage(w, r, true)
}
