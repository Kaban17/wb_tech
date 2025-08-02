package main

import (
	"testing"
)

// TestUnPack проверяет корректность работы функции unPack
func TestUnPack(t *testing.T) {
	// Определяем структуру для тестовых случаев
	type testCase struct {
		name        string
		input       string
		expected    string
		expectError bool
	}

	// Создаем набор тестовых случаев
	testCases := []testCase{
		{
			name:        "basic case with digits",
			input:       "a4bc2d5e",
			expected:    "aaaabccddddde",
			expectError: false,
		},
		{
			name:        "no digits",
			input:       "abcd",
			expected:    "abcd",
			expectError: false,
		},
		{
			name:        "invalid start with digit",
			input:       "45",
			expected:    "invalid string: starts with a digit",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "escaped digits",
			input:       "qwe\\4\\5",
			expected:    "qwe45",
			expectError: false,
		},
		{
			name:        "one escaped, one repeated",
			input:       "qwe\\45",
			expected:    "qwe44444",
			expectError: false,
		},
		{
			name:        "invalid escape sequence at end",
			input:       "a\\",
			expected:    "invalid escape sequence at the end of the string",
			expectError: true,
		},
		{
			name:        "multiple repetitions",
			input:       "ab4c2",
			expected:    "abbbbcc",
			expectError: false,
		},
		{
			name:        "multiple single repetitions",
			input:       "a2b3c",
			expected:    "aabbbc",
			expectError: false,
		},
		{
			name:        "escaped backslash",
			input:       "a\\\\3",
			expected:    "a\\\\\\",
			expectError: false,
		},
		{
			name:        "escaped digit followed by repetition",
			input:       "ab\\2c3",
			expected:    "ab2ccc",
			expectError: false,
		},
	}

	// Перебираем тестовые случаи и выполняем проверку
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			unpacked, err := unPack(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("For input %q, expected an error, but got none", tc.input)
				}
				if err != nil && err.Error() != tc.expected {
					t.Errorf("For input %q, expected error %q, but got %q", tc.input, tc.expected, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("For input %q, expected no error, but got %v", tc.input, err)
				}
				if unpacked != tc.expected {
					t.Errorf("For input %q, expected %q, but got %q", tc.input, tc.expected, unpacked)
				}
			}
		})
	}
}
