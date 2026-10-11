package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFunctionArgString(t *testing.T) {
	require.Equal(t, "", getFunctionArgString(nil, "key"))
	require.Equal(t, "", getFunctionArgString(map[string]any{"other": "val"}, "key"))
	require.Equal(t, "", getFunctionArgString(map[string]any{"key": 123}, "key"))
	require.Equal(t, "hello", getFunctionArgString(map[string]any{"key": "hello"}, "key"))
}

func TestCreateSpecFile_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	// Attempt path traversal with directory components
	err := createSpecFile("../../sub/test_spec", "content: test")
	require.NoError(t, err)

	// File must be created strictly as "test_spec.yaml" in tmpDir, not in ../../sub/
	_, err = os.Stat(filepath.Join(tmpDir, "test_spec.yaml"))
	require.NoError(t, err, "test_spec.yaml should exist in current directory")

	// Verify invalid inputs are rejected
	require.Error(t, createSpecFile("..", "content"))
	require.Error(t, createSpecFile(".", "content"))
	require.Error(t, createSpecFile("", "content"))
}

func TestCreateSQLFile_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)

	err := createSQLFile("../../../query", "SELECT 1;")
	require.NoError(t, err)

	_, err = os.Stat(filepath.Join(tmpDir, "query.sql"))
	require.NoError(t, err, "query.sql should exist in current directory")

	require.Error(t, createSQLFile("..", "SELECT 1;"))
	require.Error(t, createSQLFile(".", "SELECT 1;"))
	require.Error(t, createSQLFile("", "SELECT 1;"))
}

func TestCloudqueryTest_PathTraversal(t *testing.T) {
	res := cloudqueryTest("..")
	require.Contains(t, res, "invalid filename")

	res = cloudqueryTest(".")
	require.Contains(t, res, "invalid filename")
}
