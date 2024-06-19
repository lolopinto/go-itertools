package itertools

import (
	"iter"
	"slices"
)

func NewSeq[T any](vals ...T) iter.Seq[T] {
	return OfSlice(vals)
}

func OfSlice[T any](vals []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < len(vals); i++ {
			if !yield(vals[i]) {
				return
			}
		}
	}
}

func Enumerate[T any](s iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		var i int
		for v := range s {
			if !yield(i, v) {
				return
			}
			i++
		}
	}
}

func Map[T any, U any](mapper func(T) U, s iter.Seq[T]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range s {
			if !yield(mapper(v)) {
				return
			}
		}
	}
}

func Take[T any](s iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(s)
		defer stop()

		for i := 0; i < n; i++ {
			v, ok := next()
			if !ok || !yield(v) {
				return
			}
		}
	}
}

func Chain[T any](seqs ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range seqs {
			seq(yield)
		}
	}
}

func Count() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func Cycle[T any](s iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for {
			for v := range s {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func Repeat[T any](val T, n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; n < 0 || i < n; i++ {
			if !yield(val) {
				return
			}
		}

	}
}

func Accumulate[T any](s iter.Seq[T], op func(T, T) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		var sum T
		for v := range s {
			sum = op(sum, v)
			if !yield(sum) {
				return
			}
		}
	}
}

func Batched[T any](s iter.Seq[T], n int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		batch := make([]T, 0, n)

		for v := range s {
			if len(batch) == n {
				if !yield(batch) {
					return
				}
				batch = make([]T, 0, n)
			}

			batch = append(batch, v)
		}

		if len(batch) > 0 {
			yield(batch)
		}
	}
}

func pick[T any](vals []T, indices []int) []T {
	out := make([]T, 0, len(indices))
	for _, i := range indices {
		out = append(out, vals[i])
	}
	return out
}

func Combinations[T any](vals []T, r int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if r > len(vals) {
			return
		}

		indices := make([]int, 0, r)
		for i := range r {
			indices = append(indices, i)
		}

		yield(pick(vals, indices))

		for {
			var i int
			var found bool
			for i = r - 1; i >= 0; i-- {
				if indices[i] != i+len(vals)-r {
					found = true
					break
				}
			}

			if !found {
				return
			}

			indices[i]++
			for j := i + 1; j < r; j++ {
				indices[j] = indices[j-1] + 1
			}

			yield(pick(vals, indices))
		}
	}
}

func CombinationsWithReplacement[T any](vals []T, r int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if len(vals) == 0 && r == 0 {
			return
		}

		indices := make([]int, r)

		yield(pick(vals, indices))
		for {
			var i int
			var found bool
			for i = r - 1; i >= 0; i-- {
				if indices[i] != len(vals)-1 {
					found = true
					break
				}
			}

			if !found {
				return
			}

			nextIndex := indices[i] + 1
			for j := i; j < r; j++ {
				indices[j] = nextIndex
			}

			yield(pick(vals, indices))
		}
	}
}

func Compress[T any, C comparable](s iter.Seq[T], selectors []C) iter.Seq[T] {
	return func(yield func(T) bool) {
		var zero C
		for i, v := range Enumerate(s) {
			if i >= len(selectors) {
				return
			}

			if selectors[i] != zero {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func DropWhile[T any](pred func(T) bool, s iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var shouldYield bool
		for v := range s {
			if !pred(v) {
				shouldYield = true
			}
			if shouldYield {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func FilterFalse[T any](pred func(T) bool, s iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !pred(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func GroupBy[T comparable](s iter.Seq[T]) iter.Seq2[T, iter.Seq[T]] {
	return func(yield func(T, iter.Seq[T]) bool) {
		next, stop := iter.Pull(s)
		defer stop()

		var current T
		var ok bool

		pullGroup := func(groupValue T) iter.Seq[T] {
			return func(yield func(T) bool) {
				if !yield(groupValue) {
					return
				}

				var v T
				for {
					v, ok = next()
					if !ok {
						return
					}

					if v != groupValue {
						current = v
						return
					}

					if !yield(v) {
						return
					}
				}
			}
		}

		current, ok = next()
		if !ok {
			return
		}

		for {
			group := pullGroup(current)
			if !ok {
				return
			}

			if !yield(current, group) {
				return
			}

			// exhaust remaining group items before moving to next group
			var v T
			for {
				v, ok = next()
				if !ok {
					return
				}
				if v != current {
					current = v
					return
				}
			}
		}
	}
}

func Slice[T any](s iter.Seq[T], start, end int) iter.Seq[T] {
	return func(yield func(T) bool) {
		var i int
		for v := range s {
			if i >= start && (end < 0 || i < end) {
				if !yield(v) {
					return
				}
			}
			i++
		}
	}
}

func Pairwise[T any](s iter.Seq[T]) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		next, stop := iter.Pull(s)
		defer stop()

		a, ok := next()
		if !ok {
			return
		}
		for {
			b, ok := next()
			if !ok {
				return
			}
			if !yield(a, b) {
				return
			}
			a = b
		}
	}
}

func Permutations[T any](vals []T, r int) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(vals)
		if r > n {
			return
		}

		indices := make([]int, 0, n)
		for i := range n {
			indices = append(indices, i)
		}

		cycles := make([]int, 0, r)
		for i := n; i > n-r; i-- {
			cycles = append(cycles, i)
		}

		yield(pick(vals, indices[:r]))

		if n == 0 {
			return
		}

		for {
			var i int
			var found bool
			for i = r - 1; i >= 0; i-- {
				cycles[i]--
				if cycles[i] == 0 {
					// move to end and reset
					ind := indices[i]
					indices = append(append(indices[:i], indices[i+1:]...), ind)
					cycles[i] = n - i
				} else {
					j := n - cycles[i]
					indices[i], indices[j] = indices[j], indices[i]

					yield(pick(vals, indices[:r]))
					found = true
					break
				}
			}

			if !found {
				return
			}
		}
	}
}

func Product[T any](pool ...[]T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		n := len(pool)

		maxIndices := make([]int, n)
		for i := range n {
			maxIndices[i] = len(pool[i]) - 1
		}

		indices := make([]int, n)

		for {
			prod := make([]T, n)
			for i := range n {
				prod[i] = pool[i][indices[i]]
			}

			if !yield(prod) {
				return
			}

			if slices.Equal(indices, maxIndices) {
				return
			}

			for i := n - 1; i >= 0; i-- {
				if indices[i] < maxIndices[i] {
					indices[i]++
					break
				} else {
					indices[i] = 0
				}
			}
		}
	}
}

func ProductRepeat[T any](vals []T, repeat int) iter.Seq[[]T] {
	inputs := make([][]T, repeat)
	for i := range repeat {
		inputs[i] = vals
	}
	return Product(inputs...)
}

func TakeWhile[T any](pred func(T) bool, s iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !pred(v) || !yield(v) {
				return
			}
		}
	}
}
