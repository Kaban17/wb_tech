package repository

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidID    = errors.New("invalid id")
	ErrDoesNotExist = errors.New("does not exist")
)
