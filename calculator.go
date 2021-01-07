// Package calculator provides a library for simple calculations in Go.
package calculator

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

// Add takes two or more numbers and returns the result of adding them together.
func Add(a, b float64, extras ...float64) float64 {
	result := a + b
	for _, v := range extras {
		result += v
	}
	return result
}

// Subtract takes two or more numbers and returns the result of the
// initial number, after subtracting the remaining numbers from the
// first.
func Subtract(a, b float64, extras ...float64) float64 {
	return a - Add(0, b, extras...)
}

// Multiply takes two or more floats and returns the result of multiplying the first
// by the remaining parameters.
func Multiply(a, b float64, extras ...float64) float64 {
	result := a * b
	for _, v := range extras {
		result *= v
	}
	return result
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

// SquareRoot takes a single float and returns the square root as result,
// or an error if the input is invalid in some way.
func SquareRoot(a float64) (float64, error) {
	if a < 0 {
		// return an explicit zero by convention
		return 0, errors.New("invalid input (square root of -ve numbers is unimaginable)")
	}
	return math.Sqrt(a), nil
}

// Parse takes a single string and attempts to parse it into 3 values:
// float64, float64, string
// else returns an error
func Parse(input string) (float64, float64, string, error) {
	s := regexp.MustCompile(" +").Split(input, 3)
	if len(s) != 3 {
		return 0, 0, "", errors.New("invalid input (parameter length !3")
	}

	left, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid input (%v could not be parsed as float64)", s[0])
	}

	right, err := strconv.ParseFloat(s[2], 64)
	if err != nil {
		return 0, 0, "", fmt.Errorf("invalid input (%v could not be parsed as float64)", s[2])
	}
	operator := s[1]
	switch operator {
	case "*":
	case "/":
	case "+":
	case "-":
	default: // errors just fall out
		return 0, 0, "", fmt.Errorf("invalid input (%v could not be parsed as one of */+-)", operator)
	}
	return left, right, operator, nil

}
