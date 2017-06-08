package restTest

import (
	"sort"
	"strings"
	"time"
)

const dateTemplate = "2006-01-02"

// Date is a representation of time.Time using layout "2006-01-02".
// Implements json.Unmarshaler.
type Date struct{ time.Time }

// Implements sort.Interface to enable sorting a date slice.
type byDate []Date

func (d byDate) Len() int           { return len(d) }
func (d byDate) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d byDate) Less(i, j int) bool { return d[i].Before(d[j].Time) }

// Sorts the dates slice in ascending order.
// It makes one call to data.Len to determine n,
// and O(n*log(n)) calls to data.Less and data.Swap.
func (d byDate) Sort() { sort.Sort(d) }

// UnmarshalJSON unmarshalls byte slice into date.
func (date *Date) UnmarshalJSON(b []byte) (err error) {
	// remove quotation marks from string
	s := strings.Trim(string(b), "\"")

	// if the date string is not null, then it's not an error
	if s == "null" {
		return nil
	}
	date.Time, err = time.Parse(dateTemplate, s)
	return
}
