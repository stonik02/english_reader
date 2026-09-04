package wordnormalizer

import "testing"

func TestServiceNormalize(t *testing.T) {
	service := New(32)
	tests := []struct {
		input, want string
		valid       bool
	}{{" Hello! ", "hello", true}, {"went", "went", true}, {"<b>word</b>", "", false}, {"two words", "", false}}
	for _, test := range tests {
		got, err := service.Normalize(test.input)
		if (err == nil) != test.valid || got != test.want {
			t.Fatalf("Normalize(%q) = %q, %v", test.input, got, err)
		}
	}
}
