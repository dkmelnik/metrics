package metrics

import "errors"

var (
	ErrParse          = errors.New("parse value")
	ErrTypeNotCorrect = errors.New("type not correct")
)
