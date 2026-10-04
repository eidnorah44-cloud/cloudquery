package publish

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadAddonFromResponse(t *testing.T) {
	content := "sample addon zip content"
	sum := sha256.Sum256([]byte(content))
	expectedChecksum := hex.EncodeToString(sum[:])

	tmpDir := t.TempDir()

	t.Run("successful download creates file with 0600 mode", func(t *testing.T) {
		destFile := filepath.Join(tmpDir, "addon.zip")
		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(content)),
		}

		err := DownloadAddonFromResponse(res, expectedChecksum, destFile)
		require.NoError(t, err)

		info, err := os.Stat(destFile)
		require.NoError(t, err)
		require.Equal(t, os.FileMode(0600), info.Mode().Perm())
	})

	t.Run("fails when file already exists", func(t *testing.T) {
		destFile := filepath.Join(tmpDir, "existing.zip")
		require.NoError(t, os.WriteFile(destFile, []byte("existing"), 0600))

		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(content)),
		}

		err := DownloadAddonFromResponse(res, expectedChecksum, destFile)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists")
		require.NotContains(t, err.Error(), "is a directory")
	})

	t.Run("fails when destination is a directory", func(t *testing.T) {
		destDir := filepath.Join(tmpDir, "existing_dir")
		require.NoError(t, os.Mkdir(destDir, 0755))

		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(content)),
		}

		err := DownloadAddonFromResponse(res, expectedChecksum, destDir)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already exists: is a directory")
	})
}
