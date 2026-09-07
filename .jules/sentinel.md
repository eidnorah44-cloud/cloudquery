## 2026-03-30 - Path Traversal in LLM Function Call Arguments
**Vulnerability:** AI CLI assistant (`cli/cmd/ai.go`) executed file writes (`create_spec_file`, `create_sql_file`) and shell execution (`cloudquery_test`) using unvalidated `filename_without_extension` parameters from AI function calls.
**Learning:** Tool/function call inputs from LLM integrations must be treated as untrusted input just like web user input, as malicious prompts or crafted responses could specify relative or absolute path traversal sequences (e.g. `../../filename`).
**Prevention:** Always sanitize tool call filenames using `filepath.Clean` and `filepath.Base` to ensure file operations stay confined to base filenames within the target directory.
