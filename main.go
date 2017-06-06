package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	urlTemplate = "http://resttest.bench.co/transactions/%d.json"
)

type Transaction struct {
	Date    Date
	Ledger  string
	Amount  Amount
	Company string
}

type Page struct {
	TotalCount   int
	Page         int
	Transactions []Transaction
}

func FetchPage(pageNumber int) (*Page, error) {
	return fetchPage(fmt.Sprintf(urlTemplate, pageNumber))
}

func fetchPage(url string) (*Page, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		RemoteServerErr.Status = res.StatusCode
		return nil, RemoteServerErr
	}

	page := new(Page)
	err = json.NewDecoder(res.Body).Decode(page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

func main() {
}
