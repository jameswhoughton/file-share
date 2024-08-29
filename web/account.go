package web

import (
	"html/template"
	"io/fs"
	"net/http"
)

func getAccountHandler(templateFiles fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "templates/layout.gohtml", "templates/sidebar.gohtml", "templates/account.gohtml")

		if err != nil {
			w.Write([]byte("Template error: " + err.Error()))

			return
		}
		tmpl.ExecuteTemplate(w, "layout", nil)
	})
}
