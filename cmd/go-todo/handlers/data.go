package handlers

import (
	"net/http"

	"github.com/efreitasn/go-todo/internal/data/todo"
	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/pkg/flash"
)

// TemplateData is the data used in the rendered templates.
type TemplateData struct {
	FlashMessage *flash.Message
	Mode         string
	Auth         bool
	TodosDone    []todo.Todo
	TodosNotDone []todo.Todo
}

// SetUpTemplateDataNext is the function expected to be called after SetUpTemplateData is called.
type SetUpTemplateDataNext func(w http.ResponseWriter, r *http.Request, templateData *TemplateData)

// SetUpTemplateData sets up the template data to be rendered, e.g. reading the current flash message.
func SetUpTemplateData(next SetUpTemplateDataNext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tData TemplateData

		payload := user.PayloadFromContext(r.Context())

		if payload != nil {
			tData.Auth = true
		}

		tData.FlashMessage = flash.Read(w, r)

		next(w, r, &tData)
	}
}
