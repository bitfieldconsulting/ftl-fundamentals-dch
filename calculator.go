// Package calculator provides a library for simple calculations in Go.
package calculator

import "errors"

// Add takes two numbers and returns the result of adding them together.
func Add(a, b float64, extras ...float64) float64 {
	result := a + b
	for _, v := range extras {
		result += v
	}
	return result
}

// Subtract takes two numbers and returns the result of subtracting the second
// from the first.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply takes two floats and returns the result of multiplying the first
// by the second.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide takes two floats and returns the result of dividing the first
// by the second, or an error if the input is invalid, for example the
// second operand is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		// we need to return something of type /error/ and we want to
		// make it explicit and meaningful.
		// return an explicit zero by convention
		return 0, errors.New("invalid input (divide by zero is undefined)")
	}
	return a / b, nil
}
