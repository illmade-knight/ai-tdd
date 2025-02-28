## Shapes

we want to look at basic shapes: Rectangle, Circle and Triangle

we want tests for area and perimeter 

let's try a prompt

```aiprompt
write a golang test for basic 2D shapes: circle, rectangle, triangle, that tests for the perimeter and area. 
Do not include an implementation.
```

the result is in [shapes_v1_test](shapes_v1_test.go) (I've commented out the structs as we'll implement them in another file)

it's ok but not very comprehensive, let's increment

```aiprompt
add triangle test cases for colinear and some well known triangle edge lengths, isosceles, right angle etc
you could also expand the tests for Circle and Rectangle - for non int inputs, edge cases etc
```

the result is in [shapes_test](shapes_test.go)

I also add the prompt:

```aiprompt
we need something to check for almost equal values due to rounding errors
```

which gives us back a function:

```go
package shapes 

import "math"

func almostEqual(a, b float64) bool {
	const tolerance = 1e-6 // Adjust the tolerance as needed
	return math.Abs(a-b) < tolerance
}
```

and an explanation:

The almostEqual function compares two float64 values, a and b, and returns true if the absolute difference between them is less than the specified tolerance. You can adjust the tolerance value based on the precision you need for your comparisons.

Now, the test suite uses almostEqual to compare the calculated area and perimeter with the expected values, which should make the tests more robust and reliable when dealing with floating-point numbers.

### interesting
this is already interesting - my thinking for the tests has already expanded beyond the source material.

The expected Triangle is more general than the one in the gitbook and to pass the tests I've used a simple Heron's formula implementation 

we'll move onto [concurrency](concurrency) before doing something more complex in [clock](clock)
