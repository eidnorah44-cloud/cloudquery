## 2026-05-26 - Direct slice indexing in regexp.ReplaceAllFunc callbacks
**Learning:** `regexp.ReplaceAllFunc` passes the full matched `[]byte` slice to the callback. Calling `regexp.FindSubmatch(match)` inside the callback executes the regex search twice and allocates extra slices. Slicing fixed prefix/suffix bytes directly from `match` (e.g. `match[7 : len(match)-1]`) eliminates submatch search overhead and reduces heap allocations.
**Action:** Always slice directly from `match` in `ReplaceAllFunc` callbacks when prefix/suffix delimiters are fixed length.
