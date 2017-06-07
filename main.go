//+build !test

package main

import (
	"flag"
	"fmt"
)

func init() {
	flag.Int("concurrency", defaultConcurrency, "Number of concurrent go routines that fetch pages")
}

func main() {
	flag.Parse()

	// Get transactions from restTest API server
	ch := FetchAllTransactions()

	// Calculate running daily balances from fetched transactions
	dailyBalances := DailyBalancesFromTransactions(ch)

	// Print running daily balances
	fmt.Printf("Running Daily Balances:\n%s\n-----------\n", dailyBalances)

	// Print overall balance
	fmt.Printf("Total Balance: \t%v\n", dailyBalances.GetRunningBalance())
}
