package data

import (
	"net/http"

	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/pkg/flash"
)

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
