package itertools

import (
	"iter"
	"testing"
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

// TODO (astonm): Figure out how to make assertMatch work generically
func assertMatch(t *testing.T, gotSeq iter.Seq[int], want []int) {
	got := toSlice(gotSeq)

	if len(want) != len(got) {
		t.Fatalf("expected %v, got %v", want, got)
	}

	for i := 0; i < len(got); i++ {
		if want[i] != got[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestNewSeq(t *testing.T) {
	assertMatch(t, NewSeq(1, 2, 3), []int{1, 2, 3})
}
