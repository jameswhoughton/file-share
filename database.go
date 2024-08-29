package file_share

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/jameswhoughton/migrate"
	"github.com/jameswhoughton/migrate/pkg/migrationLog"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(conn *sql.DB) error {
	migrationFS := fs.FS(migrationFiles)

	migrationLog, err := migrationLog.NewLogSQLite(conn)

	if err != nil {
		return err
	}

	migrationsFiles, err := fs.Sub(migrationFS, "migrations")

	if err != nil {
		return err
	}

	err = migrate.Migrate(conn, migrationsFiles, &migrationLog)

	if err != nil {
		return err
	}

	return nil
}
