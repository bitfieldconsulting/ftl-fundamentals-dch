package calculator_test

import (
	"calculator"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Parallel()
	type testCase struct {
		a, b float64
		want float64
	}
	testCases := []testCase{
		{a: 2, b: 2, want: 4},
	}

	for _, tc := range testCases {
		got := calculator.Add(tc.a, tc.b)
		if tc.want != got {
			t.Errorf("want %f, got %f", tc.want, got)
		}
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
	// the value that go returns internally for divide by zero is a NaN
	// +infinity, which can be generated with var float64 = math.Inf(+1)
	_, err := calculator.Divide(1, 0)
	// fail test unless error is not nil
	if err == nil {
		t.Error("wanted err, got nil")
	}
}
