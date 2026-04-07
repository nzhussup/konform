package conform

import (
	"strings"
	"testing"

	"github.com/nzhussup/conform/internal/errs"
)

func TestValidationErrorFormattingStripsNestedPrefix(t *testing.T) {
	err := &ValidationError{
		Fields: []FieldError{
			{Path: "AppDebug", Err: errs.ValidationRequired},
		},
	}

	msg := err.Error()
	if strings.Count(msg, errs.Validation.Error()) != 1 {
		t.Fatalf("ValidationError() = %q, want exactly one validation prefix", msg)
	}
	if !strings.Contains(msg, "AppDebug: required") {
		t.Fatalf("ValidationError() = %q, want field message %q", msg, "AppDebug: required")
	}
}
