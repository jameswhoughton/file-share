package web

import (
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	file_share "github.com/jameswhoughton/file-share"
	"golang.org/x/crypto/bcrypt"
)

func getLoginHandler(templateFiles fs.FS) http.Handler {
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

func postLoginHandler(userService file_share.UserService, sessionService file_share.SessionService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		user, err := userService.GetWithCredentials(r.FormValue("email"), string(r.FormValue("password")))

		if err != nil {
			http.Redirect(w, r, "/login?invalid-credentials", http.StatusFound)
			return
		}

		session, err := sessionService.Add(file_share.Session{
			UserId:    user.Id,
			SessionId: file_share.GenerateKey(),
		})

		if err != nil {
			w.WriteHeader(500)
			log.Print(err)
			return
		}

		userSession := http.Cookie{
			Name:     "session",
			Value:    session.SessionId,
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(w, &userSession)

		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})
}

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

		http.Redirect(w, r, "/login?new-user", http.StatusFound)
	})
}

func newIsAuthenticatedMiddleware(sessionService file_share.SessionService) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			session, err := r.Cookie("session")

			if err != nil {
				switch {
				case errors.Is(err, http.ErrNoCookie):
					http.Redirect(w, r, "/login", http.StatusFound)
				default:
					log.Println(err)
					http.Error(w, "server error", http.StatusInternalServerError)
				}
				return
			}

			if !sessionService.IsValid(session.Value) {
				http.Redirect(w, r, "/login", http.StatusFound)
			}

			next.ServeHTTP(w, r)
		})
	}
}
