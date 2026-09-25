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

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/teams" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"items":[{"name":"test-team","displayName":"Test Team"}]}`))
			return
		}
		if r.URL.Path == "/plugins/test-team/source/test-plugin/versions/v1.0.0/docs" {
			w.WriteHeader(http.StatusOK)
			resp := map[string]interface{}{
				"items": []map[string]string{
					{
						"name":    "overview/getting-started",
						"content": "getting started content",
					},
					{
						"name":    "normal_doc",
						"content": "normal content",
					},
				},
			}
			b, _ := json.Marshal(resp)
			_, _ = w.Write(b)
			return
		}
		t.Logf("Unhandled request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	t.Setenv(envAPIURL, ts.URL)

	tempDir := t.TempDir()
	docsDir := filepath.Join(tempDir, "docs")

	cmd := NewCmdRoot()
	args := append([]string{"plugin", "docs", "download", "test-team/source/test-plugin@v1.0.0", "-D", docsDir}, testCommandArgs(t)...)
	cmd.SetArgs(args)

	err := cmd.Execute()
	require.NoError(t, err)

	// Verify that "overview/getting-started" was sanitized to "overview_getting-started.md" inside docsDir
	sanitizedFile := filepath.Join(docsDir, "overview_getting-started.md")
	content, err := os.ReadFile(sanitizedFile)
	require.NoError(t, err, "sanitized file should exist inside docsDir")
	assert.Equal(t, "getting started content", string(content))

	normalFile := filepath.Join(docsDir, "normal_doc.md")
	content, err = os.ReadFile(normalFile)
	require.NoError(t, err, "normal file should exist inside docsDir")
	assert.Equal(t, "normal content", string(content))
}

func TestPluginDocsDownload_PathTraversalRejection(t *testing.T) {
	t.Setenv("CLOUDQUERY_API_KEY", "testkey")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/teams" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"items":[{"name":"test-team","displayName":"Test Team"}]}`))
			return
		}
		if r.URL.Path == "/plugins/test-team/source/test-plugin/versions/v1.0.0/docs" {
			w.WriteHeader(http.StatusOK)
			resp := map[string]interface{}{
				"items": []map[string]string{
					{
						"name":    "../../evil",
						"content": "malicious content",
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

	t.Setenv(envAPIURL, ts.URL)

	tempDir := t.TempDir()
	docsDir := filepath.Join(tempDir, "docs")

	cmd := NewCmdRoot()
	args := append([]string{"plugin", "docs", "download", "test-team/source/test-plugin@v1.0.0", "-D", docsDir}, testCommandArgs(t)...)
	cmd.SetArgs(args)

	err := cmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal detected")

	evilFile := filepath.Join(tempDir, "evil.md")
	_, err = os.Stat(evilFile)
	assert.True(t, os.IsNotExist(err), "evil file should not exist outside docsDir")
}
