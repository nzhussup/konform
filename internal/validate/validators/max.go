package validators

import (
	"fmt"
	"strconv"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const MaxRuleName = "max"

func Max(f schematypes.Field, validations *[]types.ValidationResult) {
	maxArg, ok := f.ValidationArg(MaxRuleName)
	if !ok {
		return
	}

	maxValue, err := strconv.ParseFloat(maxArg, 64)
	if err != nil {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: invalid %s value %q", errs.InvalidSchema, MaxRuleName, maxArg),
		})
		return
	}

	actual, ok := numericAsFloat64(f.Value)
	if !ok {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: %s validation supports only numeric values", errs.ValidationNonNumeric, MaxRuleName),
		})
		return
	}

	if actual > maxValue {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: value %f is greater than maximum %f", errs.ValidationMax, actual, maxValue),
		})
	}
}
