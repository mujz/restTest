package main

import (
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
	*amount = Amount(n)
	return nil
}
