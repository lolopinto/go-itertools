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

func TestEnumerate(t *testing.T) {
	keys := make([]int, 0)
	vals := make([]int, 0)
	for k, v := range Enumerate(NewSeq(0, 1, 2, 3)) {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	assert.Equal(t, keys, vals)
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

func TestCombinations(t *testing.T) {
	assertSequenceMatch(t,
		Combinations([]string{"A", "B", "C", "D"}, 2),
		[][]string{[]string{"A", "B"}, []string{"A", "C"}, []string{"A", "D"}, []string{"B", "C"}, []string{"B", "D"}, []string{"C", "D"}},
	)
	assertSequenceMatch(t,
		Combinations([]int{0, 1, 2, 3}, 3),
		[][]int{[]int{0, 1, 2}, []int{0, 1, 3}, []int{0, 2, 3}, []int{1, 2, 3}},
	)
}

func TestCombinationsWithReplacement(t *testing.T) {
	assertSequenceMatch(t,
		CombinationsWithReplacement([]string{"A", "B", "C"}, 2),
		[][]string{[]string{"A", "A"}, []string{"A", "B"}, []string{"A", "C"}, []string{"B", "B"}, []string{"B", "C"}, []string{"C", "C"}},
	)
}

func TestCompress(t *testing.T) {
	assertSequenceMatch(t,
		Compress(NewSeq([]byte("ABCDEFG")...), []int{1, 0, 1, 0, 1, 1}),
		[]byte("ACEF"),
	)

	assertSequenceMatch(t,
		Compress(Count(), []bool{false, true, false, true, false, true}),
		[]int{1, 3, 5},
	)
}

func TestDropWhile(t *testing.T) {
	assertSequenceMatch(t,
		DropWhile(func(x int) bool { return x < 5 }, NewSeq(1, 4, 6, 3, 8)),
		[]int{6, 3, 8},
	)
}

func TestFilterFalse(t *testing.T) {
	assertSequenceMatch(t,
		FilterFalse(func(x int) bool { return x < 5 }, NewSeq(1, 4, 6, 3, 8)),
		[]int{6, 8},
	)
}

func TestGroupBy(t *testing.T) {
	want := map[string][]string{
		"A": []string{"A", "A", "A", "A"},
		"B": []string{"B", "B", "B"},
		"C": []string{"C", "C"},
		"D": []string{"D"},
	}

	for k, g := range GroupBy(NewSeq("A", "A", "A", "A", "B", "B", "B", "C", "C", "D")) {
		assertSequenceMatch(t, g, want[k])
	}
}

func TestSlice(t *testing.T) {
	assertSequenceMatch(t,
		Slice(NewSeq([]byte("ABCDEFG")...), 2, 4),
		[]byte("CD"),
	)

	assertSequenceMatch(t,
		Slice(NewSeq([]byte("ABCDEFG")...), 2, -1),
		[]byte("CDEFG"),
	)
}

func TestPairwise(t *testing.T) {
	want := []string{"AB", "BC", "CD", "DE", "EF", "FG"}

	var i int
	for x, y := range Pairwise(NewSeq([]byte("ABCDEFG")...)) {
		assert.Equal(t, x, want[i][0])
		assert.Equal(t, y, want[i][1])
		i++
	}
}
