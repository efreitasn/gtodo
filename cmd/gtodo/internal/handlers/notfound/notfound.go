package notfound

import (
	"net/http"

	"github.com/efreitasn/gtodo/internal/data/template"
	"github.com/efreitasn/gtodo/internal/utils"
)

// NotFound renders the 404 page.
func NotFound(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
	tData.Mode = "404"
	tData.Title = "404"

	utils.WriteTemplates(w, tData, "notfound")
}
