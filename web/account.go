package web

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	file_share "github.com/jameswhoughton/file-share"
	"golang.org/x/crypto/bcrypt"
)

func getRegistrationHandler(templateFiles fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "templates/layout.gohtml", "templates/form.gohtml", "templates/register.gohtml")

		if err != nil {
			w.Write([]byte("Template error: " + err.Error()))

			return
		}
		tmpl.ExecuteTemplate(w, "layout", nil)
	})
}

func postRegistrationHandler(userService file_share.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}

		user := file_share.User{
			Email:    r.FormValue("email"),
			Password: string(hash),
			ApiKey:   file_share.GenerateKey(),
		}

		_, err = userService.Add(user)

		if err != nil {
			log.Fatal(err)
		}

		setMessage(w, "message", "Your account has been created, please login below")

		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func getAccountHandler(templateFiles fs.FS, userService file_share.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "templates/layout.gohtml", "templates/sidebar.gohtml", "templates/account.gohtml")

		if err != nil {
			w.Write([]byte("Template error: " + err.Error()))

			return
		}

		session, err := r.Cookie("session")

		if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		user, err := userService.GetFromSessionId(session.Value)

		if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		type templateData struct {
			Title  string
			ApiKey string
		}

		err = tmpl.ExecuteTemplate(w, "layout", templateData{
			Title:  "My Account",
			ApiKey: user.ApiKey,
		})

		if err != nil {
			log.Println(err)
		}
	})
}

func putAccountHandler(userService file_share.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// if r.FormValue("password") != "" && r.FormValue("password") != r.FormValue("passwordConfirm") {

		// 	http.Redirect()
		// }

		// hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		// if err != nil {
		// 	log.Println(err)
		// }

		// user := file_share.User{
		// 	Email:    r.FormValue("email"),
		// 	Password: string(hash),
		// 	ApiKey:   file_share.GenerateKey(),
		// }
	})
}
