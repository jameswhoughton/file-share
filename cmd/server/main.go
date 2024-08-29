package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	file_share "github.com/jameswhoughton/file-share"
	"github.com/jameswhoughton/file-share/sqlite"
	"github.com/jameswhoughton/file-share/web"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	conn, err := sql.Open("sqlite3", "file-share.db")

	if err != nil {
		log.Fatal(err)
	}

	err = file_share.Migrate(conn)

	if err != nil {
		log.Fatal(err)
	}

	userService := sqlite.NewUserService(conn)
	sessionService := sqlite.NewSessionService(conn)

	mux := http.NewServeMux()

	web.AddRoutes(mux, &userService, &sessionService)

	fmt.Println("listening on port :8000")

	http.ListenAndServe(":8000", mux)
}
