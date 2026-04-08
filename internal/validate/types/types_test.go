package types

import (
	"errors"
	"reflect"
	"testing"

	schematypes "github.com/nzhussup/konform/internal/schema/types"
)

func TestValidationTypesShape(t *testing.T) {
	var called bool
	fn := ValidationFunc(func(field schematypes.Field, results *[]ValidationResult) {
		called = true
		*results = append(*results, ValidationResult{
			Field: field,
			Err:   errors.New("boom"),
		})
	})

	field := schematypes.Field{Path: "App.Name"}
	results := []ValidationResult{}
	fn(field, &results)

	if !called {
		t.Fatalf("ValidationFunc was not called")
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if !reflect.DeepEqual(results[0].Field, field) {
		t.Fatalf("result field = %#v, want %#v", results[0].Field, field)
	}
	if results[0].Err == nil || results[0].Err.Error() != "boom" {
		t.Fatalf("result err = %v, want boom", results[0].Err)
	}
}
