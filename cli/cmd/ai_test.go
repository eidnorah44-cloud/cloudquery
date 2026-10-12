package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{input: "valid_spec", expected: "valid_spec", wantErr: false},
		{input: "sub/dir/spec", expected: "spec", wantErr: false},
		{input: "../../../etc/passwd", expected: "passwd", wantErr: false},
		{input: ".", expected: "", wantErr: true},
		{input: "..", expected: "", wantErr: true},
		{input: "/", expected: "", wantErr: true},
	}

	for _, tt := range tests {
		got, err := sanitizeFilename(tt.input)
		if tt.wantErr {
			require.Error(t, err, "input: %s", tt.input)
		} else {
			require.NoError(t, err, "input: %s", tt.input)
			require.Equal(t, tt.expected, got)
		}
	}
}

func TestCreateSpecFile_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	require.NoError(t, os.Mkdir(subDir, 0755))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(subDir))

	traversalName := "../escaped_spec"
	err = createSpecFile(traversalName, "test content")
	require.NoError(t, err)

	_, errParent := os.Stat(filepath.Join(tmpDir, "escaped_spec.yaml"))
	require.True(t, os.IsNotExist(errParent), "file should not be created in parent directory")

	_, errSub := os.Stat(filepath.Join(subDir, "escaped_spec.yaml"))
	require.NoError(t, errSub, "file should be created in current working directory")
}

func TestCreateSQLFile_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	require.NoError(t, os.Mkdir(subDir, 0755))

	origWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(origWd) }()
	require.NoError(t, os.Chdir(subDir))

	traversalName := "../escaped_sql"
	err = createSQLFile(traversalName, "SELECT 1;")
	require.NoError(t, err)

	_, errParent := os.Stat(filepath.Join(tmpDir, "escaped_sql.sql"))
	require.True(t, os.IsNotExist(errParent), "file should not be created in parent directory")

	_, errSub := os.Stat(filepath.Join(subDir, "escaped_sql.sql"))
	require.NoError(t, errSub, "file should be created in current working directory")
}

func TestCloudqueryTest_InvalidFilename(t *testing.T) {
	output := cloudqueryTest("..")
	require.Contains(t, output, "cloudquery test failed: invalid filename")
}
