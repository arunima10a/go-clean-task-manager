package usecase

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInternal     = errors.New("internal server error")
)