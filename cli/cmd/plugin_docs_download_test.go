package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	cloudquery_api "github.com/cloudquery/cloudquery-api-go"
	"github.com/stretchr/testify/require"
)

func TestPluginDocsDownload_PathTraversal(t *testing.T) {
	// Create a mock server that returns items with path traversal attempts
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := cloudquery_api.ListPluginVersionDocs200Response{
			Items: []cloudquery_api.PluginDocsPage{
				{
					Name:    "../../etc/passwd",
					Content: "malicious_content_passwd",
				},
				{
					Name:    "..\\..\\Windows\\System32\\config",
					Content: "malicious_content_win",
				},
				{
					Name:    "normal_doc",
					Content: "normal_content",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	t.Setenv("CLOUDQUERY_API_URL", ts.URL)
	t.Setenv("CLOUDQUERY_API_KEY", "testkey")

	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")

	cmd := newCmdPluginDocsDownload()
	cmd.SetArgs([]string{"test-team/source/test-plugin@v1.0.0", "-D", docsDir})

	err := cmd.ExecuteContext(context.Background())
	require.NoError(t, err)

	// Ensure files were written strictly inside docsDir with sanitized names
	entries, err := os.ReadDir(docsDir)
	require.NoError(t, err)
	require.Len(t, entries, 3)

	expectedFiles := map[string]string{
		"passwd.md":     "malicious_content_passwd",
		"config.md":     "malicious_content_win",
		"normal_doc.md": "normal_content",
	}

	for _, entry := range entries {
		require.False(t, entry.IsDir())
		expectedContent, exists := expectedFiles[entry.Name()]
		require.True(t, exists, "unexpected file created: %s", entry.Name())

		content, err := os.ReadFile(filepath.Join(docsDir, entry.Name()))
		require.NoError(t, err)
		require.Equal(t, expectedContent, string(content))
	}

	// Verify no files escaped outside docsDir into tmpDir
	tmpEntries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	require.Len(t, tmpEntries, 1) // only docsDir
	require.Equal(t, "docs", tmpEntries[0].Name())
}
