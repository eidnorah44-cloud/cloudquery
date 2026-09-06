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

func TestPluginDocsDownload_PathTraversalSanitization(t *testing.T) {
	t.Setenv("CLOUDQUERY_API_KEY", "testkey")

	tempDir := t.TempDir()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/plugins/test-team/source/test-plugin/versions/v1.0.0/docs":
			checkAuthHeader(t, r)
			w.WriteHeader(http.StatusOK)
			resp := map[string]any{
				"items": []map[string]string{
					{
						"name":    "overview",
						"content": "# Overview",
					},
					{
						"name":    "tables/users",
						"content": "# Users Table",
					},
					{
						"name":    "../../etc/passwd",
						"content": "traversal attempt 1",
					},
					{
						"name":    "..\\..\\secret",
						"content": "traversal attempt 2",
					},
				},
			}
			b, err := json.Marshal(resp)
			require.NoError(t, err)
			_, err = w.Write(b)
			require.NoError(t, err)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	defer ts.Close()

	t.Setenv(envAPIURL, ts.URL)

	cmd := NewCmdRoot()
	args := append([]string{"plugin", "docs", "download", "test-team/source/test-plugin@v1.0.0", "-D", tempDir}, testCommandArgs(t)...)
	cmd.SetArgs(args)

	err := cmd.Execute()
	require.NoError(t, err)

	expectedFiles := []string{
		"overview.md",
		"tables_users.md",
		".._.._etc_passwd.md",
		".._.._secret.md",
	}

	for _, ef := range expectedFiles {
		fp := filepath.Join(tempDir, ef)
		_, err := os.Stat(fp)
		assert.NoError(t, err, "expected file %s to exist inside tempDir", ef)
	}

	entries, err := os.ReadDir(tempDir)
	require.NoError(t, err)
	assert.Len(t, entries, len(expectedFiles), "no extra files or subdirectories should be created")
}
