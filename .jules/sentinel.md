## 2025-05-18 - Secret Redaction Substring Leakage & Data Race
**Vulnerability:** Random map iteration order in secret redactor could redact shorter secret substrings before longer secrets (leaking sensitive suffixes), and concurrent map read/write access caused runtime panics.
**Learning:** In Go, string replacement for secret scrubbing must sort secret values by length descending to prevent shorter secret matches from partially corrupting and leaking longer secrets. Map accesses in log writers or redactors must be protected with `sync.RWMutex` against concurrent goroutine writes.
**Prevention:** Always sort secrets by length descending before performing replacements, and protect redactor map state with `sync.RWMutex`.
