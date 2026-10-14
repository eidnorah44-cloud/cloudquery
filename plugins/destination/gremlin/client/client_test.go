package client

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

const (
	defaultGremlinEndpoint = "ws://localhost:8182"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func TestPlugin(t *testing.T) {
	defaultInsecure := "false"
	if runtime.GOOS == "darwin" {
		defaultInsecure = "true" // x509 issue with Neptune and macOS: https://github.com/golang/go/issues/51991
	}
	insecure, _ := strconv.ParseBool(getenv("CQ_DEST_GREMLIN_INSECURE", defaultInsecure))

	ctx := context.Background()
	p := plugin.NewPlugin("gremlin", "development", New)
	s := &Spec{
		Endpoint: getenv("CQ_DEST_GREMLIN_ENDPOINT", defaultGremlinEndpoint),
		Insecure: insecure,
		Username: os.Getenv("CQ_DEST_GREMLIN_USERNAME"),
		Password: os.Getenv("CQ_DEST_GREMLIN_PASSWORD"),
	}
	s.SetDefaults()
	require.NoError(t, s.Validate())
	b, err := json.Marshal(s)
	require.NoError(t, err)

	err = p.Init(ctx, b, plugin.NewClientOptions{})
	require.NoError(t, err)

	plugin.TestWriterSuiteRunner(t,
		p,
		plugin.WriterTestSuiteTests{
			SkipInsert:       true, // we do "no PKs = all columns are PKs" in this destination
			SkipMigrate:      true,
			SkipDeleteRecord: true,
			SafeMigrations: plugin.SafeMigrations{
				AddColumn:    true,
				RemoveColumn: true,
			},
		},
		plugin.WithTestDataOptions(schema.TestSourceOptions{TimePrecision: time.Millisecond}),
	)
}

func TestGetAuthInfo_AWSHeaderRedaction(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "TEST_KEY")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "TEST_SECRET")
	t.Setenv("AWS_SESSION_TOKEN", "SENSITIVE_SESSION_TOKEN")

	var buf bytes.Buffer
	logger := zerolog.New(&buf).Level(zerolog.TraceLevel)

	c := &Client{
		logger: logger,
		spec: &Spec{
			Endpoint:  "wss://localhost:8182",
			AuthMode:  authModeAWS,
			AWSRegion: "us-east-1",
		},
	}

	authProvider, err := c.getAuthInfo(context.Background(), c.spec.Endpoint+"/gremlin")
	require.NoError(t, err)

	// Invoke the provider to trigger signing and trace logging
	authProvider.GetHeader()

	logOutput := buf.String()
	require.Contains(t, logOutput, `"iam_headers"`)
	require.Contains(t, logOutput, "*****")
	require.NotContains(t, logOutput, "SENSITIVE_SESSION_TOKEN")
}
