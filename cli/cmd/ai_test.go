package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAISanitizePath(t *testing.T) {
	tempDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tempDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origDir)
	}()

	// Test createSpecFile path traversal prevention
	err = createSpecFile("../../traversal_spec", "test spec content")
	require.NoError(t, err)

	// Check that file was created in current directory (tempDir) and NOT parent directory
	require.FileExists(t, filepath.Join(tempDir, "traversal_spec.yaml"))
	require.NoFileExists(t, filepath.Join(tempDir, "..", "traversal_spec.yaml"))

	// Test createSQLFile path traversal prevention
	err = createSQLFile("../../../traversal_sql", "SELECT 1;")
	require.NoError(t, err)

	// Check that file was created in current directory (tempDir) and NOT parent directory
	require.FileExists(t, filepath.Join(tempDir, "traversal_sql.sql"))
	require.NoFileExists(t, filepath.Join(tempDir, "..", "traversal_sql.sql"))
}
