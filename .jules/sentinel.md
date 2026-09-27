## 2025-05-18 - Snowflake Identifier Escaping in Dynamic SQL

**Vulnerability:** Snowflake dynamic SQL query clauses (`COPY INTO`, `MERGE INTO`, `UPDATE SET`, `INSERT`) were constructed using `fmt.Sprintf` with raw table and column names without double-quote escaping or identifier sanitization. Furthermore, column sanitization previously only upper-cased names without escaping internal double quotes, allowing SQL injection via crafted table/column identifiers (e.g., `col"; DROP TABLE users; --`).
**Learning:** In Go database drivers where `identifier(?)` parameters cannot be used for all query fragments (such as column projections or `MERGE` clause aliases), identifiers wrapped in double quotes must have internal double quotes escaped as `""` (`strings.ReplaceAll(name, `"`, `""`)`).
**Prevention:** Always route table and column names in dynamic SQL through a sanitization function that escapes embedded double-quote characters before wrapping in double quotes.
