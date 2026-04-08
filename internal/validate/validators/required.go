package validators

import (
	"fmt"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const RequiredRuleName = "required"

func Required(f schematypes.Field, validations *[]types.ValidationResult) {
	if schematypes.IsZeroValue(f.Value) {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: field is required", errs.ValidationRequired),
		})
	}
}
