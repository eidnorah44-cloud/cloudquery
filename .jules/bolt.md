## 2026-03-30 - Fast-path ASCII key validation for JSON sanitization
**Learning:** In Go record processing pipelines, using `regexp.ReplaceAllString` for string sanitization (e.g. `\W` key sanitization in S3 Athena mode) allocates memory and runs regex engine for every string key even when keys are valid identifiers. Replacing it with a fast-path scanning function avoids all allocations when keys are valid.
**Action:** Replace `regexp` string replacement in hot data processing loops with custom fast-path range/scanner functions.
