package template

import (
	"net/http"

	"github.com/efreitasn/gtodo/internal/data/template"
	"github.com/efreitasn/gtodo/internal/data/user"
	"github.com/efreitasn/gtodo/pkg/flash"
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

// PushAssets uses HTTP/2 server push to send the static assets that will be needed in all html pages.
func (t *Template) PushAssets(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			files := []string{
				"/static/css/style.css",
				"/static/fonts/roboto-bold.woff2",
				"/static/fonts/roboto.woff2",
			}

			for _, file := range files {
				err := pusher.Push(file, nil)

				if err != nil {
					break
				}
			}
		}

		next(w, r)
	}
}
