package validators

import (
	"fmt"
	"regexp"

	"github.com/nzhussup/konform/internal/errs"
	schematypes "github.com/nzhussup/konform/internal/schema/types"
	"github.com/nzhussup/konform/internal/validate/types"
)

const RegexRuleName = "regex"

func Regex(f schematypes.Field, validations *[]types.ValidationResult) {
	regexArg, ok := f.ValidationArg(RegexRuleName)
	if !ok {
		return
	}

	actualString := f.ValueAsString()
	ok, err := regexp.MatchString(regexArg, actualString)
	if err != nil {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: invalid regex pattern %q", errs.InvalidSchema, regexArg),
		})
		return
	}
	if !ok {
		*validations = append(*validations, types.ValidationResult{
			Field: f,
			Err:   fmt.Errorf("%w: value %q does not match regex %q", errs.ValidationRegex, actualString, regexArg),
		})
	}
}
