package errs

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestWrapDecode(t *testing.T) {
	cause := fmt.Errorf("%w: expected string, got bool", DecodeTypeMismatch)

	tests := []struct {
		name        string
		kind        error
		context     string
		cause       error
		wantErrType []error
		wantLike    []string
		wantPrefixN int
	}{
		{
			name:    "wraps decode context and strips nested decode prefix",
			kind:    DecodeSourceField,
			context: `json "App.Debug" -> AppDebug`,
			cause:   cause,
			wantErrType: []error{
				Decode, DecodeSourceField, DecodeTypeMismatch,
			},
			wantLike:    []string{`json "App.Debug" -> AppDebug`, "type mismatch", "expected string, got bool"},
			wantPrefixN: 1,
		},
		{
			name:        "wraps without cause",
			kind:        DecodeSourceRead,
			context:     `json file "config.json"`,
			cause:       nil,
			wantErrType: []error{Decode, DecodeSourceRead},
			wantLike:    []string{`json file "config.json"`},
			wantPrefixN: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapDecode(tt.kind, tt.context, tt.cause)
			for _, e := range tt.wantErrType {
				if !errors.Is(err, e) {
					t.Fatalf("WrapDecode() error = %v, want wrapped %v", err, e)
				}
			}
			for _, part := range tt.wantLike {
				if !strings.Contains(err.Error(), part) {
					t.Fatalf("WrapDecode() error = %q, want to contain %q", err.Error(), part)
				}
			}
			if got := strings.Count(err.Error(), Decode.Error()); got != tt.wantPrefixN {
				t.Fatalf("WrapDecode() prefix count = %d, want %d, error=%q", got, tt.wantPrefixN, err.Error())
			}
		})
	}
}

func TestStripDomainPrefix(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		domain error
		want   string
	}{
		{
			name:   "nil error",
			err:    nil,
			domain: Decode,
			want:   "",
		},
		{
			name:   "decode prefix stripped recursively",
			err:    fmt.Errorf("%w: %w", Decode, fmt.Errorf("%w: detail", Decode)),
			domain: Decode,
			want:   "detail",
		},
		{
			name:   "validation prefix stripped",
			err:    ValidationRequired,
			domain: Validation,
			want:   "required",
		},
		{
			name:   "validation min prefix stripped",
			err:    ValidationMin,
			domain: Validation,
			want:   "minimum value not met",
		},
		{
			name:   "validation non numeric prefix stripped",
			err:    ValidationNonNumeric,
			domain: Validation,
			want:   "non-numeric value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripDomainPrefix(tt.err, tt.domain)
			if got != tt.want {
				t.Fatalf("StripDomainPrefix() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripHelpers(t *testing.T) {
	tests := []struct {
		name string
		fn   func(error) string
		in   error
		want string
	}{
		{
			name: "StripDecodePrefix",
			fn:   StripDecodePrefix,
			in:   fmt.Errorf("%w: details", Decode),
			want: "details",
		},
		{
			name: "StripValidationPrefix",
			fn:   StripValidationPrefix,
			in:   fmt.Errorf("%w: details", Validation),
			want: "details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.in)
			if got != tt.want {
				t.Fatalf("%s() = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
