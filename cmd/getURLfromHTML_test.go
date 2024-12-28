package cmd

import (
	"reflect"
	"testing"
)

func TestGetUrlFromHtml(t *testing.T) {
	cases := []struct {
		name       string
		rawBaseUrl string
		htmlBody   string
		expected   []string
		errExists  bool
	}{
		{
			name:       "absolute URL",
			rawBaseUrl: "https://blog.boot.dev",
			htmlBody: `
			<html>
				<body>
					<a href="https://blog.boot.dev">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev"},
		},
		{
			name:       "relative URL",
			rawBaseUrl: "https://blog.boot.dev",
			htmlBody: `
			<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:       "absolute and relative URLs",
			rawBaseUrl: "https://blog.boot.dev",
			htmlBody: `
			<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
				<body>
				<a href="https://blog.boot.dev/path/two">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected: []string{"https://blog.boot.dev/path/one", "https://blog.boot.dev/path/two"},
		},
		{
			name:       "bad HTML",
			rawBaseUrl: "https://blog.boot.dev",
			htmlBody: `
			<html body>
				<a href="path/one">
					<span>Boot.dev></span>
				</a>
			</html body>
			`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:       "invalid href URL",
			rawBaseUrl: "https://blog.boot.dev",
			htmlBody: `
			<html>
				<body>
					<a href=":\\invalidURL">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected:  nil,
			errExists: true,
		},
		{
			name:       "handle invalid base URL",
			rawBaseUrl: `:\\invalidBaseURL`,
			htmlBody: `
			<html>
				<body>
					<a href="/path">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			expected:  nil,
			errExists: true,
		},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetUrlFromHtml(tc.htmlBody, tc.rawBaseUrl)
			if tc.errExists == true {
				if err == nil {
					t.Errorf("Test %v - '%s'\nFAIL: unexpected error: %v\n", i, tc.name, err)
				}
			} else if reflect.DeepEqual(actual, tc.expected) == false {
				t.Errorf("Test %v = '%s'\nFAIL: expected - '%s' actual - '%s'\n", i, tc.name, tc.expected, actual)
			}
		})
	}
}
