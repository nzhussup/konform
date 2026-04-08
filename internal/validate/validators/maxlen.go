package validators

import (
	"fmt"
	"strconv"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const MaxLenRuleName = "maxlen"

func MaxLen(f schematypes.Field, validations *[]types.ValidationResult) {
	maxLenArg, ok := f.ValidationArg(MaxLenRuleName)
	if !ok {
		return
	}

	expectedMaxLen, err := strconv.Atoi(maxLenArg)
	if err != nil {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: invalid %s value %q", errs.InvalidSchema, MaxLenRuleName, maxLenArg),
		})
		return
	}

	actualLen := len(f.ValueAsString())
	if actualLen > expectedMaxLen {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: expected maximum length %d but got %d", errs.ValidationMaxLen, expectedMaxLen, actualLen),
		})
	}
}
