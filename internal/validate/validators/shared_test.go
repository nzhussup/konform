package validators

import (
	"reflect"
	"testing"
)

func TestNumericAsFloat64(t *testing.T) {
	tests := []struct {
		name  string
		value reflect.Value
		want  float64
		ok    bool
	}{
		{
			name:  "int kind",
			value: reflect.ValueOf(-7),
			want:  -7,
			ok:    true,
		},
		{
			name:  "uint kind",
			value: reflect.ValueOf(uint(7)),
			want:  7,
			ok:    true,
		},
		{
			name:  "float32 kind",
			value: reflect.ValueOf(float32(1.5)),
			want:  1.5,
			ok:    true,
		},
		{
			name:  "float64 kind",
			value: reflect.ValueOf(2.75),
			want:  2.75,
			ok:    true,
		},
		{
			name:  "non numeric kind",
			value: reflect.ValueOf("abc"),
			want:  0,
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := numericAsFloat64(tt.value)
			if ok != tt.ok {
				t.Fatalf("numericAsFloat64() ok = %v, want %v", ok, tt.ok)
			}
			if got != tt.want {
				t.Fatalf("numericAsFloat64() value = %v, want %v", got, tt.want)
			}
		})
	}
}
