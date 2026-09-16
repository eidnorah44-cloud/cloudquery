## 2026-03-30 - Path Traversal Prevention in AI CLI Helper Functions
**Vulnerability:** Filename arguments supplied by AI assistant tool calls (`create_spec_file`, `create_sql_file`, `cloudquery_test`) were directly concatenated into file write paths and shell commands without validation, allowing potential path traversal attacks (e.g., `../../file`).
**Learning:** Functions accepting filenames from external tools or AI tool arguments must strictly sanitize path input before performing disk reads/writes or executing subprocess commands.
**Prevention:** Always validate that user-provided filename strings do not contain directory separators (`/`, `\`) or traversal elements (`..`), and restrict them strictly to base filenames.
