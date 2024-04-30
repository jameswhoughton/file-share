package main

import (
	"database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jameswhoughton/migrate"
	"github.com/jameswhoughton/migrate/pkg/migrationLog"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed templates/*.gohtml
var templateFiles embed.FS

func migrateDB() error {
	migrationDir := "migrations"

	migrationLog, err := migrationLog.Init(migrationDir + "/.log")

	if err != nil {
		return err
	}

	conn, err := sql.Open("sqlite3", "file-share.db")

	if err != nil {
		return err
	}

	return migrate.Migrate(conn, migrationDir, migrationLog)
}

func main() {

	err := migrateDB()

	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templateFiles, "templates/layout.gohtml", "templates/register.gohtml")

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

	fmt.Println("listening on port :8000")

	http.ListenAndServe(":8000", mux)
}
