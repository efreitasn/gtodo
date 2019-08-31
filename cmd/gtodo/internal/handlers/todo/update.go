package todo

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/gtodo/internal/data/template"
	"github.com/efreitasn/gtodo/internal/data/user"
	"github.com/efreitasn/gtodo/pkg/flash"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var updatePOSTSuccessMsg = &flash.Message{
	Content: "Todos updated!",
	Kind:    0,
}

var updatePOSTErrorMsg = &flash.Message{
	Content: "Error while updating todos.",
	Kind:    1,
}

// UpdateGET renders the todos to be updated.
func (t *Todo) UpdateGET(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
	tData.Mode = "update"

	t.fetchDoneNotDone(w, r, tData)
}

// UpdatePOST updates todos.
func (t *Todo) UpdatePOST(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r.ParseForm()

	doneTodosIDs, ok := r.Form["done"]
	userPayload := user.PayloadFromContext(r.Context())
	filterUserID := bson.E{
		Key:   "_user",
		Value: userPayload.ID,
	}

	if ok {
		// Filter values
		filterValue := make(bson.A, len(doneTodosIDs))

		for i, doneTodosID := range doneTodosIDs {
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

		// Update to done = true
		filterOr := bson.D{
			filterUserID,
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
			flash.Add("/update", w, r, updatePOSTErrorMsg)

			return
		}

		// Update to done = false
		filterNor := bson.D{
			filterUserID,
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
			flash.Add("/update", w, r, updatePOSTErrorMsg)

			return
		}
	} else {
		_, err := t.c.UpdateMany(
			ctx,
			bson.D{filterUserID},
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
			flash.Add("/update", w, r, updatePOSTErrorMsg)

			return
		}
	}

	flash.Add("/update", w, r, updatePOSTSuccessMsg)
}
