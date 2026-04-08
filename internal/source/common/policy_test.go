package common

import "testing"

func TestUnknownKeySuggestionModeConstants(t *testing.T) {
	if UnknownKeySuggestionError != 0 {
		t.Fatalf("UnknownKeySuggestionError = %d, want 0", UnknownKeySuggestionError)
	}
	if UnknownKeySuggestionOff != 1 {
		t.Fatalf("UnknownKeySuggestionOff = %d, want 1", UnknownKeySuggestionOff)
	}
	if UnknownKeySuggestionError == UnknownKeySuggestionOff {
		t.Fatalf("suggestion modes must be distinct")
	}
}
