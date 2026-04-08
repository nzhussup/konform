package validators

import (
	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

func Required(f schematypes.Field, validations *[]types.ValidationResult) {
	if schematypes.IsZeroValue(f.Value) {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   errs.ValidationRequired,
		})
	}
}
