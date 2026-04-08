package rules

import (
	"github.com/nzhussup/konform/internal/validate/types"
	"github.com/nzhussup/konform/internal/validate/validators"
)

const Required = "required"

var Registry = map[string]types.ValidationFunc{
	Required: validators.Required,
}

func IsSupported(name string) bool {
	_, ok := Registry[name]
	return ok
}
