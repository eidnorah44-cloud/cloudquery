package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginDocsDownloadSanitization(t *testing.T) {
	t.Setenv("CLOUDQUERY_API_KEY", "testkey")

	tempDocsDir := t.TempDir()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/plugins/test-team/source/test-plugin/versions/v1.0.0/docs" {
			w.WriteHeader(http.StatusOK)
			resp := map[string]any{
				"items": []map[string]string{
					{
						"name":    "../../traversal_posix",
						"content": "# Posix Traversal Doc",
					},
					{
						"name":    "..\\..\\traversal_win",
						"content": "# Win Traversal Doc",
					},
					{
						"name":    "normal_doc",
						"content": "# Normal Doc",
					},
				},
			}
			b, _ := json.Marshal(resp)
			_, err := w.Write(b)
			require.NoError(t, err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"not found"}`))
	}))
	defer ts.Close()

	cmd := NewCmdRoot()
	t.Setenv(envAPIURL, ts.URL)
	args := append([]string{"plugin", "docs", "download", "test-team/source/test-plugin@v1.0.0", "-D", tempDocsDir}, testCommandArgs(t)...)
	cmd.SetArgs(args)
	err := cmd.Execute()
	require.NoError(t, err)

	// Check that sanitized files exist in tempDocsDir
	expectedPosixFile := filepath.Join(tempDocsDir, ".._.._traversal_posix.md")
	expectedWinFile := filepath.Join(tempDocsDir, ".._.._traversal_win.md")
	expectedNormalFile := filepath.Join(tempDocsDir, "normal_doc.md")

	_, err = os.Stat(expectedPosixFile)
	assert.NoError(t, err, "expected sanitized posix traversal file to exist in docs dir")

	_, err = os.Stat(expectedWinFile)
	assert.NoError(t, err, "expected sanitized win traversal file to exist in docs dir")

	_, err = os.Stat(expectedNormalFile)
	assert.NoError(t, err, "expected normal file to exist in docs dir")

	// Ensure no file was created in parent directory of tempDocsDir
	parentDir := filepath.Dir(tempDocsDir)
	escapedPosixFile := filepath.Join(parentDir, "traversal_posix.md")
	escapedWinFile := filepath.Join(parentDir, "traversal_win.md")

	_, err = os.Stat(escapedPosixFile)
	assert.True(t, os.IsNotExist(err), "file should not exist outside docs dir")

	_, err = os.Stat(escapedWinFile)
	assert.True(t, os.IsNotExist(err), "file should not exist outside docs dir")
}
