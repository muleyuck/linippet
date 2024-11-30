package helper

import "iter"

// Filter returns an iterator over seq that only includes
// the values v for which f(v) is true.
func FilterSlice[V any](seq iter.Seq[V], f func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}
