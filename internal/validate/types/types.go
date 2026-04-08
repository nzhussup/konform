package types

import (
	schematypes "github.com/nzhussup/konform/internal/schema/types"
)

type ValidationResult struct {
	Field schematypes.Field
	Err   error
}

type ValidationFunc func(field schematypes.Field, results *[]ValidationResult)
