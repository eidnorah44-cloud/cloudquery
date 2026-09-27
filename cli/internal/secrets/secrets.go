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
	rep []byte
}

type SecretAwareRedactor struct {
	secrets map[string]string
	// pairs caches pre-converted []byte representations of secret values and names
	// to avoid heap allocations on every RedactBytes call in hot log/output paths.
	pairs []secretPair
}

func NewSecretAwareRedactor() *SecretAwareRedactor {
	return &SecretAwareRedactor{secrets: make(map[string]string)}
}

func (s *SecretAwareRedactor) RedactStr(msg string) string {
	return string(s.RedactBytes([]byte(msg)))
}

func (s *SecretAwareRedactor) RedactBytes(msg []byte) []byte {
	for _, pair := range s.pairs {
		msg = bytes.ReplaceAll(msg, pair.val, pair.rep)
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

		if _, ok := s.secrets[parts[1]]; !ok {
			s.secrets[parts[1]] = parts[0]
			s.pairs = append(s.pairs, secretPair{
				val: []byte(parts[1]),
				rep: []byte(parts[0]),
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
