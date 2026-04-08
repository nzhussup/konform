package validate

import (
	"github.com/nzhussup/konform/internal/errs"
	"github.com/nzhussup/konform/internal/schema"
	"github.com/nzhussup/konform/internal/validate/model"
	"github.com/nzhussup/konform/internal/validate/validations"
)

func Validate(sc *schema.Schema) ([]model.ValidationResult, error) {
	if sc == nil {
		return nil, errs.InvalidSchemaNil
	}
	var results []model.ValidationResult
	for _, f := range sc.Fields {
		validations.Required(f, &results)
	}

	return results, nil
}
