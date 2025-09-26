package olog

import "testing"

func TestExtractPackageFromFuncName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple function",
			input:    "package/path.Function",
			expected: "package/path",
		},
		{
			name:     "method with receiver",
			input:    "package/path.(*Type).Method",
			expected: "package/path",
		},
		{
			name:     "nested packages",
			input:    "github.com/user/project/internal/pkg.Function",
			expected: "github.com/user/project/internal/pkg",
		},
		{
			name:     "no dot",
			input:    "nopackage",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "complex method",
			input:    "github.com/pellared/olog.(*Logger).Info",
			expected: "github.com/pellared/olog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractPackageFromFuncName(tt.input)
			if result != tt.expected {
				t.Errorf("extractPackageFromFuncName(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
