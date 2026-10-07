## 2025-05-18 - Avoid JSON Roundtrips for In-Memory Variable Lookups
**Learning:** `ReplaceVariables` used JSON marshaling/unmarshaling and `go-funk` reflection for simple struct property access, creating heavy allocations on spec reading. Direct struct lookup and `strconv.Quote` eliminated JSON roundtrip overhead.
**Action:** Prefer direct type-safe map/struct lookups and `strconv.Quote` for escaping simple string variables.
