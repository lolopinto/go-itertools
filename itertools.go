package itertools

import (
	"iter"
)

func NewSeq[T any](vals ...T) iter.Seq[T] {
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

func Combinations[T any](vals []T, r int) iter.Seq[[]T] {
	pick := func(indices []int) []T {
		out := make([]T, 0, len(indices))
		for _, i := range indices {
			out = append(out, vals[i])
		}
		return out
	}

	return func(yield func([]T) bool) {
		if r > len(vals) {
			return
		}

		indices := make([]int, 0, r)
		for i := range r {
			indices = append(indices, i)
		}

		yield(pick(indices))

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

			yield(pick(indices))
		}
	}
}

func CombinationsWithReplacement[T any](vals []T, r int) iter.Seq[[]T] {
	pick := func(indices []int) []T {
		out := make([]T, 0, len(indices))
		for _, i := range indices {
			out = append(out, vals[i])
		}
		return out
	}

	return func(yield func([]T) bool) {
		if len(vals) == 0 && r == 0 {
			return
		}

		indices := make([]int, r)

		yield(pick(indices))
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

			yield(pick(indices))
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
