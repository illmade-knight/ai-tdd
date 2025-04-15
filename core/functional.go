package core

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
