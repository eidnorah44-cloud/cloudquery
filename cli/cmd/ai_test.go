package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "valid simple filename",
			input:     "config",
			expected:  "config",
			expectErr: false,
		},
		{
			name:      "path traversal attempt with relative directory",
			input:     "../../etc/passwd",
			expected:  "passwd",
			expectErr: false,
		},
		{
			name:      "path traversal attempt with dot dot prefix",
			input:     "../secret",
			expected:  "secret",
			expectErr: false,
		},
		{
			name:      "nested path",
			input:     "subfolder/myconfig",
			expected:  "myconfig",
			expectErr: false,
		},
		{
			name:      "empty or dot input",
			input:     "..",
			expected:  "",
			expectErr: true,
		},
		{
			name:      "single dot input",
			input:     ".",
			expected:  "",
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := sanitizeFilename(tc.input)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestCreateSpecFile_PathTraversalPrevention(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tempDir))
	defer func() { _ = os.Chdir(origDir) }()

	subDir := filepath.Join(tempDir, "subdir")
	require.NoError(t, os.Mkdir(subDir, 0755))

	// Attempt path traversal into subDir
	input := "subdir/../../escaped_spec"
	content := "kind: source"

	err = createSpecFile(input, content)
	require.NoError(t, err)

	// Verify file was created in current directory (tempDir), NOT outside or in unexpected location
	expectedPath := filepath.Join(tempDir, "escaped_spec.yaml")
	require.FileExists(t, expectedPath)

	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	require.Equal(t, content, string(data))
}

func TestCreateSQLFile_PathTraversalPrevention(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tempDir))
	defer func() { _ = os.Chdir(origDir) }()

	input := "../../../escaped_query"
	content := "SELECT 1;"

	err = createSQLFile(input, content)
	require.NoError(t, err)

	expectedPath := filepath.Join(tempDir, "escaped_query.sql")
	require.FileExists(t, expectedPath)

	data, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	require.Equal(t, content, string(data))
}
