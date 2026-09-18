## 2026-03-02 - Avoid Unconditional Metadata Conversions in Arrow Schema Transformations
**Learning:** Calling `field.Metadata.ToMap()` and `arrow.MetadataFrom()` for every field in schema transformations causes unnecessary `map[string]string` allocations and slice copies when metadata doesn't actually need mutation. Using `FindKey` beforehand to fast-path unmodified fields drastically cuts memory allocations.
**Action:** Always check if metadata keys (`FindKey`) actually exist before allocating maps to mutate Arrow field metadata.
