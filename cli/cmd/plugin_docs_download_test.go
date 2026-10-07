package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginDocsDownloadPathTraversal(t *testing.T) {
	t.Setenv("CLOUDQUERY_API_KEY", "testkey")

	tempDir := t.TempDir()
	docsDir := filepath.Join(tempDir, "docs")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/teams":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"items":[{"name":"test-team","displayName":"Test Team"}]}`))
			return
		case "/plugins/test-team/source/test-plugin/versions/v1.0.0/docs":
			w.WriteHeader(http.StatusOK)
			resp := map[string]interface{}{
				"items": []map[string]string{
					{
						"name":    "../traversal_test",
						"content": "test content traversal 1",
					},
					{
						"name":    "..\\win_traversal_test",
						"content": "test content traversal 2",
					},
					{
						"name":    "normal_doc",
						"content": "test content normal",
					},
				},
			}
			b, _ := json.Marshal(resp)
			_, _ = w.Write(b)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	cmd := NewCmdRoot()
	t.Setenv(envAPIURL, ts.URL)
	args := append([]string{"plugin", "docs", "download", "test-team/source/test-plugin@v1.0.0", "-D", docsDir}, testCommandArgs(t)...)
	cmd.SetArgs(args)

	err := cmd.Execute()
	require.NoError(t, err)

	// Verify files are saved safely inside docsDir without escaping
	expectedFile1 := filepath.Join(docsDir, ".._traversal_test.md")
	expectedFile2 := filepath.Join(docsDir, ".._win_traversal_test.md")
	expectedFile3 := filepath.Join(docsDir, "normal_doc.md")

	assert.FileExists(t, expectedFile1)
	assert.FileExists(t, expectedFile2)
	assert.FileExists(t, expectedFile3)

	// Confirm that no files were written outside docsDir
	escapedFile := filepath.Join(tempDir, "traversal_test.md")
	assert.NoFileExists(t, escapedFile)
}
