//+build !test

package main

import (
	"fmt"
)

func main() {
	// Get transactions from restTest API server
	ch := FetchAllTransactions()

	// Calculate running daily balances from fetched transactions
	dailyBalances := DailyBalancesFromTransactions(ch)

	// Print running daily balances
	fmt.Println(dailyBalances)

	// Print overall balance
	fmt.Printf("Overall Balance: %v\n", dailyBalances.GetRunningBalance())
}
