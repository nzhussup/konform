package validators

import (
	"fmt"
	"strconv"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const MinLenRuleName = "minlen"

func MinLen(f schematypes.Field, validations *[]types.ValidationResult) {
	minLenArg, ok := f.ValidationArg(MinLenRuleName)
	if !ok {
		return
	}

	expectedMinLen, err := strconv.Atoi(minLenArg)
	if err != nil {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: invalid %s value %q", errs.InvalidSchema, MinLenRuleName, minLenArg),
		})
		return
	}

	actualLen := len(f.ValueAsString())
	if actualLen < expectedMinLen {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: expected minimum length %d but got %d", errs.ValidationMinLen, expectedMinLen, actualLen),
		})
	}
}
