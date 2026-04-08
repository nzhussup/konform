package validators

import (
	"fmt"
	"strconv"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const LenRuleName = "len"

func Len(f schematypes.Field, validations *[]types.ValidationResult) {
	lenArg, ok := f.ValidationArg(LenRuleName)
	if !ok {
		return
	}

	expectedLen, err := strconv.Atoi(lenArg)
	if err != nil {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: invalid len value %q", errs.InvalidSchema, lenArg),
		})
		return
	}

	actualLen := len(f.ValueAsString())
	if actualLen != expectedLen {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: expected length %d but got %d", errs.ValidationLen, expectedLen, actualLen),
		})
	}
}
