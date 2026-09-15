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
	val []byte
	key []byte
}

type SecretAwareRedactor struct {
	secrets      map[string]string
	secretsBytes []secretPair
}

func NewSecretAwareRedactor() *SecretAwareRedactor {
	return &SecretAwareRedactor{secrets: make(map[string]string)}
}

func (s *SecretAwareRedactor) RedactStr(msg string) string {
	return string(s.RedactBytes([]byte(msg)))
}

// RedactBytes replaces sensitive values in msg with their corresponding environment variable names.
// Performance Optimization: Uses pre-converted byte slices (secretsBytes) to avoid repeated heap
// allocations for string-to-byte-slice conversions on every RedactBytes / stream Write call.
func (s *SecretAwareRedactor) RedactBytes(msg []byte) []byte {
	for _, pair := range s.secretsBytes {
		msg = bytes.ReplaceAll(msg, pair.val, pair.key)
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

		if prevKey, ok := s.secrets[parts[1]]; !ok {
			s.secrets[parts[1]] = parts[0]
			s.secretsBytes = append(s.secretsBytes, secretPair{
				val: []byte(parts[1]),
				key: []byte(parts[0]),
			})
		} else if prevKey != parts[0] {
			s.secrets[parts[1]] = parts[0]
			for i := range s.secretsBytes {
				if string(s.secretsBytes[i].val) == parts[1] {
					s.secretsBytes[i].key = []byte(parts[0])
					break
				}
			}
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
