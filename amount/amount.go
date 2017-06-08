/*
This package deals with representing monetary amounts in Go. It represents money amounts as integers. For example, 30 dollars and 5 cents would be represented as 3005. However, calling the String() method on an amount returns a string in dollars and cents. For example, the previous amount would be "30.05".
*/
package amount

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// A representation of money amounts in cents.
// Ex. 10.50 is 1050 cents. Implements json.Unmarshaler.
type Amount int

// Unmarshals byte slice into amount.
func (amount *Amount) UnmarshalJSON(b []byte) error {
	// remove quotation marks from string.
	s := strings.Trim(string(b), "\"")

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*amount = FromFloat(f)
	return nil
}

func (a Amount) String() string {
	dollars := a / 100
	cents := int(math.Abs(float64(a % 100)))
	layout := "%d.%.2d"

	// if dollars are 0 but amount is negative,
	// we need to add the sign before the 0.
	if dollars == 0 && a < 0 {
		layout = "-" + layout
	}

	return fmt.Sprintf(layout, dollars, cents)
}

// Returns Amount from float64 rounded up to the nearest 100th.
func FromFloat(f float64) Amount {
	return Amount(round(f * 100))
}

// Round float to int.
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Adds amount b to a and returns the result.
// It rounds the sum to the nearest 100th.
func (a Amount) Add(b Amount) Amount {
	return a + b
}
