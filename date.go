package main

import (
	"strings"
	"time"
)

type Date struct {
	time.Time
}

const dateTemplate = "2006-01-02"

func (date *Date) UnmarshalJSON(b []byte) (err error) {
	// remove quotation marks from string
	s := strings.Trim(string(b), "\"")

	// if the date string is not null, then it's not an error
	if s == "null" {
		date.Time = time.Time{}
		return nil
	}
	date.Time, err = time.Parse(dateTemplate, s)
	return
}
