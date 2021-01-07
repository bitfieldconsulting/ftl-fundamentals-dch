package calculator_test

import (
	"calculator"
	"fmt"
	"math"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestParse(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		input       string
		left        float64
		right       float64
		operator    string
		errExpected bool
	}{
		{name: "send empty string ''", input: "", left: 0, right: 0, operator: "", errExpected: true},
		{name: "send wrongly typed values", input: "wrongly typed values", left: 0, right: 0, operator: "", errExpected: true},
		{name: "send far too many values", input: "1 2 3 4", left: 0, right: 0, operator: "", errExpected: true},
		{name: "add 2 floats", input: "3 + 2", left: 3, right: 2, operator: "+", errExpected: false},
		{name: "multiply 2 floats", input: "3 * 2", left: 3, right: 2, operator: "*", errExpected: false},
		{name: "divide 2 floats", input: "3 / 2", left: 3, right: 2, operator: "/", errExpected: false},
		{name: "subtract 2 floats", input: "3 - 2", left: 3, right: 2, operator: "-", errExpected: false},
	}

	for _, tc := range testCases {
		left, right, operator, err := calculator.Parse(tc.input)
		// fmt.Printf("%v: want %v, got %v %v %v, err %v\n", tc.name, tc.input, left, right, operator, err)

		if tc.errExpected == false {
			if (tc.left != left) || (tc.right != right) || (tc.operator != operator) {
				// unexpected error
				t.Errorf("%v: given params %v, got %v %v %v, err %v", tc.name, tc.input, left, operator, right, err)
			}
		} else {
			// expecting an error, so ignore all resultrs
			if err == nil {
				t.Errorf("%v: given params %v, wanted err, got nil", tc.name, tc.input)
			}
		}
	}
}

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
	testCases := []struct {
		name        string
		a, b        float64
		want        float64
		errExpected bool
	}{
		{name: "divide integers", a: 2, b: 2, want: 1, errExpected: false},
		{name: "divide integers", a: 63, b: 9, want: 7, errExpected: false},
		// use math.NaN() here as if we use an actual zero, flip the
		// test to ... test the test ... it doesn't fail. we explicitly
		// use NaN to ensure that we are testing the test we want to
		// test, and a value that cannot pass the other leg.
		// one does not simply assign a literal infinite
		// it must be summoned from the darkest depths of stdlib
		// and it comes signed
		// https: //golang.org/test/zerodivide.go
		// the value that go returns internally for divide by zero is a NaN
		// +infinity, which can be generated with var float64 = math.Inf(+1)
		{name: "divide by zero", a: 63, b: 0, want: math.NaN(), errExpected: true},
	}

	for _, tc := range testCases {
		got, err := calculator.Divide(tc.a, tc.b)
		if tc.errExpected == false {
			if tc.want != got {
				t.Errorf("%v: given params %f, %f, wanted %f, got %f", tc.name, tc.a, tc.b, tc.want, got)
			}
		} else {
			// expecting an error, so tc.want is not important
			if err == nil {
				t.Errorf("%v: given params %f, %f, wanted err, got nil", tc.name, tc.a, tc.b)
			}
		}
	}
}

func TestSquareRoot(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name string
		a    float64
		want float64
	}{
		{name: "square root of 9", a: 9, want: 3},
		{name: "square root of 1", a: 1, want: 1},
		{name: "square root of 0", a: 0, want: 0},
		{name: "square root of fraction", a: 6.25, want: 2.5},
	}

	for _, tc := range testCases {
		got, _ := calculator.SquareRoot(tc.a)
		if tc.want != got {
			t.Errorf("%v: given param %f, wanted %f, got %f", tc.name, tc.a, tc.want, got)
		}
	}
}

func TestNegativeSquareRoot(t *testing.T) {
	t.Parallel()
	_, err := calculator.SquareRoot(-1)
	// fail test unless error is not nil
	if err == nil {
		t.Error("wanted err, got nil")
	}
}

func TestProperties(t *testing.T) {
	// t.Parallel()
	// initialise our property testing
	parameters := gopter.DefaultTestParameters()
	// use a fixed seed so others can see the failure
	parameters.Rng.Seed(2000)
	// just enough tests to hit the first failure case
	parameters.MinSuccessfulTests = 4000
	properties := gopter.NewProperties(parameters)
	// define a simple property
	properties.Property("", prop.ForAll(
		func(n float64) bool {
			a := calculator.Add(n, n)
			b := calculator.Multiply(n, 2)
			c, _ := calculator.Divide(b, 2)
			if (a == b) && (c == n) {
				return true
			} else {
				fmt.Printf("failed %v, got %v, %v, %v\n", n, a, b, c)
				return false
			}
		},
		gen.Float64Range(1, 100000),
		// gen.Float64(),
	))
	properties.TestingRun(t)
}
