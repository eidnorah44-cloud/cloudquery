## 2025-05-18 - Avoid Redundant Membership Checks During Slice Iteration

**Learning:** Iterating over a slice and checking `slices.Contains(slice, element)` inside the loop creates an unnecessary $O(N^2)$ operation, as the element being iterated over is already guaranteed to be in the slice.
**Action:** When mapping or filtering slice elements via map lookups, iterate directly and check map key existence (`if dest, ok := m[key]; ok`) rather than calling `slices.Contains`.
