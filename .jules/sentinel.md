## 2026-03-31 - Overlapping Secret Redaction Order and Determinism
**Vulnerability:** Non-deterministic secret redaction order in Go map iteration when redacting both raw environment assignment strings (`KEY=VALUE`) and standalone secret values (`VALUE`).
**Learning:** Iterating over a Go map for string replacement (`bytes.ReplaceAll`) causes non-deterministic replacement order. If a shorter secret is replaced first, it can alter full assignment strings (e.g. `KEY=VALUE` becoming `KEY=KEY`) before the full string replacement rule runs.
**Prevention:** Always sort redactor targets by string length in descending order prior to executing replacement passes so longer/more specific replacement patterns match first.
