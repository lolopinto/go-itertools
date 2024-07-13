package main

import (
	"fmt"
	"iter"

	"github.com/astonm/go-itertools"
)

func main() {

	// test_seq()
	// test_ofslice()
	// test_enumerate()
	// test_map()
	// test_take()
	// test_chain()
	// test_count()
	// test_cycle()
	// test_repeat()
	// test_accumulate()
	// test_batched()
	// test_combinations()
	// test_combinations_with_replacements()
	// test_compress()
	// test_drop_while()
	// test_filter_false()
	// test_group_by()
	// test_slice()
	// test_pairwise()
	// test_permutations()
	// test_product()
	// test_product_repeat()
	// test_take_while()
	// test_tee()
	// test_zip()
	// test_pullzip3()
	test_pullzip4()
}

func get_fruits() []string {
	return []string{"apple", "banana", "cherry", "date", "elderberry", "fig", "grape"}
}

func get_numbers() []int {
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
}

func get_romans() []string {
	return []string{"I", "II", "III", "IV", "V", "VI", "VII"}
}

func get_starks() []string {
	return []string{
		"Rob",
		"Sansa",
		"Jon",
		"Arya",
		"Bran",
		"Rickon",
	}
}

func test_seq() {
	// i like this API better than OfSlice but see why both are needed
	// not sure what to do with it in this dumb example tho lol
	// would like a range function that returns a Seq. would make some of these more powerful
	v := itertools.NewSeq(1, 2, 3, 4, 5)
	for x := range v {
		println(x)
	}
}

func test_ofslice() {
	v := itertools.OfSlice(get_fruits())
	for x := range v {
		println(x)
	}
}

func test_enumerate() {
	v := itertools.Enumerate(itertools.OfSlice(get_fruits()))
	for x, y := range v {
		println(x, y)
	}
}

func mapped_fruits(fns ...func(x string) string) iter.Seq[string] {
	if len(fns) > 1 {
		panic("only support one function for now")
	}
	if len(fns) == 0 {
		return itertools.Map(func(x string) string { return x + "!" }, itertools.OfSlice(get_fruits()))
	}

	return itertools.Map(fns[0], itertools.OfSlice(get_fruits()))
}

func mapped_fruits_slice(fns ...func(x string) string) []string {
	var l []string
	for v := range mapped_fruits(fns...) {
		l = append(l, v)
	}
	return l
}

func test_map() {
	// looks consistent with starmap but it seems to make more sense for the slice to be first
	v := mapped_fruits()
	for x := range v {
		println(x)
	}
}

func test_take() {
	// would be nice if some of these could just take a slice and convert it to a Seq as needed
	// the need to call OfSlice is a bit annoying
	v := itertools.Take(itertools.OfSlice(get_fruits()), 3)
	for x := range v {
		println(x)
	}

	println("")

	v2 := itertools.Take(itertools.OfSlice(get_fruits()), 10)
	for x := range v2 {
		println(x)
	}
}

func test_chain() {
	// another place where it's annoying to have to call OfSlice
	v := itertools.Chain(itertools.OfSlice(get_fruits()), mapped_fruits(), mapped_fruits(func(x string) string { return x + "!!" }), mapped_fruits(func(x string) string { return x + "!!!" }))
	for x := range v {
		println(x)
	}
}

func test_count() {
	// not clear exactly when this would be useful but makes sense
	// seems like it needs start, stop and a way to do something like this

	//  counter = itertools.count(start=10, step=2)

	// first_ten = [next(counter) for _ in range(10)]

	v := itertools.Count()
	for x := range v {
		if x > 10 {
			break
		}
		println(x)
	}
}

func test_cycle() {

	cycle := itertools.Cycle(itertools.OfSlice(get_fruits()))
	apple := 0
	for x := range cycle {
		if x == "apple" {
			apple++
		}
		if apple > 2 {
			break
		}
		println(x)
	}
}

func test_repeat() {
	v := itertools.Repeat("hello", 3)
	for x := range v {
		println(x)
	}
}

func test_accumulate() {
	fns := []func(x int, y int) int{func(x, y int) int { return x + y }, func(x, y int) int { return x * y }}

	for _, fn := range fns {
		// i would like the multiply function to work here
		v := itertools.Accumulate(itertools.OfSlice([]int{1, 2, 3, 4, 5}), fn)

		for x := range v {
			println(x)
		}
		println("")
	}
}

// this feels like it needs a better name. not sure lol
func test_batched() {
	v := itertools.Batched(itertools.OfSlice(get_fruits()), 3)

	for x := range v {
		fmt.Printf("%v\n", x)
	}
}

