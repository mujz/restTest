package restTest

import (
	"fmt"
	"strings"
	"sync"
)

// Representation of financial transaction.
type Transaction struct {
	// Date of the transaction in layout 2006-01-02.
	Date   Date
	Ledger string
	// Transactoun amount with a precision of up to 2 decimal places.
	Amount Amount
	// Company name
	Company string
}

// Data structure for representing dates and their balances.
// It is optimized for efficient and fast sorting by date.
type DailyBalances struct {
	// Store days, which are the keys of the dailyBalances map, in a separate slice.
	// This makes it more efficient for sorting; sorting a slice is much faster than
	// sorting a map.
	days     []Date
	balances map[Date]Amount
}

// Returns the transaction's fields formatted as JSON
func (t Transaction) String() string {
	return fmt.Sprintf("{\n\tDate: %v,\n\tLedger: %s,\n\tAmount: %v,\n\tCompany: %s\n}", t.Date.Format(dateTemplate), t.Ledger, t.Amount, t.Company)
}

// Returns the daily balances fields formatted as day:	balance
func (db DailyBalances) String() string {
	var s []string
	for _, day := range db.days {
		s = append(s, fmt.Sprintf("%s:\t%.2f", day.Format(dateTemplate), db.balances[day]))
	}
	return strings.Join(s, "\n")
}

// Adds each day's balance to the next starting from second day.
// Completes in O(n) number of iterations.
func (db *DailyBalances) setRunningDailyBalances() {
	for i := 1; i < len(db.days); i++ {
		db.balances[db.days[i]] = db.balances[db.days[i]].Add(db.balances[db.days[i-1]])
	}
}

// Returns the last day's balance.
func (db DailyBalances) GetRunningBalance() Amount {
	return db.balances[db.days[len(db.days)-1]]
}

// Sorts the daily balances by date in ascending order.
// It implements the standard package's sort.Sort.
func (db *DailyBalances) Sort() {
	byDate(db.days).Sort()
}

// Receives transaction slices over the channel, sorts them, and calculates
// their running daily balances. It returns after the channel is closed.
//
// Blocks until it finishes processing all transactions.
func DailyBalancesFromTransactions(ch chan []Transaction) DailyBalances {
	var (
		wg    sync.WaitGroup
		mutex = &sync.Mutex{}

		db = DailyBalances{balances: make(map[Date]Amount)}
	)

	// Waits for transaction slices to come then launches a go routine for each
	// to loop over each transaction and add it to the daily balance.
	for {
		ts, more := <-ch
		if !more {
			break
		}

		wg.Add(1)
		go func(ts []Transaction) {
			defer wg.Done()
			for _, t := range ts {
				// Must lock to read map and increment amount
				mutex.Lock()

				// if day doesn't already exist, add it to the days slice
				if _, ok := db.balances[t.Date]; !ok {
					db.days = append(db.days, t.Date)
				}

				// increment daily balance
				db.balances[t.Date] = db.balances[t.Date].Add(t.Amount)

				mutex.Unlock()
			}
		}(ts)
	}

	// Wait until all transactions have been processed
	wg.Wait()

	// Sort days slice
	db.Sort()

	// Calculate running daily balances
	db.setRunningDailyBalances()

	return db
}
