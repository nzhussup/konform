package model

import "github.com/nzhussup/konform/internal/schema"

type ValidationResult struct {
	Field schema.Field
	Err   error
}
