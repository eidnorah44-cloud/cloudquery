## 2025-05-18 - Avoid heap allocations in stream redactors with pre-converted byte slices & fast path presence checks
**Learning:** `bytes.ReplaceAll` in Go stream output redactors always performs slice allocations if called blindly. Pre-converting secret strings to byte slices at configuration time and using `strings.Contains` / `bytes.Contains` for fast-path checks completely avoids allocations on non-matching log lines.
**Action:** Always check `strings.Contains` or `bytes.Contains` before calling string/bytes modification functions in hot streaming paths.
