package main

import (
	"flag"
	"fmt"

	"github.com/mujz/restTest"
)

var (
	concurrency = flag.Int("concurrency", restTest.DefaultConcurrency, "Number of concurrent go routines that fetch pages")
)

func main() {
	flag.Parse()
	restTest.Concurrency = *concurrency

	// Get transactions from restTest API server
	ch := restTest.FetchAllTransactions()

	// Calculate running daily balances from fetched transactions
	dailyBalances := restTest.DailyBalancesFromTransactions(ch)

	// Print running daily balances
	fmt.Printf("Running Daily Balances:\n%s\n-----------\n", dailyBalances)

	// Print overall balance
	fmt.Printf("Total Balance: \t%v\n", dailyBalances.GetRunningBalance())
}
