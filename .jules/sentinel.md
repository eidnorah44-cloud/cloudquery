## 2026-03-31 - Path Traversal in Addon Document Image Processing

**Vulnerability:** Markdown image parser in `cli/internal/publish/images/images.go` accepted relative (`../`) and absolute local paths without verifying directory confinement, allowing arbitrary local file read and exfiltration during addon publication.

**Learning:** `filepath.Join(baseDir, relativePath)` alone does not restrict files to `baseDir` if `relativePath` contains directory traversal sequences (`..`).

**Prevention:** Use `filepath.Rel` after `filepath.Clean` to strictly verify that target paths reside within the intended base directory before accessing local files.
