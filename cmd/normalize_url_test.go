package cmd

import (
	"testing"
)

func TestNormalizeUrl(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		expected  string
		errExists bool
	}{
		{
			name:     "remove scheme",
			input:    "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove trailing '/'",
			input:    "blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "normalizing case in host/path",
			input:    "BlOg.BoOt.DeV/pAtH",
			expected: "blog.boot.dev/path",
		},
		{
			name:      "handle invalid url",
			input:     "://invalidUrl",
			expected:  "",
			errExists: true,
		},
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NormalizeUrl(tc.input)
			if err != nil && tc.errExists != true {
				t.Errorf("Test %v - '%s'\nFAIL: unexpected error: %v\n", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v = '%s'\nFAIL: expected - '%s' actual - '%s'\n", i, tc.name, tc.expected, actual)
			}
		})
	}
}
