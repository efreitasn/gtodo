package about

import (
	"net/http"

	"github.com/efreitasn/gtodo/internal/data/template"
	"github.com/efreitasn/gtodo/internal/utils"
)

// About renders the about page.
func About(w http.ResponseWriter, r *http.Request) {
	tData := template.DataFromContext(r.Context())
	tData.Mode = "about"
	tData.Title = "About"

	utils.WriteTemplates(w, tData, "about")
}
