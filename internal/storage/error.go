package storage

import (
	"errors"
)

var (
	ErrNoRows  = errors.New("no result")
	ErrUnique  = errors.New("filed must be unique")
	ErrNotNull = errors.New("field must be not null")
)
