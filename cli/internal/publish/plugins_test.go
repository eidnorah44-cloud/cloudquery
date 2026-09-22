package publish

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTransportWithRegistryAuth(t *testing.T) {
	t.Run("default transport", func(t *testing.T) {
		tr := newTransportWithRegistryAuth(true, "test-token")
		require.NotNil(t, tr)
		require.NotNil(t, tr.baseTransport)
		require.NotNil(t, tr.baseTransport.TLSClientConfig)
		require.True(t, tr.baseTransport.TLSClientConfig.InsecureSkipVerify)
		require.Equal(t, "test-token", tr.registryAuth)
	})

	t.Run("nil TLSClientConfig on default transport", func(t *testing.T) {
		orig := http.DefaultTransport
		defer func() {
			http.DefaultTransport = orig
		}()

		// Set DefaultTransport TLSClientConfig to nil to simulate custom/uninitialized default transport
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = nil
		http.DefaultTransport = customTransport

		tr := newTransportWithRegistryAuth(false, "test-token-2")
		require.NotNil(t, tr)
		require.NotNil(t, tr.baseTransport)
		require.NotNil(t, tr.baseTransport.TLSClientConfig)
		require.False(t, tr.baseTransport.TLSClientConfig.InsecureSkipVerify)
		require.Equal(t, "test-token-2", tr.registryAuth)
	})
}
