package restTest

import (
	"sort"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

type ByDate []Date

func (d ByDate) Len() int           { return len(d) }
func (d ByDate) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByDate) Less(i, j int) bool { return d[i].Before(d[j].Time) }

// Sorts the dates slice in ascending order.
// It makes one call to data.Len to determine n,
// and O(n*log(n)) calls to data.Less and data.Swap.
func (d ByDate) Sort() { sort.Sort(d) }

const dateTemplate = "2006-01-02"

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
