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

func TestLen(t *testing.T) {
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

	tests := []struct {
		name        string
		field       schema.Field
		initial     []types.ValidationResult
		wantTotal   int
		wantErrType error
		wantLike    string
	}{
		{
			name:      "missing len rule is ignored",
			field:     makeStringField("Name", "abcd", nil),
			wantTotal: 0,
		},
		{
			name:      "equal length passes",
			field:     makeStringField("Name", "abcd", map[string]string{"len": "4"}),
			wantTotal: 0,
		},
		{
			name:        "mismatched length returns validation len",
			field:       makeStringField("Name", "abcd", map[string]string{"len": "3"}),
			wantTotal:   1,
			wantErrType: errs.ValidationLen,
			wantLike:    "expected length 3 but got 4",
		},
		{
			name:        "invalid len arg returns invalid schema",
			field:       makeStringField("Name", "abcd", map[string]string{"len": "abc"}),
			wantTotal:   1,
			wantErrType: errs.InvalidSchema,
			wantLike:    `invalid len value "abc"`,
		},
		{
			name:      "non string value uses fmt string length",
			field:     makeIntField("Port", 1234, map[string]string{"len": "4"}),
			wantTotal: 0,
		},
		{
			name:  "appends after existing validations",
			field: makeStringField("Name", "abcd", map[string]string{"len": "3"}),
			initial: []types.ValidationResult{
				{Err: errors.New("existing")},
			},
			wantTotal:   2,
			wantErrType: errs.ValidationLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := append([]types.ValidationResult(nil), tt.initial...)

			validators.Len(tt.field, &results)

			if got := len(results); got != tt.wantTotal {
				t.Fatalf("len(results) = %d, want %d", got, tt.wantTotal)
			}
			if tt.wantErrType == nil {
				return
			}

			got := results[len(results)-1]
			if got.Field.Path != tt.field.Path {
				t.Fatalf("result field Path = %q, want %q", got.Field.Path, tt.field.Path)
			}
			if !errors.Is(got.Err, tt.wantErrType) {
				t.Fatalf("result error = %v, want wrapped %v", got.Err, tt.wantErrType)
			}
			if tt.wantLike != "" && !strings.Contains(got.Err.Error(), tt.wantLike) {
				t.Fatalf("result error = %q, want to contain %q", got.Err.Error(), tt.wantLike)
			}
		})
	}
}
