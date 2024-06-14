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
