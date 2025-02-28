package basic

import "fmt"

// Function to be tested (variadic sum)
func Sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

// SumMixed function now accepts a mix of ints and []int
func SumMixed(inputs ...interface{}) (int, error) {
	total := 0
	for _, input := range inputs {
		switch v := input.(type) {
		case int:
			total += v
		case []int:
			for _, num := range v {
				total += num
			}
		default:
			return 0, fmt.Errorf("invalid input type: %T", input) // Return error for invalid type
		}
	}
	return total, nil
}
