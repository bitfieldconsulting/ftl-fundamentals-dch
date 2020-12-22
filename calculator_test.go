package calculator_test

import (
	"calculator"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()
	var want float64 = 4
	got := calculator.Add(2, 2)
	if want != got {
		t.Errorf("want %f, got %f", want, got)
	}
}

func TestSubtract(t *testing.T) {
	t.Parallel()
	var want float64 = 2
	got := calculator.Subtract(4, 2)
	if want != got {
		t.Errorf("want %f, got %f", want, got)
	}
}

func TestMultiply(t *testing.T) {
	t.Parallel()
	var want float64 = 12
	got := calculator.Multiply(3, 4)
	if want != got {
		t.Errorf("want %f, got %f", want, got)
	}
}

func TestDivide(t *testing.T) {
	t.Parallel()
	var want float64 = 7
	got, _ := calculator.Divide(63, 9)
	if want != got {
		t.Errorf("want %f, got %f", want, got)
	}
}

func TestDivideByZero(t *testing.T) {
	t.Parallel()
	// one does not simply assign a literal infinite
	// it must be summoned from the darkest depths of stdlib
	// and it comes signed
	// https: //golang.org/test/zerodivide.go
	// var want float64 = math.Inf(+1)
	// got := calculator.Divide(1, 0)
	// if want != got {
	// 	t.Errorf("want %f, got %f", want, got)
	// }
}
