package sqlite

import (
	"blog/internal/storage"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

func wrapNoRows(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNoRows
	} else {
		return err
	}
}

func wrapUnique(err error) error {
	var sqlErr sqlite3.Error
	if errors.As(err, &sqlErr) &&
		sqlErr.Code == sqlite3.ErrConstraint &&
		sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return storage.ErrUnique
	} else {
		return err
	}
}

func wrapNotNull(err error) error {
	var sqlErr sqlite3.Error
	if errors.As(err, &sqlErr) &&
		sqlErr.Code == sqlite3.ErrConstraint &&
		sqlErr.ExtendedCode == sqlite3.ErrConstraintNotNull {
		return storage.ErrNotNull
	} else {
		return err
	}
}
func noWrap(err error) error {
	return err
}
