package secrets

import (
	"bytes"
	"io"
	"slices"
	"strings"
)

var allowedEnvPrefixes = []string{
	"_CQ_TEAM_NAME=",
	"HOME=",

	// injected by EKS, do not contain any sensitive information regardless
	"AWS_STS_REGIONAL_ENDPOINTS", "AWS_DEFAULT_REGION", "AWS_REGION",
	"AWS_WEB_IDENTITY_TOKEN_FILE", "AWS_ROLE_ARN", "AWS_ROLE_SESSION_NAME",
}

// minRedactingLength is the minimum length of an environment variable value for it to be redacted
const minRedactingLength = 4

type SecretAwareRedactor struct {
	secrets map[string]string
}

func NewSecretAwareRedactor() *SecretAwareRedactor {
	return &SecretAwareRedactor{secrets: make(map[string]string)}
}

func (s *SecretAwareRedactor) RedactStr(msg string) string {
	return string(s.RedactBytes([]byte(msg)))
}

func (s *SecretAwareRedactor) RedactBytes(msg []byte) []byte {
	keys := make([]string, 0, len(s.secrets))
	for v := range s.secrets {
		keys = append(keys, v)
	}
	// Sort keys by length in descending order to handle longer overlapping secrets first
	// (e.g. full assignment "KEY=VALUE" before standalone "VALUE"), with tie-breaking for determinism.
	slices.SortFunc(keys, func(a, b string) int {
		if len(a) != len(b) {
			return len(b) - len(a) // descending length
		}
		return strings.Compare(a, b)
	})

	for _, v := range keys {
		msg = bytes.ReplaceAll(msg, []byte(v), []byte(s.secrets[v]))
	}
	return msg
}

// AddSecretEnv registers environment variables for secret redaction.
// Security concern: prevents sensitive credentials, tokens, and raw environment key=value
// assignments from leaking in CLI output, error messages, or logs.
func (s *SecretAwareRedactor) AddSecretEnv(envs []string) {
	for _, v := range envs {
		if slices.ContainsFunc(allowedEnvPrefixes, func(prefix string) bool { return strings.HasPrefix(v, prefix) }) {
			continue
		}

		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			continue
		}

		// Redact raw KEY=VALUE assignment strings if logged in full
		if len(v) >= minRedactingLength {
			s.secrets[v] = parts[0]
		}

		// Redact the secret value itself if it meets the minimum redacting length
		if len(parts[1]) >= minRedactingLength {
			s.secrets[parts[1]] = parts[0]
		}
	}
}

type SecretAwareWriter struct {
	out      io.Writer
	redactor *SecretAwareRedactor
}

func NewSecretAwareWriter(out io.Writer, redactor *SecretAwareRedactor) *SecretAwareWriter {
	return &SecretAwareWriter{out: out, redactor: redactor}
}

func (s SecretAwareWriter) Write(p []byte) (n int, err error) {
	return s.out.Write(s.redactor.RedactBytes(p))
}
