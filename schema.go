package conform

import (
	"fmt"
	"reflect"
)

type field struct {
	GoName          string
	Path            string
	ConfName        string
	EnvName         string
	DefaultValue    string
	HasDefaultValue bool
	Required        bool
	Type            reflect.Type
	Value           reflect.Value
}

type schema struct {
	fields []field
}

func buildSchema(target any) (*schema, error) {
	v := reflect.ValueOf(target)
	if !v.IsValid() || v.Kind() != reflect.Pointer || v.IsNil() {
		return nil, fmt.Errorf("%w: target must be a non-nil pointer to a struct", ErrInvalidTarget)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: target must point to a struct", ErrInvalidTarget)
	}

	s := &schema{}
	if err := collectFields(v, v.Type(), "", &s.fields); err != nil {
		return nil, err
	}
	return s, nil
}

func collectFields(v reflect.Value, t reflect.Type, parentPath string, fields *[]field) error {
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		fieldValue := v.Field(i)

		if structField.PkgPath != "" {
			continue
		}

		path := structField.Name
		if parentPath != "" {
			path = parentPath + "." + structField.Name
		}

		confName := structField.Tag.Get("conf")
		envName := structField.Tag.Get("env")
		defaultValue := structField.Tag.Get("default")
		hasDefaultValue := defaultValue != ""
		required := structField.Tag.Get("required") == "true"

		*fields = append(*fields, field{
			GoName:          structField.Name,
			Path:            path,
			ConfName:        confName,
			EnvName:         envName,
			DefaultValue:    defaultValue,
			HasDefaultValue: hasDefaultValue,
			Required:        required,
			Type:            structField.Type,
			Value:           fieldValue,
		})

		if structField.Type.Kind() == reflect.Struct {
			if err := collectFields(fieldValue, structField.Type, path, fields); err != nil {
				return err
			}
		}
	}
	return nil
}

func isZeroValue(v reflect.Value) bool {
	return v.IsZero()
}
