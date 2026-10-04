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

// secretPair holds pre-computed byte slices and string representations of secret key-value pairs.
// Pre-computing these avoids heap allocations during log redaction on hot paths.
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

// RedactStr redacts secret values in msg.
// Performance optimization: We check if msg contains any secret using strings.Contains first.
// If no secret is found (the common case for logs), we return the string directly without
// converting to byte slices or making heap allocations.
func (s *SecretAwareRedactor) RedactStr(msg string) string {
	if len(msg) == 0 || len(s.pairs) == 0 {
		return msg
	}
	hasSecret := false
	for _, pair := range s.pairs {
		if strings.Contains(msg, pair.valueStr) {
			hasSecret = true
			break
		}
	}
	if !hasSecret {
		return msg
	}
	return string(s.RedactBytes([]byte(msg)))
}

// RedactBytes redacts secret values in msg.
// Performance optimization: We reuse pre-computed []byte slices for secret values and names,
// and check bytes.Contains before calling bytes.ReplaceAll. This avoids slice allocations and
// buffer copying when secrets are not present in the byte stream.
func (s *SecretAwareRedactor) RedactBytes(msg []byte) []byte {
	if len(msg) == 0 || len(s.pairs) == 0 {
		return msg
	}
	for _, pair := range s.pairs {
		if bytes.Contains(msg, pair.value) {
			msg = bytes.ReplaceAll(msg, pair.value, pair.key)
		}
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

		secretVal := parts[1]
		secretName := parts[0]

		if oldName, exists := s.secrets[secretVal]; exists {
			if oldName != secretName {
				s.secrets[secretVal] = secretName
				for i := range s.pairs {
					if s.pairs[i].valueStr == secretVal {
						s.pairs[i].key = []byte(secretName)
						break
					}
				}
			}
		} else {
			s.secrets[secretVal] = secretName
			s.pairs = append(s.pairs, secretPair{
				valueStr: secretVal,
				value:    []byte(secretVal),
				key:      []byte(secretName),
			})
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
