package main

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	user "github.com/jameswhoughton/file-share/internal"
	"github.com/jameswhoughton/migrate"
	"github.com/jameswhoughton/migrate/pkg/migrationLog"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

//go:embed templates/*.gohtml
var templateFiles embed.FS

//go:embed static/*
var publicFiles embed.FS

func migrateDB() (*sql.DB, error) {
	migrationDir := "migrations"

	conn, err := sql.Open("sqlite3", "file-share.db")

	if err != nil {
		return nil, err
	}

	migrationLog, err := migrationLog.NewLogSQLite(conn)

	if err != nil {
		return nil, err
	}

	err = migrate.Migrate(conn, migrationDir, &migrationLog)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func generateKey() string {
	key := make([]byte, 32)
	rand.Read(key)

	return string(key)
}

func main() {

	conn, err := migrateDB()

	if err != nil {
		log.Fatal(err)
	}

	userModel := user.NewUserModel(conn)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "templates/layout.gohtml", "templates/register.gohtml")

		if err != nil {
			w.Write([]byte("Template error: " + err.Error()))

			return
		}
		tmpl.ExecuteTemplate(w, "layout", nil)
	})

	mux.Handle("GET /login", getLoginHandler())

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		user, err := userModel.GetWithCredentials(r.FormValue("email"), string(r.FormValue("password")))

		fmt.Println(user)

		if err != nil {
			http.Redirect(w, r, "/login?invalid-credentials", http.StatusFound)
			return
		}

		userSession := http.Cookie{
			Name:     "session",
			Value:    strconv.Itoa(user.Id),
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}

		http.SetCookie(w, &userSession)

		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})

	mux.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}

		form := user.Form{
			Email:    r.FormValue("email"),
			Password: string(hash),
			ApiKey:   generateKey(),
		}

		_, err = userModel.Add(form)

		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/login?new-user", http.StatusFound)
	})

	mux.HandleFunc("GET /dashbaord/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "templates/layout.gohtml", "templates/dashboard.gohtml")

		if err != nil {
			w.Write([]byte("Template error: " + err.Error()))

			return
		}
		tmpl.ExecuteTemplate(w, "layout", nil)
	})

	mux.HandleFunc("GET /file/{hash}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)

		// Verify hash, if invalid return 404

		// Return file
	})

	mux.HandleFunc("POST /file-hash", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)

		// Receive hash of file to be uploaded
		// DATA: file hash, api key, file name, lifetime, file size

		// Create new record in the files table

		// Return the id of the file record
	})

	mux.HandleFunc("PUT /file-chunk/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)

		// Receive and store chunk of file

		// return 206 code if still waiting for more data

		// Once all chunks received, use file size and hash to verify file

		// return 200 if all chunks received and verified

		// return 422 if the file is invalid
	})

	var staticFS = fs.FS(publicFiles)
	htmlContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.StripPrefix("/static/", http.FileServer(http.FS(htmlContent)))

	// Serve static files
	mux.Handle("GET /static/", fs)

	fmt.Println("listening on port :8000")

	http.ListenAndServe(":8000", mux)
}
