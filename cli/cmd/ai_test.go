package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple filename",
			input:    "my_config",
			expected: "my_config",
		},
		{
			name:     "path traversal attempt with relative dots",
			input:    "../../etc/passwd",
			expected: "passwd",
		},
		{
			name:     "path traversal with leading slashes",
			input:    "/etc/shadow",
			expected: "shadow",
		},
		{
			name:     "flag injection attempt with leading dash",
			input:    "--option",
			expected: "option",
		},
		{
			name:     "hidden file attempt with leading dot",
			input:    ".hidden_file",
			expected: "hidden_file",
		},
		{
			name:     "empty filename resulting from dots only",
			input:    "..",
			expected: "cloudquery",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAIFilenameSanitization(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(tmpDir))
	defer func() {
		_ = os.Chdir(origDir)
	}()

	t.Run("createSpecFile prevents path traversal", func(t *testing.T) {
		err := createSpecFile("../../../test_spec", "content: test")
		require.NoError(t, err)

		// Check file was created in current directory with sanitized name
		_, err = os.Stat("test_spec.yaml")
		assert.NoError(t, err, "file should exist in local temp dir with sanitized name")

		// Check file was not created outside
		_, err = os.Stat(filepath.Join(tmpDir, "..", "test_spec.yaml"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("createSQLFile prevents path traversal", func(t *testing.T) {
		err := createSQLFile("../../test_sql", "SELECT 1;")
		require.NoError(t, err)

		// Check file was created in current directory with sanitized name
		_, err = os.Stat("test_sql.sql")
		assert.NoError(t, err, "file should exist in local temp dir with sanitized name")
	})
}
