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

func TestRegex(t *testing.T) {
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
			name:      "missing regex rule is ignored",
			field:     makeStringField("Name", "example_db", nil),
			wantTotal: 0,
		},
		{
			name:      "substring regex match passes",
			field:     makeStringField("Name", "example_db", map[string]string{"regex": "db"}),
			wantTotal: 0,
		},
		{
			name:      "anchored exact regex match passes",
			field:     makeStringField("Name", "db", map[string]string{"regex": "^db$"}),
			wantTotal: 0,
		},
		{
			name:        "anchored regex mismatch returns validation regex",
			field:       makeStringField("Name", "example_db", map[string]string{"regex": "^db$"}),
			wantTotal:   1,
			wantErrType: errs.ValidationRegex,
			wantLike:    `value "example_db" does not match regex "^db$"`,
		},
		{
			name:        "invalid regex pattern returns invalid schema",
			field:       makeStringField("Name", "example_db", map[string]string{"regex": "["}),
			wantTotal:   1,
			wantErrType: errs.InvalidSchema,
			wantLike:    `invalid regex pattern "["`,
		},
		{
			name:      "numeric field value is matched via ValueAsString",
			field:     makeIntField("Port", 5432, map[string]string{"regex": `^[0-9]+$`}),
			wantTotal: 0,
		},
		{
			name:  "appends after existing validations",
			field: makeStringField("Name", "example_db", map[string]string{"regex": "^db$"}),
			initial: []types.ValidationResult{
				{Err: errors.New("existing")},
			},
			wantTotal:   2,
			wantErrType: errs.ValidationRegex,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := append([]types.ValidationResult(nil), tt.initial...)

			validators.Regex(tt.field, &results)

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
