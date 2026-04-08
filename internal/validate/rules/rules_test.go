package rules_test

import (
	"reflect"
	"testing"

	"github.com/nzhussup/konform/internal/schema"
	"github.com/nzhussup/konform/internal/validate/rules"
	"github.com/nzhussup/konform/internal/validate/types"
	"github.com/nzhussup/konform/internal/validate/validators"
)

func TestIsSupported(t *testing.T) {
	tests := []struct {
		name string
		rule string
		want bool
	}{
		{
			name: "required is supported",
			rule: validators.RequiredRuleName,
			want: true,
		},
		{
			name: "min is supported",
			rule: validators.MinRuleName,
			want: true,
		},
		{
			name: "max is supported",
			rule: validators.MaxRuleName,
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
	validator, ok := rules.Registry[validators.RequiredRuleName]
	if !ok {
		t.Fatalf("Registry missing %q rule", validators.RequiredRuleName)
	}

	zero := ""
	field := schema.Field{
		GoName:      "Name",
		Path:        "Name",
		Validations: map[string]string{validators.RequiredRuleName: ""},
		Type:        reflect.TypeOf(""),
		Value:       reflect.ValueOf(&zero).Elem(),
	}

	var results []types.ValidationResult
	validator(field, &results)

	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
}

func TestRegistryMinValidator(t *testing.T) {
	validator, ok := rules.Registry[validators.MinRuleName]
	if !ok {
		t.Fatalf("Registry missing %q rule", validators.MinRuleName)
	}

	v := 9
	field := schema.Field{
		GoName:      "Port",
		Path:        "Port",
		Validations: map[string]string{validators.MinRuleName: "10"},
		Type:        reflect.TypeOf(0),
		Value:       reflect.ValueOf(&v).Elem(),
	}

	var results []types.ValidationResult
	validator(field, &results)

	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
}

func TestRegistryMaxValidator(t *testing.T) {
	validator, ok := rules.Registry[validators.MaxRuleName]
	if !ok {
		t.Fatalf("Registry missing %q rule", validators.MaxRuleName)
	}

	v := 11
	field := schema.Field{
		GoName:      "Port",
		Path:        "Port",
		Validations: map[string]string{validators.MaxRuleName: "10"},
		Type:        reflect.TypeOf(0),
		Value:       reflect.ValueOf(&v).Elem(),
	}

	var results []types.ValidationResult
	validator(field, &results)

	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
}
