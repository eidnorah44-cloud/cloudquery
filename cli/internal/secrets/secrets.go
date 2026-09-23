package secrets

import (
	"bytes"
	"cmp"
	"io"
	"slices"
	"strings"
	"sync"
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

type secretKV struct {
	val []byte
	key []byte
}

type SecretAwareRedactor struct {
	mu      sync.RWMutex
	secrets map[string]string
}

func NewSecretAwareRedactor() *SecretAwareRedactor {
	return &SecretAwareRedactor{secrets: make(map[string]string)}
}

func (s *SecretAwareRedactor) RedactStr(msg string) string {
	return string(s.RedactBytes([]byte(msg)))
}

func (s *SecretAwareRedactor) RedactBytes(msg []byte) []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Security concern: Longer secrets containing shorter secrets as substrings must be redacted first.
	// Redacting a shorter secret first would corrupt the longer secret and leak sensitive substrings.
	items := make([]secretKV, 0, len(s.secrets))
	for v, k := range s.secrets {
		items = append(items, secretKV{val: []byte(v), key: []byte(k)})
	}
	slices.SortFunc(items, func(a, b secretKV) int {
		if c := cmp.Compare(len(b.val), len(a.val)); c != 0 {
			return c
		}
		return bytes.Compare(a.val, b.val)
	})

	for _, item := range items {
		msg = bytes.ReplaceAll(msg, item.val, item.key)
	}
	return msg
}

func (s *SecretAwareRedactor) AddSecretEnv(envs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
