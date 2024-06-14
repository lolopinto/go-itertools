package main

import "iter"

func Count() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i + 1) {
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
