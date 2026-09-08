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

func TestPluginDocsDownload_PathTraversal(t *testing.T) {
	t.Setenv("CLOUDQUERY_API_KEY", "testkey")

	tempDir := t.TempDir()
	docsDir := filepath.Join(tempDir, "docs")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/plugins/test_team/source/test_plugin/versions/v1.0.0/docs":
			w.WriteHeader(http.StatusOK)
			resp := map[string]interface{}{
				"items": []map[string]string{
					{
						"name":    "../../../../etc/passwd",
						"content": "malicious_content",
					},
					{
						"name":    "overview",
						"content": "overview_content",
					},
				},
				"total": 2,
			}
			b, _ := json.Marshal(resp)
			_, err := w.Write(b)
			require.NoError(t, err)
		case "/teams":
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"items":[{"name":"test_team","displayName":"Test Team"}]}`))
			require.NoError(t, err)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer ts.Close()

	cmd := NewCmdRoot()
	t.Setenv(envAPIURL, ts.URL)
	args := append([]string{"plugin", "docs", "download", "test_team/source/test_plugin@v1.0.0", "-D", docsDir}, testCommandArgs(t)...)
	cmd.SetArgs(args)
	err := cmd.Execute()
	require.NoError(t, err)

	// Verify that passwd was sanitized to passwd.md inside docsDir, and not written to /etc/passwd or outside docsDir
	expectedFile := filepath.Join(docsDir, "passwd.md")
	_, err = os.Stat(expectedFile)
	assert.NoError(t, err, "expected file %s to exist inside docsDir", expectedFile)

	content, err := os.ReadFile(expectedFile)
	require.NoError(t, err)
	assert.Equal(t, "malicious_content", string(content))

	// Verify normal file overview.md exists inside docsDir
	overviewFile := filepath.Join(docsDir, "overview.md")
	_, err = os.Stat(overviewFile)
	assert.NoError(t, err, "expected file %s to exist inside docsDir", overviewFile)
}
