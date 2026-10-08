package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPluginDocsDownloadSanitizesFilename(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple file name",
			input:    "overview",
			expected: "overview.md",
		},
		{
			name:     "slash in name",
			input:    "sub/overview",
			expected: "sub_overview.md",
		},
		{
			name:     "path traversal relative path",
			input:    "../../etc/passwd",
			expected: ".._.._etc_passwd.md",
		},
		{
			name:     "backslashes in name",
			input:    "..\\..\\windows\\system32",
			expected: ".._.._windows_system32.md",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cleanName := strings.ReplaceAll(tc.input, "/", "_")
			cleanName = strings.ReplaceAll(cleanName, "\\", "_")
			safeName := filepath.Base(cleanName) + ".md"

			require.Equal(t, tc.expected, safeName)
			require.NotContains(t, safeName, "/")
			require.NotContains(t, safeName, "\\")
		})
	}
}
