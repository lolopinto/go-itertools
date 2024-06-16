package itertools

import (
	"iter"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toSlice[T any](s iter.Seq[T]) []T {
	maxLen := 20
	next, _ := iter.Pull(s)
	out := make([]T, 0)
	for {
		if len(out) > maxLen {
			panic("iterator ran too long")
		}

		v, ok := next()
		if !ok {
			break
		}
		out = append(out, v)
	}
	return out
}

func assertSequenceMatch[V any](t *testing.T, gotSeq iter.Seq[V], want []V) {
	got := toSlice(gotSeq)
	assert.Equal(t, want, got)
}

func TestNewSeq(t *testing.T) {
	assertSequenceMatch(t, NewSeq(1, 2, 3), []int{1, 2, 3})
}

func TestTake(t *testing.T) {
	assertSequenceMatch(t, Take(NewSeq(1, 2, 3), 0), []int{})
	assertSequenceMatch(t, Take(NewSeq(1, 2, 3, 4, 5, 6), 3), []int{1, 2, 3})
}

func TestChain(t *testing.T) {
	assertSequenceMatch(t,
		Chain(NewSeq(1, 2, 3), NewSeq(4, 5, 6)),
		[]int{1, 2, 3, 4, 5, 6},
	)
}

func TestCount(t *testing.T) {
	assertSequenceMatch(t, Take(Count(), 3), []int{0, 1, 2})
}

func TestCycle(t *testing.T) {
	assertSequenceMatch(t, Take(Cycle(NewSeq(1, 2, 3)), 5), []int{1, 2, 3, 1, 2})
}

func TestRepeat(t *testing.T) {
	assertSequenceMatch(t, Repeat("a", 5), []string{"a", "a", "a", "a", "a"})
}

func TestAccumulate(t *testing.T) {
	runningSums := Accumulate(NewSeq(1, 2, 3), func(x, y int) int { return x + y })
	assertSequenceMatch(t, runningSums, []int{1, 3, 6})
}

func TestBatched(t *testing.T) {
	assertSequenceMatch(t,
		Batched(NewSeq(1, 2, 3, 4, 5, 6, 7, 8), 3),
		[][]int{[]int{1, 2, 3}, []int{4, 5, 6}, []int{7, 8}},
	)
}
