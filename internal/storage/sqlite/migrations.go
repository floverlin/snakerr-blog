package sqlite

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func MustMigration(db *sql.DB, dir string) {
	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}
	if err := goose.Up(db, dir); err != nil {
		panic(err)
	}
}
