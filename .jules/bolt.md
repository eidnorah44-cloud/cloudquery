## 2025-01-01 - Caching byte representations in output stream redactors
**Learning:** `SecretAwareWriter` wraps all CLI output and log writers, calling `RedactBytes` on every output write. Converting secret strings to `[]byte` during map iteration in `RedactBytes` generates millions of unnecessary heap allocations in stream-heavy CLI applications.
**Action:** Pre-convert and store secret keys/values as `[]byte` pairs during secret registration to achieve zero key/value conversion allocations during stream redaction.
