# Sentinel Security Journal

## 2025-05-18 - Path Traversal in Plugin Docs Download
**Vulnerability:** In `cli/cmd/plugin_docs_download.go`, plugin documentation page names returned by API endpoints were joined to local directories using `strings.ReplaceAll(item.Name, string(filepath.Separator), "_")`. On Windows, `/` in names was ignored; on Unix, `\` was ignored, and path traversal sequences (`../` or `..\`) could allow arbitrary file creation outside the target `docsDir`.
**Learning:** Replacing only OS-native `filepath.Separator` is insufficient when dealing with inputs from remote servers or foreign OS paths.
**Prevention:** Normalize all separators (`\`) to `/`, sanitize using `filepath.Base(filepath.Clean(...))`, and strip or convert all `/` and `\` characters before joining with local directories.
