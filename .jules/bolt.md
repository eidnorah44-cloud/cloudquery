## 2026-09-23 - Fast-path exact type matching before regex parsing in plugin arrow type converters
**Learning:** Type conversion functions in plugin SDK / destination plugins (e.g. pgarrow, snowflake) were running `regexp.FindAllStringSubmatch` for parameterized types on every single column type string before checking standard type exact matches in `switch t`.
**Action:** Always place exact `switch t` matches first, use `strings.HasPrefix` guards before calling regex parsers, and prefer `FindStringSubmatch` over `FindAllStringSubmatch` when only the first match is needed.
