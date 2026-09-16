package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		expectErr bool
	}{
		{
			name:      "valid filename",
			filename:  "test_config",
			expectErr: false,
		},
		{
			name:      "empty filename",
			filename:  "",
			expectErr: true,
		},
		{
			name:      "relative path traversal with dots",
			filename:  "../etc/passwd",
			expectErr: true,
		},
		{
			name:      "subdirectory traversal",
			filename:  "subdir/config",
			expectErr: true,
		},
		{
			name:      "windows path separator traversal",
			filename:  "subdir\\config",
			expectErr: true,
		},
		{
			name:      "dot dot alone",
			filename:  "..",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateFilename(tt.filename)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.filename, got)
			}
		})
	}
}

func TestCreateSpecAndSQLFile_Security(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origDir)
	}()
	require.NoError(t, os.Chdir(tempDir))

	// Test valid file creation
	err = createSpecFile("valid_spec", "key: value")
	require.NoError(t, err)
	require.FileExists(t, filepath.Join(tempDir, "valid_spec.yaml"))

	err = createSQLFile("valid_sql", "SELECT 1;")
	require.NoError(t, err)
	require.FileExists(t, filepath.Join(tempDir, "valid_sql.sql"))

	// Test path traversal rejection
	err = createSpecFile("../traversal_spec", "key: value")
	require.Error(t, err)
	require.NoFileExists(t, filepath.Join(tempDir, "..", "traversal_spec.yaml"))

	err = createSQLFile("sub/traversal_sql", "SELECT 1;")
	require.Error(t, err)

	// Test cloudqueryTest rejection with path traversal
	res := cloudqueryTest("../traversal_test")
	require.Contains(t, res, "cloudquery test failed: invalid filename")
}
