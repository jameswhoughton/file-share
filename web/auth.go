package web

import (
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	file_share "github.com/jameswhoughton/file-share"
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

		message, err := getMessage(w, r, "message")

		if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
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
			setMessage(w, "message", "credentials are invalid")
			http.Redirect(w, r, "/login", http.StatusFound)
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

		http.Redirect(w, r, "/account", http.StatusFound)
	})
}

func getLogoutHandler(sessionService file_share.SessionService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentSesion, err := r.Cookie("session")

		if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		err = sessionService.Destroy(currentSesion.Value)

		if err != nil {
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}

		userSession := http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(w, &userSession)

		http.Redirect(w, r, "/login", http.StatusFound)
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
