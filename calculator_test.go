package calculator_test

import (
	"calculator"
	"testing"
)

func TestAddSubMul(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name   string
		fn     func(float64, float64, ...float64) float64
		a, b   float64
		extras []float64
		want   float64
	}{
		{name: "add 2 floats", fn: calculator.Add, a: 1, b: 2, want: 3},
		{name: "add several floats", fn: calculator.Add, a: 3, b: 4, extras: []float64{4.0}, want: 11},
		{name: "add only variadic floats", fn: calculator.Add, a: 0, b: 0, extras: []float64{1, 2, 3}, want: 6},
		{name: "add fractional floats", fn: calculator.Add, a: 0.25, b: 0.5, extras: []float64{0.25}, want: 1},

		{name: "subtract floats", fn: calculator.Subtract, a: 4, b: 2, want: 2},
		{name: "subtract variadic floats", fn: calculator.Subtract, a: 4, b: 2, extras: []float64{0.5, -0.5}, want: 2},

		{name: "multiply floats", fn: calculator.Multiply, a: 4, b: 2, want: 8},
		{name: "multiply floats with a zero", fn: calculator.Multiply, a: 4, b: 0, want: 0},
		{name: "multiply variadic floats with negative fractions", fn: calculator.Multiply, a: 4, b: 2, extras: []float64{0.5, -0.5}, want: -2},
	}

	for _, tc := range testCases {
		got := tc.fn(tc.a, tc.b, tc.extras...)
		if tc.want != got {
			t.Errorf("%v: given params %f,%f, and slice %v want %f, got %f", tc.name, tc.a, tc.b, tc.extras, tc.want, got)
		}
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
