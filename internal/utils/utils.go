package utils

import (
	"fmt"
	"html/template"
	"net/http"
)

// WriteTemplates writes the templates related to the provided file names to a ResponseWriter.
func WriteTemplates(w http.ResponseWriter, data interface{}, fileNames ...string) {
	filePaths := make([]string, 1, len(fileNames)+1)
	filePaths[0] = "web/templates/index.html"

	for _, fileName := range fileNames {
		filePaths = append(filePaths, fmt.Sprintf("web/templates/%v.html", fileName))
	}

	templates := template.Must(template.ParseFiles(filePaths...))

	templates.ExecuteTemplate(w, "index", data)
}
