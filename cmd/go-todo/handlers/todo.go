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

/*
 * Messages
 */

// Utils
var fetchDoneNotDoneErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while fetching the list of todos.",
}

// Add
var addPOSTErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while adding the todo.",
}

var addPOSTSuccessMsg = &flash.Message{
	Kind:    0,
	Content: "Todo added!",
}

// Update
var updateGETErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while fetching the list of todos.",
}

var updatePOSTSuccessMsg = &flash.Message{
	Content: "Todos updated!",
	Kind:    0,
}

var updatePOSTErrorMsg = &flash.Message{
	Content: "Error while updating todos.",
	Kind:    1,
}

// Delete
var deletePOSTSuccessMsg = &flash.Message{
	Content: "Todos deleted!",
	Kind:    0,
}

var deletePOSTErrorMsg = &flash.Message{
	Content: "Error while deleting todos.",
	Kind:    1,
}

/*
 * Utils
 */

func (t *Todo) fetchDoneNotDone(w http.ResponseWriter, r *http.Request, tData *TemplateData) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Done
	todosDone, err := todo.FetchDone(ctx, t.c)

	if err != nil {
		if tData.FlashMessage == nil {
			tData.FlashMessage = fetchDoneNotDoneErrorMsg
		}

		utils.WriteTemplates(w, tData, tData.Mode, "no-todos")

		return
	}

	// Not done
	todosNotDone, err := todo.FetchNotDone(ctx, t.c)

	if err != nil {
		if tData.FlashMessage == nil {
			tData.FlashMessage = fetchDoneNotDoneErrorMsg
		}

		utils.WriteTemplates(w, tData, tData.Mode, "no-todos")

		return
	}

	tData.TodosDone = todosDone
	tData.TodosNotDone = todosNotDone

	utils.WriteTemplates(w, tData, tData.Mode, "no-todos")
}

/*
 * List
 */

// ListGET lists the todos.
func (t *Todo) ListGET(w http.ResponseWriter, r *http.Request, tData *TemplateData) {
	tData.Mode = "list"

	t.fetchDoneNotDone(w, r, tData)
}

/*
 * Add
 */

// AddGET renders the form to add a todo.
func (t *Todo) AddGET(w http.ResponseWriter, r *http.Request, tData *TemplateData) {
	tData.Mode = "add"

	utils.WriteTemplates(w, tData, "add")
}

// AddPOST adds a todo the to the db.
func (t *Todo) AddPOST(w http.ResponseWriter, r *http.Request) {
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
			"/add",
			w,
			r,
			addPOSTErrorMsg,
			true,
		)

		return
	}

	flash.Add(
		"/add",
		w,
		r,
		addPOSTSuccessMsg,
		true,
	)
}

/*
 * Update
 */

// UpdateGET renders the todos to be updated.
func (t *Todo) UpdateGET(w http.ResponseWriter, r *http.Request, tData *TemplateData) {
	tData.Mode = "update"

	t.fetchDoneNotDone(w, r, tData)
}

// UpdatePOST updates todos.
func (t *Todo) UpdatePOST(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r.ParseForm()

	doneTodosIDs, ok := r.Form["done"]

	if ok {
		// Filter values
		filterValue := make(bson.A, len(doneTodosIDs))

		for i, doneTodosID := range doneTodosIDs {
			oID, err := primitive.ObjectIDFromHex(doneTodosID)

			if err != nil {
				flash.Add("/update", w, r, updatePOSTErrorMsg, true)

				return
			}

			filterValue[i] = bson.D{
				{
					Key:   "_id",
					Value: oID,
				},
			}
		}

		// Update to done = true
		filterOr := bson.D{
			{
				Key:   "$or",
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
			flash.Add("/update", w, r, updatePOSTErrorMsg, true)

			return
		}

		// Update to done = false
		filterNor := bson.D{
			{
				Key:   "$nor",
				Value: filterValue,
			},
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
			flash.Add("/update", w, r, updatePOSTErrorMsg, true)

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
			flash.Add("/update", w, r, updatePOSTErrorMsg, true)

			return
		}
	}

	flash.Add("/update", w, r, updatePOSTSuccessMsg, true)
}

/*
 * Delete
 */

// DeleteGET renders the todos to be deleted.
func (t *Todo) DeleteGET(w http.ResponseWriter, r *http.Request, tData *TemplateData) {
	tData.Mode = "delete"

	t.fetchDoneNotDone(w, r, tData)
}

// DeletePOST deletes todos from the db.
func (t *Todo) DeletePOST(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r.ParseForm()

	todosToDeleteIDs, ok := r.Form["delete"]

	if ok {
		filterValue := make(bson.A, len(todosToDeleteIDs))

		for i, doneTodosID := range todosToDeleteIDs {
			oID, err := primitive.ObjectIDFromHex(doneTodosID)

			if err != nil {
				flash.Add("/update", w, r, updatePOSTErrorMsg, true)

				return
			}

			filterValue[i] = bson.D{
				{
					Key:   "_id",
					Value: oID,
				},
			}
		}

		_, err := t.c.DeleteMany(
			ctx,
			bson.D{
				{
					Key:   "$or",
					Value: filterValue,
				},
			},
		)

		if err != nil {
			flash.Add("/delete", w, r, deletePOSTErrorMsg, true)

			return
		}
	}

	flash.Add(
		"/delete",
		w,
		r,
		deletePOSTSuccessMsg,
		true,
	)
}
