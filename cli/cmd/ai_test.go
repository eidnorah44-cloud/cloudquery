package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAISanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{input: "valid_name", expected: "valid_name", hasError: false},
		{input: "../../../etc/passwd", expected: "passwd", hasError: false},
		{input: "nested/dir/myconfig", expected: "myconfig", hasError: false},
		{input: ".", expected: "", hasError: true},
		{input: "..", expected: "", hasError: true},
		{input: "/", expected: "", hasError: true},
	}

	for _, tt := range tests {
		got, err := sanitizeFilename(tt.input)
		if tt.hasError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		}
	}
}

func TestAICreateSpecAndSQLFile_PathTraversal(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origDir) }()

	require.NoError(t, os.Chdir(tempDir))

	// Attempting path traversal with createSpecFile
	err = createSpecFile("../../outside_spec", "content")
	require.NoError(t, err)

	// Verify file was created in current tempDir and not outside
	_, err = os.Stat(filepath.Join(tempDir, "outside_spec.yaml"))
	require.NoError(t, err)

	// Attempting path traversal with createSQLFile
	err = createSQLFile("../../outside_sql", "SELECT 1;")
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tempDir, "outside_sql.sql"))
	require.NoError(t, err)
}
