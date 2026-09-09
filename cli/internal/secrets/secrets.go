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

type secretPair struct {
	valueStr string
	value    []byte
	key      []byte
}

type SecretAwareRedactor struct {
	secrets map[string]string
	pairs   []secretPair
}

func NewSecretAwareRedactor() *SecretAwareRedactor {
	return &SecretAwareRedactor{secrets: make(map[string]string)}
}

// RedactStr redacts secrets from a string.
// Fast path: if no secrets are registered or present in msg, it returns msg without any heap allocations.
func (s *SecretAwareRedactor) RedactStr(msg string) string {
	if len(s.pairs) == 0 || len(msg) == 0 {
		return msg
	}
	// Check if any secret is present before converting string to byte slice.
	// This avoids allocating []byte(msg) and string(redacted) when no secret is in the message.
	hasSecret := false
	for _, p := range s.pairs {
		if strings.Contains(msg, p.valueStr) {
			hasSecret = true
			break
		}
	}
	if !hasSecret {
		return msg
	}
	return string(s.RedactBytes([]byte(msg)))
}

// RedactBytes redacts secrets from a byte slice.
// Fast path: if no secrets are registered or present in msg, it returns msg without modifying it.
func (s *SecretAwareRedactor) RedactBytes(msg []byte) []byte {
	if len(s.pairs) == 0 || len(msg) == 0 {
		return msg
	}
	// Check if any secret is present before calling bytes.ReplaceAll.
	// Pre-converted byte slices (s.pairs) avoid allocations in the hot path.
	hasSecret := false
	for _, p := range s.pairs {
		if bytes.Contains(msg, p.value) {
			hasSecret = true
			break
		}
	}
	if !hasSecret {
		return msg
	}
	for _, p := range s.pairs {
		msg = bytes.ReplaceAll(msg, p.value, p.key)
	}
	return msg
}

func (s *SecretAwareRedactor) AddSecretEnv(envs []string) {
	for _, v := range envs {
		if slices.ContainsFunc(allowedEnvPrefixes, func(prefix string) bool { return strings.HasPrefix(v, prefix) }) {
			continue
		}

		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 || len(parts[1]) < minRedactingLength {
			continue
		}

		s.secrets[parts[1]] = parts[0]
	}

	s.pairs = make([]secretPair, 0, len(s.secrets))
	for v, k := range s.secrets {
		s.pairs = append(s.pairs, secretPair{
			valueStr: v,
			value:    []byte(v),
			key:      []byte(k),
		})
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
