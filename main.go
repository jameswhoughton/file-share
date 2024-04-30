package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/jameswhoughton/migrate"
	"github.com/jameswhoughton/migrate/pkg/migrationLog"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	migrationDir := "migrations"

	migrationLog, err := migrationLog.Init(migrationDir + "/.log")

	if err != nil {
		log.Fatalln(err)
	}

	conn, err := sql.Open("sqlite3", "file-share.db")

	if err != nil {
		log.Fatalln(err)
	}

	err = migrate.Migrate(conn, migrationDir, migrationLog)

	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()

	http.HandleFunc("GET /file", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	http.HandleFunc("POST /file", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	fmt.Println("listening on port :8000")

	http.ListenAndServe(":8000", mux)
}
