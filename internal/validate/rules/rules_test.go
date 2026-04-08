package rules_test

import (
	"reflect"
	"testing"

	"github.com/nzhussup/konform/internal/schema"
	"github.com/nzhussup/konform/internal/validate/rules"
	"github.com/nzhussup/konform/internal/validate/types"
)

func TestIsSupported(t *testing.T) {
	tests := []struct {
		name string
		rule string
		want bool
	}{
		{
			name: "required is supported",
			rule: rules.Required,
			want: true,
		},
		{
			name: "unknown rule is not supported",
			rule: "unknown",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rules.IsSupported(tt.rule); got != tt.want {
				t.Fatalf("IsSupported(%q) = %v, want %v", tt.rule, got, tt.want)
			}
		})
	}
}

func TestRegistryRequiredValidator(t *testing.T) {
	validator, ok := rules.Registry[rules.Required]
	if !ok {
		t.Fatalf("Registry missing %q rule", rules.Required)
	}

	zero := ""
	field := schema.Field{
		GoName:      "Name",
		Path:        "Name",
		Validations: map[string]string{rules.Required: ""},
		Type:        reflect.TypeOf(""),
		Value:       reflect.ValueOf(&zero).Elem(),
	}

	var results []types.ValidationResult
	validator(field, &results)

	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
}
