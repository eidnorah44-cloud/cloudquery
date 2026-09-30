## 2025-05-18 - Atomic Map Update Pattern for Thread-Safe Maps
**Learning:** In thread-safe wrappers around maps (`ConcurrentMap`), replacing a read-then-write pattern (`Get` then `Add`) with a single conditional atomic operation (`SetIfGreater`) eliminates redundant mutex lock acquisitions, prevents redundant map hashing/lookups, and eliminates potential race conditions between read and write calls.
**Action:** Use single-lock compound operations (e.g. `SetIfGreater` or `Upsert`) when updating thread-safe map structures conditionally based on current map state.
