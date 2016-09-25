package gcd

import (
	"testing"
)

type Pair struct {
	a uint64
	b uint64
}

func TestSwapUint64(t *testing.T) {
	var tests = []struct {
		input    Pair
		expected Pair
	}{
		{
			input: Pair{
				a: 1000,
				b: 2222,
			},
			expected: Pair{
				a: 2222,
				b: 1000,
			},
		},
	}

	for _, test := range tests {
		swap(&test.input.a, &test.input.b)
		if test.input.a != test.expected.a || test.input.b != test.expected.b {
			t.Errorf("Actual: %#v, Expected: %#v", test.input, test.expected)
		}
		t.Logf("%#v", test.input)
	}
}

func TestGCD(t *testing.T) {
	var tests = []struct {
		input    Pair
		expected uint64
	}{
		{
			input: Pair{
				a: 1071,
				b: 462,
			},
			expected: 21,
		},
		{
			input: Pair{
				a: 462,
				b: 1071,
			},
			expected: 21,
		},
		{
			input: Pair{
				a: 0,
				b: 1071,
			},
			expected: 1071,
		},
		{
			input: Pair{
				a: 333,
				b: 0,
			},
			expected: 333,
		},
		{
			input: Pair{
				a: 65537,
				b: 97,
			},
			expected: 1,
		},
	}

	for caseid, test := range tests {
		res := GCD(test.input.a, test.input.b)
		if res != test.expected {
			t.Errorf("Case #%d, Actual: %d, Expected: %d", caseid+1, res, test.expected)
		}
		t.Logf("GCD(%d, %d) = %d", test.input.a, test.input.b, res)
	}
}
