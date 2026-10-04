## 2026-03-08 - Fast-path checking for stream redactor & string transformers

**Learning:** `SecretAwareRedactor` processes all CLI log streams and strings. Re-converting map keys to `[]byte` and invoking `bytes.ReplaceAll` unconditionally creates redundant heap allocations and buffer copies on every log line, even when no secrets are present.
**Action:** Pre-convert and cache slice/string representations of keys during configuration (`AddSecretEnv`), and use `strings.Contains` or `bytes.Contains` before performing replacement operations to short-circuit allocation on clean streams.
