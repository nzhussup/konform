package validators_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/nzhussup/konform/internal/errs"
	"github.com/nzhussup/konform/internal/schema"
	"github.com/nzhussup/konform/internal/validate/types"
	"github.com/nzhussup/konform/internal/validate/validators"
)

func TestMin(t *testing.T) {
	makeIntField := func(name string, value int, rules map[string]string) schema.Field {
		v := value
		return schema.Field{
			GoName:      name,
			Path:        name,
			Validations: rules,
			Type:        reflect.TypeOf(0),
			Value:       reflect.ValueOf(&v).Elem(),
		}
	}
	makeFloatField := func(name string, value float64, rules map[string]string) schema.Field {
		v := value
		return schema.Field{
			GoName:      name,
			Path:        name,
			Validations: rules,
			Type:        reflect.TypeOf(0.0),
			Value:       reflect.ValueOf(&v).Elem(),
		}
	}
	makeStringField := func(name, value string, rules map[string]string) schema.Field {
		v := value
		return schema.Field{
			GoName:      name,
			Path:        name,
			Validations: rules,
			Type:        reflect.TypeOf(""),
			Value:       reflect.ValueOf(&v).Elem(),
		}
	}

	tests := []struct {
		name        string
		field       schema.Field
		wantCount   int
		wantErrType error
		wantLike    string
	}{
		{
			name:      "missing min rule is ignored",
			field:     makeIntField("Port", 5, nil),
			wantCount: 0,
		},
		{
			name:      "numeric value greater than min passes",
			field:     makeIntField("Port", 12, map[string]string{"min": "10"}),
			wantCount: 0,
		},
		{
			name:      "numeric value equal to min passes",
			field:     makeFloatField("Rate", 1.5, map[string]string{"min": "1.5"}),
			wantCount: 0,
		},
		{
			name:        "numeric value below min returns validation min",
			field:       makeIntField("Port", 9, map[string]string{"min": "10"}),
			wantCount:   1,
			wantErrType: errs.ValidationMin,
		},
		{
			name:        "non numeric field with min returns non numeric validation error",
			field:       makeStringField("Name", "abc", map[string]string{"min": "2"}),
			wantCount:   1,
			wantErrType: errs.ValidationNonNumeric,
			wantLike:    "min validation supports only numeric values",
		},
		{
			name:        "invalid min arg returns invalid schema",
			field:       makeIntField("Port", 10, map[string]string{"min": "ten"}),
			wantCount:   1,
			wantErrType: errs.InvalidSchema,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []types.ValidationResult

			validators.Min(tt.field, &results)

			if got := len(results); got != tt.wantCount {
				t.Fatalf("len(results) = %d, want %d", got, tt.wantCount)
			}
			if tt.wantCount == 0 {
				return
			}
			if !errors.Is(results[0].Err, tt.wantErrType) {
				t.Fatalf("result error = %v, want wrapped %v", results[0].Err, tt.wantErrType)
			}
			if tt.wantLike != "" && !strings.Contains(results[0].Err.Error(), tt.wantLike) {
				t.Fatalf("result error = %q, want to contain %q", results[0].Err.Error(), tt.wantLike)
			}
		})
	}
}
