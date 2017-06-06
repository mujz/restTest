package main

import "testing"

func TestAmountUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input      string
		expected   float64
		shouldPass bool
	}{
		{"100.40", 100.40, true},
		{"0.40", 0.40, true},
		{".40", 0.40, true},
		{".4", 0.4, true},
		{"100", 100, true},
		{"00", 0, true},
		{"0", 0, true},
		{"0.0", 0, true},
		{"0100.040", 100.04, true},
		{"100.04", 100.04, true},
		{".4 ", 0, false},
		{"100.04f", 0, false},
	}

	for _, tc := range tests {
		actual := new(Amount)
		err := actual.UnmarshalJSON([]byte(tc.input))
		if tc.shouldPass && err != nil {
			t.Fatal(err)
		} else if !tc.shouldPass && err == nil {
			t.Fatalf("Expected %s test to fail, but it passed instead", tc.input)
		}

		if expected := Amount(tc.expected); expected != *actual {
			t.Errorf("Expected amount %f, Got %f", expected, *actual)
		}
	}
}
