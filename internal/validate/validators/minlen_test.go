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

func TestMinLen(t *testing.T) {
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
		initial     []types.ValidationResult
		wantTotal   int
		wantErrType error
		wantLike    string
	}{
		{
			name:      "missing minlen rule is ignored",
			field:     makeStringField("Name", "abc", nil),
			wantTotal: 0,
		},
		{
			name:      "length equal to minlen passes",
			field:     makeStringField("Name", "abc", map[string]string{"minlen": "3"}),
			wantTotal: 0,
		},
		{
			name:      "length greater than minlen passes",
			field:     makeStringField("Name", "abcd", map[string]string{"minlen": "3"}),
			wantTotal: 0,
		},
		{
			name:        "length smaller than minlen returns validation min len",
			field:       makeStringField("Name", "ab", map[string]string{"minlen": "3"}),
			wantTotal:   1,
			wantErrType: errs.ValidationMinLen,
			wantLike:    "expected minimum length 3 but got 2",
		},
		{
			name:        "invalid minlen arg returns invalid schema",
			field:       makeStringField("Name", "abc", map[string]string{"minlen": "abc"}),
			wantTotal:   1,
			wantErrType: errs.InvalidSchema,
			wantLike:    `invalid minlen value "abc"`,
		},
		{
			name:  "appends after existing validations",
			field: makeStringField("Name", "ab", map[string]string{"minlen": "3"}),
			initial: []types.ValidationResult{
				{Err: errors.New("existing")},
			},
			wantTotal:   2,
			wantErrType: errs.ValidationMinLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := append([]types.ValidationResult(nil), tt.initial...)

			validators.MinLen(tt.field, &results)

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
