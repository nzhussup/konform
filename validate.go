package conform

import "fmt"

func validateRequired(sc *schema) error {
	if sc == nil {
		return fmt.Errorf("%w: nil schema", ErrInvalidSchema)
	}

	var missing []FieldError

	for _, f := range sc.fields {
		if !f.Required {
			continue
		}

		if isZeroValue(f.Value) {
			missing = append(missing, FieldError{
				Path: f.Path,
				Err:  fmt.Errorf("%w: required", ErrValidation),
			})
		}
	}

	if len(missing) > 0 {
		return &ValidationError{
			Fields: missing,
		}
	}

	return nil
}
