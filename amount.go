package main

import (
	"math"
	"strconv"
	"strings"
)

type Amount float64

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

func (a Amount) Add(b Amount) Amount {
	return Amount(toFixed(float64(a+b), 2))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
