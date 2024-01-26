package apperrors

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrParse          = errors.New("parse value")
	ErrTypeNotCorrect = errors.New("type not correct")
)
