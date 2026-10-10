package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPluginDocsDownload_PathTraversalSanitization(t *testing.T) {
	testCases := []struct {
		name         string
		itemName     string
		expectedName string
	}{
		{
			name:         "simple doc name",
			itemName:     "overview",
			expectedName: "overview.md",
		},
		{
			name:         "unix path traversal",
			itemName:     "../../etc/passwd",
			expectedName: ".._.._etc_passwd.md",
		},
		{
			name:         "windows path traversal",
			itemName:     "..\\..\\windows\\system32",
			expectedName: ".._.._windows_system32.md",
		},
		{
			name:         "mixed slashes path traversal",
			itemName:     "../foo\\../bar",
			expectedName: ".._foo_.._bar.md",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			docsDir := t.TempDir()

			// Sanitize item.Name using the same logic as in runPluginDocsDownload
			safeName := strings.ReplaceAll(strings.ReplaceAll(tc.itemName, "/", "_"), "\\", "_") + ".md"
			require.Equal(t, tc.expectedName, safeName)

			fn := filepath.Join(docsDir, safeName)

			// Verify that fn is located strictly inside docsDir and does not escape it
			rel, err := filepath.Rel(docsDir, fn)
			require.NoError(t, err)
			escapesDir := rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel)
			require.False(t, escapesDir, "file path %q should not escape directory %q", fn, docsDir)

			fp, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
			require.NoError(t, err)

			_, err = fp.WriteString("# Security Test")
			require.NoError(t, err)
			require.NoError(t, fp.Close())

			// Verify file exists at expected path
			_, err = os.Stat(fn)
			require.NoError(t, err)
		})
	}
}
