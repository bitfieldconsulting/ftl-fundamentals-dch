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


func TestEvaluate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		input       string
		res         float64
		errExpected bool
	}{
		{name: "send empty string ''", input: "", errExpected: true},
		{name: "send wrongly typed values", input: "wrongly typed values", errExpected: true},
		{name: "send far too many values", input: "1 2 3 4", errExpected: true},
		{name: "add 2 floats", input: "3 + 2", res: 5, errExpected: false},
		{name: "multiply 2 floats", input: "3 * 2", res: 6, errExpected: false},
		{name: "divide 2 floats", input: "3 / 2", res: 1.5, errExpected: false},
		{name: "subtract 2 floats", input: "3 - 2", res: 1, errExpected: false},
		{name: "parse sneaky floats", input: ".6 + -0.6", res: 0, errExpected: false},
	}

	for _, tc := range testCases {
		res, err := calculator.Evaluate(tc.input)

		if tc.errExpected == true {
			// expecting an error, so ignore results
			if err == nil {
				t.Errorf("%v: given params %v, wanted err, got nil", tc.name, tc.input)
			}
		} else {
			if tc.res != res {
				// unexpected error
				t.Errorf("%v: given params %v, got %v", tc.name, tc.input, res)
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
		{
			name: "add 2 floats",
			fn: calculator.Add,
			a: 1, b: 2,
			want: 3,
		},
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
		{name: "recurring decimal", a: 2, b: 3, want: 0.666667, errExpected: false},
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
		t.Run(tc.name, func(t *testing.T) {
			got, err := calculator.Divide(tc.a, tc.b)
			errReceived := err != nil
			if tc.errExpected != errReceived {
				t.Fatalf("unexpected error status %v", err)
			}
			if !tc.errExpected && !closeEnough(tc.want, got) {
				t.Errorf("given params %f, %f, wanted %f, got %f", tc.a, tc.b, tc.want, got)
			}
		})
	}
}

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < 0.000001
}

func TestSquareRoot(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		input       float64
		want        float64
		errExpected bool
	}{
		{name: "square root of 9", input: 9, want: 3},
		{name: "square root of 1", input: 1, want: 1},
		{name: "square root of 0", input: 0, want: 0},
		{name: "square root of fraction", input: 6.25, want: 2.5},
		{name: "square root of -1", input: -1, want: 3, errExpected: true}, // inconceivable!
	}

	for _, tc := range testCases {
		got, err := calculator.SquareRoot(tc.input)
		errReceived := err != nil
		if tc.errExpected != errReceived {
			t.Fatalf("unexpected error status %v", err)
		}
		if !tc.errExpected && tc.want != got {
			t.Errorf("given params %f, wanted %f, got %f", tc.input, tc.want, got)
		}
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
