package apperrors

import "errors"

// Global errors can be used in different layers.
var (
	ErrNotFound       = errors.New("not found")
	ErrParse          = errors.New("parse value")
	ErrTypeNotCorrect = errors.New("type not correct")
)
