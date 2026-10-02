## 2025-05-18 - Cross-Platform Path Traversal in File Downloads
**Vulnerability:** In `cli/cmd/plugin_docs_download.go`, doc item names were sanitized using `strings.ReplaceAll(name, string(filepath.Separator), "_")`.
**Learning:** `filepath.Separator` is platform dependent (`/` on Unix, `\` on Windows). Using `filepath.Separator` allowed cross-platform path traversal payloads (e.g. forward slashes on Windows or backslashes on Unix) to bypass sanitization.
**Prevention:** Always replace both `/` and `\` slashes explicitly when sanitizing untrusted external filenames, and validate relative path resolution using `filepath.Rel` against the target output directory.
