package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAISpecAndSQLFileCreationPathTraversal(t *testing.T) {
	// Create a temporary directory to act as working directory
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(origDir)
	}()

	// Subdirectory that should NOT be escaped into
	subDir := filepath.Join(tmpDir, "subdir")
	require.NoError(t, os.Mkdir(subDir, 0755))

	err = os.Chdir(subDir)
	require.NoError(t, err)

	// Attempt path traversal input: "../traversal_spec"
	traversalInput := "../traversal_spec"
	specContent := "kind: source\nname: test"

	err = createSpecFile(traversalInput, specContent)
	require.NoError(t, err)

	// Verify file was created as base filename in current directory (subDir)
	expectedFileInSubdir := filepath.Join(subDir, "traversal_spec.yaml")
	require.FileExists(t, expectedFileInSubdir)

	// Verify file was NOT created in parent directory
	unexpectedFileInParent := filepath.Join(tmpDir, "traversal_spec.yaml")
	require.NoFileExists(t, unexpectedFileInParent)

	// Attempt path traversal input for SQL file
	sqlTraversalInput := "../traversal_sql"
	sqlContent := "SELECT 1;"

	err = createSQLFile(sqlTraversalInput, sqlContent)
	require.NoError(t, err)

	// Verify SQL file was created in current directory (subDir)
	expectedSQLInSubdir := filepath.Join(subDir, "traversal_sql.sql")
	require.FileExists(t, expectedSQLInSubdir)

	// Verify SQL file was NOT created in parent directory
	unexpectedSQLInParent := filepath.Join(tmpDir, "traversal_sql.sql")
	require.NoFileExists(t, unexpectedSQLInParent)
}
