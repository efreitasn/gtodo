package about

import (
	"net/http"

	"github.com/efreitasn/go-todo/internal/data/template"
	"github.com/efreitasn/go-todo/internal/utils"
)

// About renders the about page.
func About(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
	tData.Mode = "about"

	utils.WriteTemplates(w, tData, "about")
}
