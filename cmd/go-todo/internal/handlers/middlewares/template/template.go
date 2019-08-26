package template

import (
	"net/http"

	"github.com/efreitasn/go-todo/internal/data/template"
	"github.com/efreitasn/go-todo/internal/data/user"
	"github.com/efreitasn/go-todo/pkg/flash"
)

// Template is the representation of all the template-related middlewares.
type Template struct{}

// New creates an Auth struct.
func New() *Template {
	return &Template{}
}

// SetUpTemplateData sets up the template data to be rendered.
func (t *Template) SetUpTemplateData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tData template.Data

		payload := user.PayloadFromContext(r.Context())

		if payload != nil {
			tData.Auth = true
		}

		tData.FlashMessage = flash.Read(w, r)

		newR := r.WithContext(template.ContextWithData(r.Context(), &tData))

		next(w, newR)
	}
}
