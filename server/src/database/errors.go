package database

import "errors"

// errors for the database to use
var (
	ErrKeyNotFound = errors.New("key not found")
)