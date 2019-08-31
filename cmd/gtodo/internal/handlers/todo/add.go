package todo

import (
	"context"
	"net/http"
	"time"

	"github.com/efreitasn/gtodo/internal/data/template"
	"github.com/efreitasn/gtodo/internal/data/todo"
	"github.com/efreitasn/gtodo/internal/data/user"
	"github.com/efreitasn/gtodo/internal/utils"
	"github.com/efreitasn/gtodo/pkg/flash"
)

var addPOSTErrorMsg = &flash.Message{
	Kind:    1,
	Content: "Error while adding the todo.",
}

var addPOSTSuccessMsg = &flash.Message{
	Kind:    0,
	Content: "Todo added!",
}

// AddGET renders the form to add a todo.
func (t *Todo) AddGET(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
	tData.Mode = "add"
	tData.Title = "Add"

	utils.WriteTemplates(w, tData, "add")
}

// AddPOST adds a todo the to the db.
func (t *Todo) AddPOST(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r.ParseForm()

	userPayload := user.PayloadFromContext(r.Context())

	todoToBeInserted := todo.InsertTodo{
		Title:     r.Form.Get("title"),
		CreatedAt: time.Now(),
		UserID:    userPayload.ID,
	}

	_, err := t.c.InsertOne(ctx, todoToBeInserted)

	if err != nil {
		flash.Add("/add", w, r, addPOSTErrorMsg)

		return
	}

	flash.Add("/add", w, r, addPOSTSuccessMsg)
}
