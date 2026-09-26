## 2026-09-26 - Function Slice Allocation in Hot Parser Loops
**Learning:** Instantiating slices of function pointers or closures (e.g. `parsers := []func(...)`) inside frequently called functions causes slice and header allocations on every invocation, even when the functions are package-level identifiers.
**Action:** Call helper/parser functions directly or use package-level slice variables for static parser functions to achieve zero heap allocations on hot conversion paths.
