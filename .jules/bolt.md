## 2026-03-31 - Avoid `FindAllStringSubmatch` and parser closure slices in type converters

**Learning:** `regexp.FindAllStringSubmatch(s, -1)` allocates nested 3D string slices (`[][][]string`) to capture all occurrences in string `s`, even when only the first match is needed. In schema/type conversion functions executed per column, replacing `FindAllStringSubmatch` with `FindStringSubmatch` and replacing dynamic closure slice allocations (`parsers := []func...`) with direct function calls reduces memory allocations by ~80% (from 956 B/op to 184 B/op) and decreases latency by ~39%.

**Action:** Look for `FindAllStringSubmatch` or transient closure slices in hot paths (like type parsers or column schema converters) and replace them with `FindStringSubmatch` and direct function calls.
