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

func TestMaxLen(t *testing.T) {
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
			name:      "missing max_len rule is ignored",
			field:     makeStringField("Name", "abc", nil),
			wantTotal: 0,
		},
		{
			name:      "length equal to max_len passes",
			field:     makeStringField("Name", "abc", map[string]string{"max_len": "3"}),
			wantTotal: 0,
		},
		{
			name:      "length smaller than max_len passes",
			field:     makeStringField("Name", "ab", map[string]string{"max_len": "3"}),
			wantTotal: 0,
		},
		{
			name:        "length larger than max_len returns validation max len",
			field:       makeStringField("Name", "abcd", map[string]string{"max_len": "3"}),
			wantTotal:   1,
			wantErrType: errs.ValidationMaxLen,
			wantLike:    "expected maximum length 3 but got 4",
		},
		{
			name:        "invalid max_len arg returns invalid schema",
			field:       makeStringField("Name", "abc", map[string]string{"max_len": "abc"}),
			wantTotal:   1,
			wantErrType: errs.InvalidSchema,
			wantLike:    `invalid max_len value "abc"`,
		},
		{
			name:  "appends after existing validations",
			field: makeStringField("Name", "abcd", map[string]string{"max_len": "3"}),
			initial: []types.ValidationResult{
				{Err: errors.New("existing")},
			},
			wantTotal:   2,
			wantErrType: errs.ValidationMaxLen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := append([]types.ValidationResult(nil), tt.initial...)

			validators.MaxLen(tt.field, &results)

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
