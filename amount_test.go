package restTest

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
		{"100.400", 100.40, true},
		{"100.499", 100.50, true},
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

func TestAdd(t *testing.T) {
	tests := []struct {
		a        Amount
		b        Amount
		expected Amount
	}{
		{Amount(10.10), Amount(10.10), Amount(20.20)},
		{Amount(10.11), Amount(10.99), Amount(21.10)},
		{Amount(-10.10), Amount(10.10), Amount(0)},
		{Amount(10.10), Amount(-10.10), Amount(0)},
		{Amount(-10.10), Amount(-10.10), Amount(-20.20)},
	}

	for _, tc := range tests {
		actual := tc.a.Add(tc.b)

		if tc.expected != actual {
			t.Errorf("Expected sum %v, Got %v", tc.expected, actual)
		}
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		input    float64
		expected int
	}{
		{10.49, 10},
		{10.49999, 10},
		{10.50, 11},
		{10.9999, 11},
		{10.0, 10},
	}

	for _, tc := range tests {
		actual := round(tc.input)

		if tc.expected != actual {
			t.Errorf("Expected rounded number %d, Got %v", tc.expected, actual)
		}
	}
}

func TestToFixed(t *testing.T) {
	tests := []struct {
		input     float64
		precision int
		expected  float64
	}{
		{10.491, 2, 10.49},
		{10.491, 1, 10.5},
		{10.49999, 3, 10.5},
		{10.4999, 4, 10.4999},
		{10.4, 2, 10.4},
		{10.5, 0, 11},
	}

	for _, tc := range tests {
		actual := toFixed(tc.input, tc.precision)

		if tc.expected != actual {
			t.Errorf("Expected %d fixed precision number %v, Got %v", tc.precision, tc.expected, actual)
		}
	}
}
