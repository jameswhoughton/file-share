package main

import (
	"errors"
	"log"
	"net/http"
)

type middleware func(http.Handler) http.Handler

func authenticatedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

		//user, err :=

		next.ServeHTTP(w, r)
	})
}
