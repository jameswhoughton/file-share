package main

import (
	"html/template"
	"log"
	"net/http"
)

func getLoginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "templates/layout.gohtml", "templates/form.gohtml", "templates/login.gohtml")

		if err != nil {
			w.Write([]byte("Template error: " + err.Error()))

			return
		}
		type loginData struct {
			Title   string
			Message string
		}
		var message string

		if r.URL.Query().Has("new-user") {
			message = "Your account has been created, please login"
		}

		err = tmpl.ExecuteTemplate(w, "layout", loginData{
			Title:   "Login",
			Message: message,
		})

		if err != nil {
			log.Print(err)
		}
	})
}
