package storage

import (
	"testing"
)

func TestURLToPath(t *testing.T) {
	testCases := []struct {
		name     string
		rawURL   string
		expected string
		hasError bool
	}{
		{
			name:     "Root URL",
			rawURL:   "https://example.com",
			expected: "example.com/index.html",
			hasError: false,
		},
		{
			name:     "Root URL with slash",
			rawURL:   "https://example.com/",
			expected: "example.com/index.html",
			hasError: false,
		},
		{
			name:     "Path without extension",
			rawURL:   "https://example.com/about",
			expected: "example.com/about/index.html",
			hasError: false,
		},
		{
			name:     "Path with trailing slash",
			rawURL:   "https://example.com/news/",
			expected: "example.com/news/index.html",
			hasError: false,
		},
		{
			name:     "Path to a file",
			rawURL:   "https://example.com/style.css",
			expected: "example.com/style.css",
			hasError: false,
		},
		{
			name:     "Path to an HTML file",
			rawURL:   "https://example.com/articles/post.html",
			expected: "example.com/articles/post.html",
			hasError: false,
		},
		{
			name:     "Invalid URL",
			rawURL:   "://invalid-url",
			expected: "",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := URLToPath(tc.rawURL)

			if tc.hasError {
				if err == nil {
					t.Errorf("Expected an error, but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Did not expect an error, but got: %v", err)
			}

			if actual != tc.expected {
				t.Errorf("Expected path %q, but got %q", tc.expected, actual)
			}
		})
	}
}
