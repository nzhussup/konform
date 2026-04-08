package types

import (
	"reflect"
	"testing"
)

func TestFieldMethods(t *testing.T) {
	targetValue := "value"
	field := Field{
		DefaultValue: "default",
		Validations: map[string]string{
			"required": "",
			"min":      "3",
		},
		Value: reflect.ValueOf(&targetValue).Elem(),
	}

	if !field.HasDefaultValue() {
		t.Fatalf("HasDefaultValue() = false, want true")
	}
	if !field.HasValidation("required") {
		t.Fatalf("HasValidation(required) = false, want true")
	}
	if field.HasValidation("max") {
		t.Fatalf("HasValidation(max) = true, want false")
	}

	arg, ok := field.ValidationArg("min")
	if !ok || arg != "3" {
		t.Fatalf("ValidationArg(min) = (%q, %v), want (%q, true)", arg, ok, "3")
	}

	_, ok = field.ValidationArg("max")
	if ok {
		t.Fatalf("ValidationArg(max) ok = true, want false")
	}

	if got := field.ValueAsString(); got != "value" {
		t.Fatalf("ValueAsString() = %q, want %q", got, "value")
	}
}

func TestFieldMethodsEmptyValidations(t *testing.T) {
	field := Field{}
	if field.HasValidation("required") {
		t.Fatalf("HasValidation(required) = true, want false")
	}
	if _, ok := field.ValidationArg("required"); ok {
		t.Fatalf("ValidationArg(required) ok = true, want false")
	}
}

func TestIsZeroValue(t *testing.T) {
	var zeroInt int
	if !IsZeroValue(reflect.ValueOf(zeroInt)) {
		t.Fatalf("IsZeroValue(0) = false, want true")
	}

	nonZeroInt := 1
	if IsZeroValue(reflect.ValueOf(nonZeroInt)) {
		t.Fatalf("IsZeroValue(1) = true, want false")
	}
}
