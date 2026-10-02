package cmd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPluginDocsDownload_PathTraversal(t *testing.T) {
	docsDir := t.TempDir()

	tests := []struct {
		name             string
		itemName         string
		expectedSafeName string
	}{
		{
			name:             "normal file name",
			itemName:         "overview",
			expectedSafeName: "overview.md",
		},
		{
			name:             "path traversal with forward slash",
			itemName:         "../../etc/passwd",
			expectedSafeName: ".._.._etc_passwd.md",
		},
		{
			name:             "path traversal with backslash",
			itemName:         "..\\..\\windows\\system32",
			expectedSafeName: ".._.._windows_system32.md",
		},
		{
			name:             "sub directory name",
			itemName:         "tables/users",
			expectedSafeName: "tables_users.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			safeName, fn, err := sanitizeDocFileName(docsDir, tt.itemName)
			require.NoError(t, err)
			require.Equal(t, tt.expectedSafeName, safeName)

			// Ensure target file stays strictly inside docsDir
			rel, err := filepath.Rel(docsDir, fn)
			require.NoError(t, err)
			require.NotEqual(t, "..", rel)
			require.False(t, strings.HasPrefix(rel, ".."+string(filepath.Separator)), "path escaped docsDir: %s", rel)
		})
	}
}
