package validators

import (
	"fmt"
	"slices"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const OneOfRuleName = "oneof"

func OneOf(f schematypes.Field, validations *[]types.ValidationResult) {
	oneOfArg, ok := f.ValidationArg(OneOfRuleName)
	if !ok {
		return
	}

	options := parsePipeSeparated(oneOfArg)
	value := f.ValueAsString()

	if slices.Contains(options, value) {
		return
	}

	*validations = append(*validations, types.ValidationResult{
		Field: f,
		Err:   fmt.Errorf("%w: expected one of %v but got %q. check if you have pipe separated values", errs.ValidationOneOf, options, value),
	})
}
