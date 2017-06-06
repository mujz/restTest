package main

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
