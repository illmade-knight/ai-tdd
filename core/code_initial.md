### initial code to pass tests

OK here is our initial code to pass the tests.
We've included an additional accumulator function for string
so first review the code and then add tests for the string generic

```go

func Reduce[A, B any](collection []A, f func(B, A) B, initialValue B) B {
	var result = initialValue
	for _, x := range collection {
		result = f(result, x)
	}
	return result
}

func sumInts(accumulator, currentItem int) int {
	// Basic sum. We can consider potential overflow later if required.
	return accumulator + currentItem
}

func sumStrings(accumulator, currentItem string) string {
	// Basic sum. We can consider potential overflow later if required.
	return accumulator + currentItem
}

```