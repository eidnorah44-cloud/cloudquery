## 2026-03-31 - Avoid `time.After` in high-throughput Go streaming select loops
**Learning:** In streaming pipelines like `TransformerPipeline`, using `time.After(duration)` inside `for { select { ... } }` loops allocates a new `time.Timer` on every iteration (every record processed). These timers remain active in Go's runtime timer heap until they expire, leading to high allocation rate, heap fragmentation, and GC churn.
**Action:** Always replace `time.After` in high-frequency event loops with `ticker := time.NewTicker(duration)` and `defer ticker.Stop()` outside the loop.
