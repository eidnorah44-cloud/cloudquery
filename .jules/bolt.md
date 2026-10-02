## 2026-03-29 - Avoid converting Arrow Field Metadata to map unless mutations are required
**Learning:** `arrow.Field.Metadata.ToMap()` and `arrow.MetadataFrom()` allocate new maps and slices on every field during schema transformations. Checking for key presence via `field.Metadata.FindKey(...)` before converting to a map allows skipping allocation entirely for fields without metadata mutations, reducing CPU time by ~30-35% and allocations by ~35%.
**Action:** When transforming Arrow schemas or record fields, inspect metadata key presence with `FindKey` before allocating map copies.
