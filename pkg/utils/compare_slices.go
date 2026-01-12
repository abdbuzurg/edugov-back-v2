// Package utils provides small helpers shared across services.
package utils

// CompareSlices compares two slices by id.
//
// Rules:
//   - In slice2: id == 0 => new items (created)
//   - In slice2: id != 0 => existing items (updated/existing)
//   - Removed items are those in slice1 with id not present among slice2 existing ids.
//
// Assumptions:
//   - No duplicate non-zero IDs inside the same slice.
func CompareSlices[T any](
	slice1 []T,
	slice2 []T,
	getID func(T) int64,
) (updated []T, created []T, removed []T) {
	// Build a set of "existing" ids from slice2 (id != 0),
	// and split slice2 into created vs updated.
	slice2Existing := make(map[int64]struct{}, len(slice2))

	for _, it := range slice2 {
		id := getID(it)
		if id == 0 {
			created = append(created, it)
			continue
		}
		slice2Existing[id] = struct{}{}
		updated = append(updated, it)
	}

	// Removed = items in slice1 whose id not present in slice2Existing
	for _, it := range slice1 {
		id := getID(it)
		if id == 0 {
			// Typically slice1 shouldn't contain new items; ignore if it does.
			continue
		}
		if _, ok := slice2Existing[id]; !ok {
			removed = append(removed, it)
		}
	}

	return updated, created, removed
}