func test_combinations() {
	// why does this one take a slice and not a sequence??
	// seems inconsistent.
	// to be clear, i like the slice since i have earlier feedback about this but confusing...
	numbers := []int{2, 3, 4}
	for _, num := range numbers {
		combos := itertools.Combinations(get_fruits(), num)

		for x := range combos {
			fmt.Printf("%v\n", x)
		}
		println("")
	}
}

// i don't fully understand this one
// it's like the above but adds current item?
func test_combinations_with_replacements() {
	numbers := []int{2}
	for _, num := range numbers {
		combos := itertools.CombinationsWithReplacement(get_fruits(), num)

		for x := range combos {
			fmt.Printf("%v\n", x)
		}
		println("")
	}
}

func test_compress() {
	fruits := itertools.OfSlice(get_fruits())
	// would like something that can fill this automatically. not sure what that would look like
	// itertools.Count(length) -> fill if boolean passes
	// a list comprehension would do this easily
	selectors := []bool{true, false, true, false, true, false, true}

	v := itertools.Compress(fruits, selectors)
	for x := range v {
		println(x)
	}
}

func test_drop_while() {
	fruits := itertools.OfSlice(get_fruits())

	v := itertools.DropWhile(func(x string) bool { return x != "date" }, fruits)
	for x := range v {
		println(x)
	}

	println("")

	// idk what's happening here
	v2 := itertools.DropWhile(func(x string) bool { return x > "date" }, fruits)
	for x := range v2 {
		println(x)
	}
}

func test_filter_false() {
	fruits := itertools.OfSlice(get_fruits())

	// i'd kind also want something like this that can take an index
	// can this be combined?
	v := itertools.FilterFalse(func(x string) bool { return x == "date" }, fruits)
	for x := range v {
		println(x)
	}
}

func test_group_by() {
	// i'd like to see an example of this
}

func test_slice() {
	fruits := itertools.OfSlice(get_fruits())
	v := itertools.Slice(fruits, 2, 10)

	for x := range v {
		println(x)
	}
}

func test_pairwise() {
	// great! OfSlice still annoying but getting used to it
	fruits := itertools.OfSlice(get_fruits())
	for x, y := range itertools.Pairwise(fruits) {
		println(x, y)
	}
}

func test_permutations() {
	// fruits := itertools.OfSlice(get_fruits())
	// why is this a slice?
	for perms := range itertools.Permutations(get_fruits(), 3) {
		// would like a print function!
		// println(perms)
		// var l []string
		// for _, v := range perms {
		//  l = append(l, v)
		// }
		// println(strings.Join(l, ", "))

		fmt.Printf("%v\n", perms)

	}
}

func test_product() {
	for v := range itertools.Product(get_fruits(), mapped_fruits_slice(func(x string) string { return x + "!!!" })) {
		fmt.Printf("%v\n", v)
	}
}

func test_product_repeat() {
	// idk what this does
	for v := range itertools.ProductRepeat(get_fruits(), 5) {
		fmt.Printf("%v\n", v)
	}
}

func test_take_while() {
	// just made this the same as drop_while to compare
	fruits := itertools.OfSlice(get_fruits())

	v := itertools.TakeWhile(func(x string) bool { return x != "date" }, fruits)
	for x := range v {
		println(x)
	}

	println("")

	// idk what's happening here
	v2 := itertools.TakeWhile(func(x string) bool { return x > "date" }, fruits)
	for x := range v2 {
		println(x)
	}
}

func test_tee() {
	// idk what this is doing??
	// why is this printing ints?
	fruits := itertools.OfSlice(get_fruits())

	for v := range itertools.Tee(fruits, 2) {
		println(v)
	}
}

func test_zip() {
	fruits := itertools.OfSlice(get_fruits())
	numbers := itertools.OfSlice(get_numbers())
	for v1, v2 := range itertools.Zip(fruits, numbers) {
		println(v1, v2)
	}
}

func test_pullzip3() {
	fruits := itertools.OfSlice(get_fruits())
	numbers := itertools.OfSlice(get_numbers())
	romans := itertools.OfSlice(get_romans())

	// why return 2 here?
	// this API took me a while to figure this out lol
	// is there a simpler way to do this?
	fn1, _ := itertools.PullZip3(fruits, numbers, romans)
	for true {
		v1, v2, v3, b := fn1()
		println(v1, v2, v3, b)
		if !b {
			break
		}
	}
}

func test_pullzip4() {
	fruits := itertools.OfSlice(get_fruits())
	numbers := itertools.OfSlice(get_numbers())
	romans := itertools.OfSlice(get_romans())
	starks := itertools.OfSlice(get_starks())

	fn1, _ := itertools.PullZip4(fruits, numbers, romans, starks)
	for true {
		v1, v2, v3, v4, b := fn1()
		println(v1, v2, v3, v4, b)
		if !b {
			break
		}
	}

}
