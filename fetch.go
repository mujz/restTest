package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
)

const (
	// REST API url
	urlTemplate = "http://resttest.bench.co/transactions/%d.json"
	// Maximum number of transaction per page
	TransactionsPerPage = 10
	// default number of concurrent go routines to fetch pages
	MaxIdleConnections = 100
)

var (
	Concurrency = 20 // number of concurrent go routines that fetch pages
)

type Page struct {
	TotalCount   int
	Page         int
	Transactions []Transaction
}

func init() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = MaxIdleConnections
}

func (p Page) String() string {
	return fmt.Sprintf("{\n\tTotal Count: %d,\n\tPage: %d,\n\tTransactions: %v\n}",
		p.TotalCount, p.Page, p.Transactions)
}

func FetchPage(pageNumber int) (*Page, error) {
	return fetchPage(pageURL(pageNumber, urlTemplate))
}

func pageURL(n int, urlTemplate string) string {
	return fmt.Sprintf(urlTemplate, n)
}

func fetchPage(url string) (*Page, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, HTTPError{res.Status, res.StatusCode}
	}

	page := new(Page)
	err = json.NewDecoder(res.Body).Decode(page)
	if err != nil {
		return nil, err
	}

	return page, nil
}

// Fetches a all pages from the restTest API and
// returns the slice of transactions (max transactions per slice = 10)
// from each page over a channel. It closes the channel once it all
// transactions are put to the channel.
//
// Panics if encountes an error
func FetchAllTransactions() chan []Transaction {
	ch := make(chan []Transaction)
	go fetchAllTransactions(ch, urlTemplate, Concurrency)
	return ch
}

// Fetches the first page to get total number of pages to fetch.
// Then launches a go rountine to fetch each page. After the last
// channel is put in the channel, it closes the channel.
//
// It only launces as many go routines as the concurrency value
func fetchAllTransactions(ch chan []Transaction, urlTemplate string, concurrency int) {
	// Fetch the first page
	p, err := fetchPage(pageURL(1, urlTemplate))
	if err != nil {
		panic(err)
	}

	// Put the first page's transactions in the channel
	ch <- p.Transactions

	// Calc the number of remaining pages to fetch
	pageCount := int(
		math.Floor(
			float64(
				(p.TotalCount-1)/TransactionsPerPage,
			),
		) + 1,
	)

	// Close the channel if there are no more pages
	if pageCount < 2 {
		close(ch)
		return
	}

	// Wait for all go routines to finish, then close the channel
	var wg sync.WaitGroup
	wg.Add(pageCount - 1)

	// If an error occurs in a child go routine,
	// send it over the channel and panic from the parent
	done := make(chan error, 1)
	// Semaphore lock to limit the number of go routines
	sem := make(chan bool, concurrency)

	for i := 2; i <= pageCount; i++ {
		sem <- true // increment semaphore
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()

			// Fetch page
			p, err := fetchPage(pageURL(i, urlTemplate))
			if err != nil {
				done <- err
				return
			}

			// Put page's transactions in channel
			ch <- p.Transactions
		}(i)
	}

	// Wait for all go routines to finish then close the done channel
	go func() {
		wg.Wait()
		close(done)
	}()

	// Waits until an error or channel close.
	// Panics if error. Otherwise closes the transactions channel
	err = <-done
	if err != nil {
		panic(err)
	}
	close(ch)
}
