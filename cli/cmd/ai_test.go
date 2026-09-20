package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAIFiles_PathTraversal(t *testing.T) {
	// Change working directory to a temporary directory for safety during tests
	tempDir := t.TempDir()
	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origWd)
	}()
	require.NoError(t, os.Chdir(tempDir))

	subDir := filepath.Join(tempDir, "subdir")
	require.NoError(t, os.MkdirAll(subDir, 0755))

	// Test createSpecFile with path traversal payload
	traversalName := "../subdir/test_spec"
	err = createSpecFile(traversalName, "kind: spec")
	require.NoError(t, err)

	// Verify file was written to current directory (tempDir/test_spec.yaml), NOT subDir/test_spec.yaml
	require.FileExists(t, filepath.Join(tempDir, "test_spec.yaml"))
	require.NoFileExists(t, filepath.Join(subDir, "test_spec.yaml"))

	// Test createSQLFile with path traversal payload
	err = createSQLFile(traversalName, "SELECT 1;")
	require.NoError(t, err)

	// Verify file was written to current directory (tempDir/test_spec.sql), NOT subDir/test_spec.sql
	require.FileExists(t, filepath.Join(tempDir, "test_spec.sql"))
	require.NoFileExists(t, filepath.Join(subDir, "test_spec.sql"))
}
