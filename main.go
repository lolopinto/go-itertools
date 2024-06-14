package main

import (
	"fmt"
	"iter"
)

func main() {
	for x := range Take(Count(), 10) {
		fmt.Println(x)
	}
	fmt.Println()

	for y := range ToSeq([]int{5, 4, 3, 2, 1}) {
		fmt.Println(y)
	}
	fmt.Println()

	together := Chain[int](Take(Count(), 5), ToSeq([]int{5, 4, 3, 2, 1}))
	for a := range Take(together, 7) {
		fmt.Println(a)
	}
}

func Count() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; ; i++ {
			if !yield(i + 1) {
				return
			}
		}
	}
}

func ToSeq[T any](vals []T) iter.Seq[T] {
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
