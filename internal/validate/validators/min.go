package validators

import (
	"fmt"
	"strconv"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const MinRuleName = "min"

func Min(f schematypes.Field, validations *[]types.ValidationResult) {
	minArg, ok := f.ValidationArg(MinRuleName)
	if !ok {
		return
	}

	minValue, err := strconv.ParseFloat(minArg, 64)
	if err != nil {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: invalid %s value %q", errs.InvalidSchema, MinRuleName, minArg),
		})
		return
	}

	actual, ok := numericAsFloat64(f.Value)
	if !ok {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: %s validation supports only numeric values", errs.ValidationNonNumeric, MinRuleName),
		})
		return
	}

	if actual < minValue {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: value %f is less than minimum %f", errs.ValidationMin, actual, minValue),
		})
	}
}
