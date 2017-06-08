package money

import "testing"

func TestAmountUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input      string
		expected   float64
		shouldPass bool
	}{
		{"100.40", 10040, true},
		{"0.40", 40, true},
		{".40", 40, true},
		{".4", 40, true},
		{"100", 10000, true},
		{"00", 0, true},
		{"0", 0, true},
		{"0.0", 0, true},
		{"0100.040", 10004, true},
		{"100.04", 10004, true},
		{"100.400", 10040, true},
		{"100.499", 10050, true},
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
			t.Errorf("Expected amount %d, Got %d", expected, *actual)
		}
	}
}

func TestFromFloat(t *testing.T) {
	tests := []struct {
		in       float64
		expected Amount
	}{
		{10.45, Amount(1045)},
		{10.045, Amount(1005)},
		{10.0, Amount(1000)},
	}

	for _, tc := range tests {
		if actual := FromFloat(tc.in); actual != tc.expected {
			t.Errorf("Expected amount %d, got %d", tc.expected, actual)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		in       Amount
		expected string
	}{
		{Amount(100), "1.00"},
		{Amount(10000), "100.00"},
		{Amount(5), "0.05"},
		{Amount(-500), "-5.00"},
		{Amount(-5), "-0.05"},
	}

	for _, tc := range tests {
		if actual := tc.in.String(); actual != tc.expected {
			t.Errorf("Expected amount string %s, Got %s", tc.expected, actual)
		}
	}
}
