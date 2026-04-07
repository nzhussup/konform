package conform

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidTarget = errors.New("conform: invalid target")
	ErrInvalidSchema = errors.New("conform: invalid schema")
	ErrDecode        = errors.New("conform: decode error")
	ErrValidation    = errors.New("conform: validation failed")
)

type FieldError struct {
	Path string
	Err  error
}

func (e FieldError) Error() string {
	if e.Path == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Path, e.Err)
}

type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	if e == nil || len(e.Fields) == 0 {
		return ErrValidation.Error()
	}

	var b strings.Builder
	b.WriteString("conform: validation failed:")
	for _, field := range e.Fields {
		b.WriteString("\n  - ")
		b.WriteString(field.Error())
	}
	return b.String()
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}
