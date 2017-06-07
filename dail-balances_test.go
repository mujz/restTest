package main

import (
	"testing"
	"time"
)

func TestDailyBalancesString(t *testing.T) {
	d := []Date{
		newDate("2016-04-01"),
		newDate("2016-04-02"),
	}
	db := DailyBalances{
		d,
		map[Date]Amount{
			d[0]: Amount(100.49),
			d[1]: Amount(199.50),
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
		map[Date]Amount{
			d[0]: Amount(100.49),
			d[1]: Amount(100.49),
			d[2]: Amount(199.50),
		},
	}

	db.setRunningDailyBalances()

	tests := []struct {
		expected Amount
		actual   Amount
	}{
		{Amount(100.49), db.balances[d[0]]},
		{Amount(200.98), db.balances[d[1]]},
		{Amount(400.48), db.balances[d[2]]},
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
		map[Date]Amount{
			d[0]: Amount(100.49),
			d[1]: Amount(100.49),
			d[2]: Amount(199.50),
		},
	}

	db.setRunningDailyBalances()

	actual := db.GetRunningBalance()

	if expected := Amount(400.48); expected != actual {
		t.Errorf("Expected running balance: %.2f, Got:%.2f", expected, actual)
	}
}

func TestDailyBalancesSort(t *testing.T) {
	d0 := newDate("2016-04-03")
	d1 := newDate("2016-04-01")
	d2 := newDate("2016-04-02")
	db := DailyBalances{
		[]Date{d0, d1, d2},
		map[Date]Amount{
			d0: Amount(100.49),
			d1: Amount(300.51),
			d2: Amount(200.50),
		},
	}

	db.Sort()

	expected := []struct {
		date   Date   // expected date
		amount Amount // expected amount
	}{
		{d1, Amount(300.51)},
		{d2, Amount(200.50)},
		{d0, Amount(100.49)},
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
					{days[2], "L1", Amount(-1674.45), "C1"},
					{days[0], "L2", Amount(100.01), "C2"},
					{days[1], "L3", Amount(200.02), "C3"},
				}, []Transaction{
					{days[2], "L4", Amount(100.49), "C1"},
					{days[1], "L5", Amount(50.25), "C2"},
					{days[1], "L6", Amount(39.14), "C3"},
				},
			},
			expected: DailyBalances{
				days,
				map[Date]Amount{
					days[0]: Amount(100.01),
					days[1]: Amount(389.42),
					days[2]: Amount(-1184.54),
				},
			},
		},
		// Test case 2: Empty transactions set
		{[][]Transaction{[]Transaction{}}, DailyBalances{[]Date{}, map[Date]Amount{}}},
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
