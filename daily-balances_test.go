package restTest

import (
	"testing"
	"time"

	"github.com/mujz/restTest/amount"
)

func TestDailyBalancesString(t *testing.T) {
	d := []Date{
		newDate("2016-04-01"),
		newDate("2016-04-02"),
	}
	db := DailyBalances{
		d,
		map[Date]amount.Amount{
			d[0]: amount.Amount(10049),
			d[1]: amount.Amount(19950),
		},
	}

	actual := db.String()

	if expected := "2016-04-01:\t100.49\n2016-04-02:\t199.50"; expected != actual {
		t.Errorf("Expected daily balances string:\n%s\n\nGot:\n%s", expected, actual)
	}
}

func TestSetRunningDailyBalances(t *testing.T) {
	d := []Date{
		newDate("2016-04-01"),
		newDate("2016-04-02"),
		newDate("2016-04-03"),
	}
	db := DailyBalances{
		d,
		map[Date]amount.Amount{
			d[0]: amount.Amount(10049),
			d[1]: amount.Amount(10049),
			d[2]: amount.Amount(19950),
		},
	}

	db.setRunningDailyBalances()

	tests := []struct {
		expected amount.Amount
		actual   amount.Amount
	}{
		{amount.Amount(10049), db.balances[d[0]]},
		{amount.Amount(20098), db.balances[d[1]]},
		{amount.Amount(40048), db.balances[d[2]]},
	}

	for _, tc := range tests {
		if tc.expected != tc.actual {
			t.Errorf("Expected balance %.2f, Got %.2f", tc.expected, tc.actual)
		}
	}
}

func TestDailyBalancesGetRunningBalance(t *testing.T) {
	d := []Date{
		newDate("2016-04-01"),
		newDate("2016-04-02"),
		newDate("2016-04-03"),
	}
	db := DailyBalances{
		d,
		map[Date]amount.Amount{
			d[0]: amount.Amount(10049),
			d[1]: amount.Amount(10049),
			d[2]: amount.Amount(19950),
		},
	}

	db.setRunningDailyBalances()

	actual := db.GetRunningBalance()

	if expected := amount.Amount(40048); expected != actual {
		t.Errorf("Expected running balance: %.2f, Got:%.2f", expected, actual)
	}
}

func TestDailyBalancesSort(t *testing.T) {
	d0 := newDate("2016-04-03")
	d1 := newDate("2016-04-01")
	d2 := newDate("2016-04-02")
	db := DailyBalances{
		[]Date{d0, d1, d2},
		map[Date]amount.Amount{
			d0: amount.Amount(10049),
			d1: amount.Amount(30051),
			d2: amount.Amount(20050),
		},
	}

	db.Sort()

	expected := []struct {
		date   Date          // expected date
		amount amount.Amount // expected amount
	}{
		{d1, amount.Amount(30051)},
		{d2, amount.Amount(20050)},
		{d0, amount.Amount(10049)},
	}
	actual := db.days

	for i, e := range expected {
		if e.date != actual[i] {
			t.Errorf("Expected date %v, Got %v", e.date, actual[i])
		}

		if a := db.balances[actual[i]]; e.amount != a {
			t.Errorf("Expected amount %.2f, Got %.2f", e.amount, a)
		}
	}
}

func TestDailyBalancesFromTransactions(t *testing.T) {
	days := []Date{
		newDate("2016-04-01"),
		newDate("2016-04-02"),
		newDate("2016-04-03"),
	}

	tests := []struct {
		transactions [][]Transaction
		expected     DailyBalances
	}{
		// Test case 1
		{
			transactions: [][]Transaction{
				[]Transaction{
					{days[2], "L1", amount.Amount(-167445), "C1"},
					{days[0], "L2", amount.Amount(10001), "C2"},
					{days[1], "L3", amount.Amount(20002), "C3"},
				}, []Transaction{
					{days[2], "L4", amount.Amount(10049), "C1"},
					{days[1], "L5", amount.Amount(5025), "C2"},
					{days[1], "L6", amount.Amount(3914), "C3"},
				},
			},
			expected: DailyBalances{
				days,
				map[Date]amount.Amount{
					days[0]: amount.Amount(10001),
					days[1]: amount.Amount(38942),
					days[2]: amount.Amount(-118454),
				},
			},
		},
		// Test case 2: Empty transactions set
		{[][]Transaction{[]Transaction{}}, DailyBalances{[]Date{}, map[Date]amount.Amount{}}},
	}

	for _, tc := range tests {
		// Send transactions over channel
		ch := make(chan []Transaction)
		go func(ch chan []Transaction) {
			for _, t := range tc.transactions {
				ch <- t
			}
			close(ch)
		}(ch)

		actual := DailyBalancesFromTransactions(ch)

		if actual.String() != tc.expected.String() {
			t.Errorf("Expected daily balances:\n%v\n---\nGot:\n%v", tc.expected, actual)
		}

	}
}

func newDate(date string) Date {
	var (
		d   Date
		err error
	)
	d.Time, err = time.Parse(dateTemplate, date)
	if err != nil {
		panic(err)
	}
	return d
}
