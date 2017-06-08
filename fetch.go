package restTest

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"sync"
)

const (
	// REST API url.
	urlTemplate = "http://resttest.bench.co/transactions/%d.json"
	// Maximum number of transaction per page.
	transactionsPerPage = 10
	// Maximum number of idle http connections.
	maxIdleConnections = 100
	// DefaultConcurrency is the default number of concurrent go routines to fetch pages.
	DefaultConcurrency = 20
)

var (
	// Concurrency is the number of concurrent go routines that fetch pages.
	Concurrency = DefaultConcurrency
)

// Page represents a slice of transactions.
type Page struct {
	// Total number of transactions (in this page plus all other pages).
	TotalCount int
	// Page number. Index starts 1 (not 0).
	Page int
	// Page's transactions. Must not exceed 10 entries.
	Transactions []Transaction
}

func init() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = maxIdleConnections
}

// Returns the page's fields formatted as JSON.
func (p Page) String() string {
	return fmt.Sprintf("{\n\tTotal Count: %d,\n\tPage: %d,\n\tTransactions: %v\n}",
		p.TotalCount, p.Page, p.Transactions)
}

// FetchPage fetches the page from the restTest API server and decodes it into Page.
// Returns HTTPError if response status is not 200.
func FetchPage(pageNumber int) (*Page, error) {
	return fetchPage(pageURL(pageNumber, urlTemplate))
}

// Returns page url from base url template and page number
// urlTemplate must specify where the pageNumber goes with %d.
func pageURL(n int, urlTemplate string) string {
	return fmt.Sprintf(urlTemplate, n)
}

// Calls HTTP GET to the passed url and decodes the response body into Page struct.
// returns HTTPError if response status is not 200
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

// FetchAllTransactions fetches all pages from the restTest API and
// puts the slice of transactions (max transactions per slice = 10)
// from each page over a channel. It closes the channel once all
// transactions are put to the channel.
//
// Panics if encounters an error.
func FetchAllTransactions() chan []Transaction {
	ch := make(chan []Transaction)
	go fetchAllTransactions(ch, urlTemplate, Concurrency)
	return ch
}

// Fetches the first page to get total number of pages to fetch.
// Then launches a go routine to fetch each page. After the last
// transaction is put in the channel, it closes the channel.
//
// It only launches as many go routines as the passed concurrency flag
func fetchAllTransactions(ch chan []Transaction, urlTemplate string, concurrency int) {
	// Fetch the first page
	p, err := fetchPage(pageURL(1, urlTemplate))
	if err != nil {
		panic(err)
	}

	// Put the first page's transactions in the channel
	ch <- p.Transactions

	// Calculate the number of remaining pages to fetch
	pageCount := int(
		math.Floor(
			float64(
				(p.TotalCount-1)/transactionsPerPage,
			),
		) + 1,
	)

	// Close the channel if there are no more pages
	if pageCount < 2 {
		close(ch)
		return
	}

	// Initialize WaitGroup
	var wg sync.WaitGroup
	wg.Add(pageCount - 1)

	// If an error occurs in a child go routine,
	// send it over the channel and panic from the parent
	done := make(chan error, 1)
	// Semaphore to limit the number of go routines
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
