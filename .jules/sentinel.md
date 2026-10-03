## 2026-03-31 - Path Traversal in Documentation Image Processing
**Vulnerability:** Markdown image parser in `cli/internal/publish/images/images.go` accepted relative file paths (e.g. `![](../secret.txt)`) without validating that the resolved path remained inside the document root directory (`docDir`).
**Learning:** Joining paths using `filepath.Join` or `filepath.IsAbs` alone does not prevent directory traversal if relative path components like `..` are present.
**Prevention:** Always use `filepath.Rel` against the clean base directory and check for `..` prefixes or `..` equality after resolving local file paths.
