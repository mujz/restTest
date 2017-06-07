package restTest

import (
	"testing"
	"time"
)

func TestDateUnmarshalJSON(t *testing.T) {
	tests := []struct {
		input      string
		shouldPass bool
	}{
		{"2006-01-02", true},
		{"null", true},
		{"2006-01-00", false},
		{"2006-13-03", false},
		{"1990", false},
		{"1990 March 30", false},
	}

	for _, tc := range tests {
		actual := new(Date)
		err := actual.UnmarshalJSON([]byte(tc.input))
		if tc.shouldPass && err != nil {
			t.Fatal(err)
		} else if !tc.shouldPass {
			if err == nil {
				t.Fatalf("Expected %s test to fail, but it passed instead", tc.input)
			}
			return
		}

		if (actual.Time != time.Time{}) {
			if a := actual.Format(dateTemplate); a != tc.input {
				t.Errorf("Expected date %v, Got %v", tc.input, a)
			}
		}
	}
}

func TestByDateSort(t *testing.T) {
	ts := []struct {
		input    []string
		expected []string
	}{
		{
			[]string{"2006-01-02", "2006-01-04", "2006-01-03"},
			[]string{"2006-01-02", "2006-01-03", "2006-01-04"},
		},
		{
			[]string{"2005-03-01", "2004-04-05", "2005-03-01"},
			[]string{"2004-04-05", "2005-03-01", "2005-03-01"},
		},
	}

	type testCase struct {
		input    []Date
		expected []Date
	}

	tests := make([]testCase, len(ts))

	for i, t := range ts {
		for j, in := range t.input {
			var input Date
			var expected Date
			input.Time, _ = time.Parse(dateTemplate, in)
			expected.Time, _ = time.Parse(dateTemplate, t.expected[j])
			tests[i].input = append(tests[i].input, input)
			tests[i].expected = append(tests[i].expected, expected)
		}
	}

	for _, tc := range tests {
		ByDate(tc.input).Sort()
		for i, a := range tc.input {
			if e := tc.expected[i]; !e.Equal(a.Time) {
				t.Errorf("Expected date %v, Got %v", e, a)
			}
		}
	}

}

func TestByDateLen(t *testing.T) {
	tc := ByDate([]Date{Date{}, Date{}})
	actual := tc.Len()
	if expected := 2; actual != expected {
		t.Errorf("Expected length %d, got %d", expected, actual)
	}
}

func TestByDateSwap(t *testing.T) {
	var (
		d1 Date
		d2 Date
	)
	d1.Time, _ = time.Parse(dateTemplate, "2016-01-01")
	d2.Time, _ = time.Parse(dateTemplate, "2016-01-02")
	tc := ByDate([]Date{d1, d2})
	tc.Swap(0, 1)
	if tc[0] != d2 || tc[1] != d1 {
		t.Errorf("Swap failed. \n1. Expected: %v, Got: %v\n2. Expected: %v, Got: %v", d2, tc[0], d1, tc[1])
	}
}

func TestByDateLess(t *testing.T) {
	var (
		d1 Date
		d2 Date
		d3 Date
	)
	d1.Time, _ = time.Parse(dateTemplate, "2016-01-01")
	d2.Time, _ = time.Parse(dateTemplate, "2016-01-02")
	d3.Time, _ = time.Parse(dateTemplate, "2015-01-02")
	tc := ByDate([]Date{d1, d2, d3})
	if !tc.Less(0, 1) {
		t.Errorf("Expected date %v, to be before %v", tc[0], tc[1])
	}
	if tc.Less(0, 2) {
		t.Errorf("Expected date %v, to be after %v", tc[0], tc[2])
	}
}
