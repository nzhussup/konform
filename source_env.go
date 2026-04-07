package conform

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

type envSource struct{}

func (s envSource) load(sc *schema) error {
	if sc == nil {
		return fmt.Errorf("%w: nil schema", ErrInvalidSchema)
	}

	for _, field := range sc.fields {
		envName := field.EnvName
		if envName == "" {
			continue
		}

		raw, ok := os.LookupEnv(envName)
		if !ok {
			continue
		}

		if err := setFieldValue(field, raw); err != nil {
			return fmt.Errorf("%w: env %q -> %s: %w", ErrDecode, envName, field.Path, err)
		}
	}

	return nil
}

func setFieldValue(field field, raw string) error {
	if !field.Value.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	switch field.Type.Kind() {
	case reflect.String:
		field.Value.SetString(raw)
		return nil

	case reflect.Int:
		v, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("invalid int value %q", raw)
		}
		field.Value.SetInt(int64(v))
		return nil

	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("invalid bool value %q", raw)
		}
		field.Value.SetBool(v)
		return nil

	default:
		return fmt.Errorf("unsupported field type %s", field.Type)
	}
}
