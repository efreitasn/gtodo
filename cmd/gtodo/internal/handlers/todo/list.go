package todo

import (
	"net/http"

	"github.com/efreitasn/gtodo/internal/data/template"
)

// ListGET lists the todos.
func (t *Todo) ListGET(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
	tData.Mode = "list"
	tData.Title = "List"

	t.fetchDoneNotDone(w, r, tData)
}
