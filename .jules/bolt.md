# Bolt's Journal

## 2026-03-31 - Package-level pre-compilation of regular expressions
**Learning:** Re-compiling regular expressions inside frequently called helper functions (such as `extractYamlFromMarkdownCodeBlock` in CLI initialization) causes unnecessary string parsing and heap allocation overhead on every invocation.
**Action:** Always pre-compile regular expressions at package level using `var re = regexp.MustCompile(...)` when the pattern is static.
