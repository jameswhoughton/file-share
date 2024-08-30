package web

import (
	"embed"
	"errors"
	"log"
	"net/http"

	file_share "github.com/jameswhoughton/file-share"
)

type middleware func(http.Handler) http.Handler

//go:embed templates/*.gohtml
var templateFiles embed.FS

//go:embed static/*
var publicFiles embed.FS

func AddRoutes(
	mux *http.ServeMux,
	userService file_share.UserService,
	sessionService file_share.SessionService,
) {
	// Middleware
	isAuthed := newIsAuthenticatedMiddleware(sessionService)

	// Root redirect
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session")

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

		http.Redirect(w, r, "/account", http.StatusFound)

	})

	// Static files
	mux.Handle("GET /static/", getStaticFilesHandler(publicFiles))

	// Authentication
	mux.Handle("GET /login", getLoginHandler(templateFiles))
	mux.Handle("POST /login", postLoginHandler(userService, sessionService))
	mux.Handle("GET /register", getRegistrationHandler(templateFiles))
	mux.Handle("POST /register", postRegistrationHandler(userService))
	mux.Handle("GET /logout", getLogoutHandler(sessionService))

	// Account
	mux.Handle("GET /account", isAuthed(getAccountHandler(templateFiles, userService)))
}
