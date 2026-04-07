package conform

import "fmt"

func applyDefaults(sc *schema) error {
	if sc == nil {
		return fmt.Errorf("%w: nil schema", ErrInvalidSchema)
	}

	for _, f := range sc.fields {
		if !f.HasDefaultValue {
			continue
		}

		if !isZeroValue(f.Value) {
			continue
		}

		if err := setFieldValue(f, f.DefaultValue); err != nil {
			return fmt.Errorf("%w: invalid default for %s: %w", ErrDecode, f.Path, err)
		}
	}

	return nil
}
