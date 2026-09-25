## 2025-05-18 - Path Traversal Prevention in Cross-Platform File Extraction
**Vulnerability:** Document download logic in `cli/cmd/plugin_docs_download.go` attempted path separator replacement using `string(filepath.Separator)`. On non-Windows platforms, backslashes `\` or POSIX slashes `/` on Windows were ignored, allowing cross-platform path traversal.
**Learning:** `filepath.Separator` only reflects the host OS separator. Slashes from untrusted remote payloads can use either `/` or `\` regardless of host OS.
**Prevention:** Validate remote file paths using `filepath.IsLocal` before processing, sanitize both `/` and `\`, and verify `filepath.Rel` does not escape target directory.
