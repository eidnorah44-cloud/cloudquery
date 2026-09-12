## 2025-05-18 - Fast-path & Direct Lookup for Variable Interpolation
**Learning:** In Go spec parsing, converting structs to JSON maps (`json.Marshal` + `json.Unmarshal`) for dynamic property lookup (via `go-funk`) created massive overhead (25-54 allocations per call, 4.4-9.8 µs).
**Action:** Always provide a `strings.Contains` fast-path check to bypass regex and lookup overhead when no variable markers exist, and use direct struct lookups before falling back to generic map reflection.
