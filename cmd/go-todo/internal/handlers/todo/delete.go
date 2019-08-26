package todo

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/go-todo/internal/data/template"
	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/pkg/flash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var deletePOSTSuccessMsg = &flash.Message{
	Content: "Todos deleted!",
	Kind:    0,
}

var deletePOSTErrorMsg = &flash.Message{
	Content: "Error while deleting todos.",
	Kind:    1,
}

// DeleteGET renders the todos to be deleted.
func (t *Todo) DeleteGET(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
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
		userPayload := user.PayloadFromContext(r.Context())

		filterValue := make(bson.A, len(todosToDeleteIDs))

		for i, doneTodosID := range todosToDeleteIDs {
			oID, err := primitive.ObjectIDFromHex(doneTodosID)

			if err != nil {
				flash.Add("/update", w, r, updatePOSTErrorMsg)

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
					Key:   "_user",
					Value: userPayload.ID,
				},
				{
					Key:   "$or",
					Value: filterValue,
				},
			},
		)

		if err != nil {
			flash.Add("/delete", w, r, deletePOSTErrorMsg)

			return
		}
	}

	flash.Add("/delete", w, r, deletePOSTSuccessMsg)
}
