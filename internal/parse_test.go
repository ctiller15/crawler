package internal

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "multiple relative and absolute URLs",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">Link One</a>
				<a href="https://example.com/path/two">Link Two</a>
				<a href="/path/three">Link Three</a>
			</body>
		</html>
		`,
			expected: []string{"https://example.com/path/one", "https://example.com/path/two", "https://example.com/path/three"},
		},
		{
			name:     "no anchor tags",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<div>No links here!</div>
				<p>Still no links...</p>
			</body>
		</html>
		`,
			expected: []string{},
		},
		{
			name:     "nested anchor tags",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<div>
					<a href="/nested">Nested Link</a>
					<div>
						<a href="/deeply/nested">Deeply Nested Link</a>
					</div>
				</div>
			</body>
		</html>			
		`,
			expected: []string{
				"https://example.com/nested",
				"https://example.com/deeply/nested",
			},
		},
		{
			name:     "invalid href values",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<a>No href attribute</a>
				<a href="">Empty href</a>
				<a href="#fragment">Fragment only</a>
			</body>
		</html>			
		`,
			expected: []string{},
		},
		{
			name:     "special character handling",
			inputURL: "https://example.com",
			inputBody: `
		<html>
			<body>
				<a href="/path with spaces">Space in URL</a>
				<a href="/path?param=value&another=val">Query Parameters</a>
				<a href="/weird|path">Weird Characters</a>
			</body>
		</html>
		`,
			expected: []string{
				"https://example.com/path%20with%20spaces",
				"https://example.com/path?param=value&another=val",
				"https://example.com/weird%7Cpath",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
