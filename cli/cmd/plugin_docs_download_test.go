package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPluginDocsDownload_PathTraversal(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "normal_doc",
			expected: "normal_doc.md",
		},
		{
			input:    "../../etc/passwd",
			expected: ".._.._etc_passwd.md",
		},
		{
			input:    "..\\..\\windows\\system32\\config",
			expected: ".._.._windows_system32_config.md",
		},
		{
			input:    "nested/folder/doc",
			expected: "nested_folder_doc.md",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := sanitizeDocName(tc.input)
			require.Equal(t, tc.expected, result)
			require.False(t, strings.Contains(result, "/"))
			require.False(t, strings.Contains(result, "\\"))
			require.Equal(t, filepath.Base(result), result)
		})
	}
}
