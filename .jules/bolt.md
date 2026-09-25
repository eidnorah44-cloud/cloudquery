## 2026-09-25 - Contiguous Slice Allocation in Go Batch Transforms
**Learning:** Allocating 2D slices (`[][]any`) row-by-row in hot loops causes $N+1$ allocations per batch. Allocating a single underlying contiguous `[]any` slice of size $N \times M$ and assigning 3-index subslices (`data[low : high : max]`) reduces memory allocations to 2 per batch and prevents slice capacity overflow side-effects.
**Action:** Use contiguous single-allocation backing arrays with 3-index slicing when transforming matrix or 2D batch structures in Go.
