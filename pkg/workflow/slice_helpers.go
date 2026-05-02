// This file provides generic slice utilities.
//
// These utilities are used throughout the workflow compilation process
// to safely manipulate string slices. They complement the map utilities
// in map_helpers.go.
//
// # Key Functions
//
// Slice Operations:
//   - mergeUnique() - Merge two slices while preserving order and deduplicating
//   - excludeFromSlice() - Create a new slice excluding specified values

package workflow

// excludeFromSlice returns a new slice containing the items from base
// that do not appear in the exclude set. Order of remaining items is preserved.
// Always returns a fresh slice (never aliases base) even when no items are removed.
func excludeFromSlice(base []string, exclude ...string) []string {
	if len(exclude) == 0 {
		return append([]string(nil), base...)
	}
	excluded := make(map[string]bool, len(exclude))
	for _, v := range exclude {
		excluded[v] = true
	}
	result := make([]string, 0, len(base))
	for _, v := range base {
		if !excluded[v] {
			result = append(result, v)
		}
	}
	return result
}

// mergeUnique returns a deduplicated slice that starts with base and appends any
// items from extra that are not already present in base.  Order is preserved.
func mergeUnique(base []string, extra ...string) []string {
	seen := make(map[string]bool, len(base)+len(extra))
	result := make([]string, 0, len(base)+len(extra))
	for _, v := range base {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	for _, v := range extra {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
