package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAISanitizeFilename(t *testing.T) {
	tempDir := t.TempDir()

	// Change working directory to tempDir for test isolation
	origDir, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origDir)
	}()
	require.NoError(t, os.Chdir(tempDir))

	t.Run("createSpecFile prevents path traversal", func(t *testing.T) {
		maliciousPath := "../../malicious_spec"
		content := "name: test"

		err := createSpecFile(maliciousPath, content)
		require.NoError(t, err)

		// Check that file was created in current directory with base name
		expectedPath := filepath.Join(tempDir, "malicious_spec.yaml")
		require.FileExists(t, expectedPath)

		// Verify file content
		data, err := os.ReadFile(expectedPath)
		require.NoError(t, err)
		require.Equal(t, content, string(data))

		// Check that no file was created outside tempDir
		outsidePath := filepath.Join(tempDir, "..", "..", "malicious_spec.yaml")
		require.NoFileExists(t, outsidePath)
	})

	t.Run("createSQLFile prevents path traversal", func(t *testing.T) {
		maliciousPath := "../../../malicious_sql"
		content := "SELECT 1;"

		err := createSQLFile(maliciousPath, content)
		require.NoError(t, err)

		expectedPath := filepath.Join(tempDir, "malicious_sql.sql")
		require.FileExists(t, expectedPath)

		data, err := os.ReadFile(expectedPath)
		require.NoError(t, err)
		require.Equal(t, content, string(data))
	})
}
