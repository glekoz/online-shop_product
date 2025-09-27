package models

import "errors"

const (
	UniqueErrCode = "23505"
)

var (
	ErrNotFound        = errors.New("no result found")
	ErrInternal        = errors.New("something goes wrong")
	ErrUniqueViolation = errors.New("already exists")
)
