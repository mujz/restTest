package restTest

import (
	"math"
	"strconv"
	"strings"
)

// A representation of money amounts as dollars and 2-digit number cents.
// Ex. 10.50 is 10 dollars and 50 cents.
// Implements json.Unmarshaler
type Amount float64

// Unmarshals byte slice into amount
func (amount *Amount) UnmarshalJSON(b []byte) error {
	// remove quotation marks from string
	s := strings.Trim(string(b), "\"")

	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*amount = Amount(toFixed(n, 2))
	return nil
}

// Adds amount b to a and returns the result.
// It rounds the sum to the nearest 100th.
func (a Amount) Add(b Amount) Amount {
	return Amount(toFixed(float64(a+b), 2))
}

// Round float to int
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Round number up to the passed precision
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
