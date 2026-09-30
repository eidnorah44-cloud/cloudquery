# Sentinel Security Journal

## 2025-05-18 - Path Traversal in AI assistant file creation and test execution
**Vulnerability:** In `cli/cmd/ai.go`, AI tool calls `create_spec_file`, `create_sql_file`, and `cloudquery_test` took an unsanitized `filename_without_extension` argument directly from API parameters, allowing path traversal (e.g., `../../something`) when writing files or calling `cloudquery test`.
**Learning:** `filepath.Base` or path validation should be used when receiving file names or file name prefixes from external/API sources.
**Prevention:** Sanitize or strip path components using `filepath.Base` before combining file path extensions or passing to file system operations.
