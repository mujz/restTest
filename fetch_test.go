package restTest

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/mujz/restTest/money"
)

var (
	mockPageStr = `{
	"totalCount": %d,
	"page": %d,
	"transactions": [{
			"Date": "2013-12-13",
			"Ledger": "Insurance Expense",
			"Amount": "-117.81",
			"Company": "LONDON DRUGS 78 POSTAL VANCOUVER BC"
	}, {
			"Date": "2013-12-13",
			"Ledger": "Equipment Expense",
			"Amount": "-520.85",
			"Company": "ECHOSIGN xxxxxxxx6744 CA xx8.80 USD @ xx0878"
	}, {
			"Date": "2013-12-13",
			"Ledger": "Equipment Expense",
			"Amount": "-5518.17",
			"Company": "APPLE STORE #R280 VANCOUVER BC"
	}, {
			"Date": "2013-12-12",
			"Ledger": "Postage & Shipping Expense",
			"Amount": "-30.69",
			"Company": "DHL YVR GW RICHMOND BC"
	}, {
			"Date": "2013-12-12",
			"Ledger": "Office Expense",
			"Amount": "-42.53",
			"Company": "FEDEX xxxxx5291 MISSISSAUGA ON"
	}, {
			"Date": "2013-12-20",
			"Ledger": "Equipment Expense",
			"Amount": "-1874.75",
			"Company": "NINJA STAR WORLD VANCOUVER BC"
	}, {
			"Date": "2013-12-12",
			"Ledger": "Postage & Shipping Expense",
			"Amount": "-30.69",
			"Company": "DHL YVR GW RICHMOND BC"
	}, {
			"Date": "2013-12-12",
			"Ledger": "Office Expense",
			"Amount": "-42.53",
			"Company": "FEDEX xxxxx5291 MISSISSAUGA ON"
	}, {
			"Date": "2013-12-12",
			"Ledger": "Web Hosting & Services Expense",
			"Amount": "-63.01",
			"Company": "GROWINGCITY.COM xxxxxx4926 BC"
	}, {
			"Date": "2013-12-12",
			"Ledger": "Business Meals & Entertainment Expense",
			"Amount": "-91.12",
			"Company": "NESTERS MARKET #x0064 VANCOUVER BC"
	}]
}`
	emptyPageJSON = []byte(`{
	"totalCount": 0,
	"page": 1,
	"transactions": []
	}`)
	mockPage Page
	_        = json.Unmarshal([]byte(fmt.Sprintf(mockPageStr, 10, 1)), &mockPage)
)

type restTestHandler struct {
	// If status is not 200, it writes it back with an empty body.
	status int
	// Number of transactions to serve.
	totalCount int
	// If payload is not provided, the handler uses the default mock page.
	payload []byte // optional
}

// Mock server
func (h *restTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.status != http.StatusOK {
		w.WriteHeader(h.status)
	} else {
		pageNumber, err := strconv.Atoi(strings.Trim(r.URL.Path, "/"))
		if err != nil {
			panic(err)
		}

		// return 404 if attempting to get too many pages and totalCount is non-negative
		if h.totalCount > 0 {
			pageCount := int(math.Floor(float64(h.totalCount-1)/transactionsPerPage) + 1)
			if pageNumber > pageCount {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		// Serve handler's payload or default mock page
		if h.payload != nil {
			w.Write(h.payload)
			return
		}
		w.Write([]byte(fmt.Sprintf(mockPageStr, h.totalCount, pageNumber)))
	}
}

// Test without a mock server
func TestFetchPageRemote(t *testing.T) {
	p, err := FetchPage(1)
	if err != nil {
		t.Fatal(err)
	}
	if n := p.Page; n != 1 {
		t.Errorf("Expected page number to be %d, got %d", 1, n)
	}
}

func TestFetchPage(t *testing.T) {
	// start a mock server
	handler := restTestHandler{}
	mockServer := httptest.NewServer(&handler)
	defer mockServer.Close()

	tests := []struct {
		status     int
		totalCount int
		payload    []byte
		url        string
		pageNumber int
		shouldPass bool
	}{
		{http.StatusOK, 10, nil, pageURL(1, mockServer.URL+"/%d"), 1, true},

		{http.StatusNotFound, 0, nil, pageURL(1, mockServer.URL+"/%d"), 1, false},
		{http.StatusOK, 0, nil, "invalid url", 0, false},
		{http.StatusOK, 0, []byte("Not JSON"), pageURL(1, mockServer.URL+"/%d"), 0, false},
	}

	for _, tc := range tests {
		handler.status = tc.status
		handler.totalCount = tc.totalCount
		handler.payload = tc.payload

		// Test success case
		p, err := fetchPage(tc.url)

		if tc.shouldPass {
			if err != nil {
				t.Fatal(err)
			}

			// Assertions
			if n := p.Page; n != tc.pageNumber {
				t.Errorf("Expected page number %d, got %d", tc.pageNumber, n)
			}
			if c := p.TotalCount; c != tc.totalCount {
				t.Errorf("Expected total count %d, got %d", tc.totalCount, c)
			}

			n := len(p.Transactions)
			if tc.totalCount > transactionsPerPage && n != transactionsPerPage {
				t.Errorf("Expected number of transactions %d, got %d", tc.totalCount, n)
			} else if tc.totalCount <= transactionsPerPage && n != tc.totalCount {
				t.Errorf("Expected number of transactions %d, got %d", transactionsPerPage, n)
			}

			// Assert first transaction
			tr := p.Transactions[0]
			if expected := newDate("2013-12-13"); !tr.Date.Equal(expected.Time) {
				t.Errorf("Expected transaction date %v, got %v", expected, tr.Date)
			}
			if expected := money.Amount(-11781); tr.Amount != expected {
				t.Errorf("Expected transaction amount %s, got %s", expected, tr.Amount)
			}
			if tr.Ledger == "" {
				t.Errorf("Transaction ledger is empty.")
			}
			if tr.Company == "" {
				t.Errorf("Transaction Company is empty.")
			}
		} else {
			if e, ok := err.(HTTPError); ok && e.StatusCode != tc.status {
				t.Error(err)
			}
		}
	}
}

// Test FetchAllPages without a mock server (i.e. against the real server)
func TestFetchAllPagesFromRemoteServer(t *testing.T) {
	// first we need to know how many many transactions to expect
	p, err := FetchPage(1)
	if err != nil {
		t.Fatal(err)
	}
	expectedCount := p.TotalCount

	// now get all transactions from remote server
	ch := FetchAllTransactions()

	var all []Transaction
	for {
		actual, more := <-ch
		if !more {
			break
		}
		all = append(all, actual...)

		// assert we didn't get more than the expected transactions per page count
		if a := len(actual); a > transactionsPerPage {
			t.Errorf("Expected transactions per page less than or equal to %d, Got %d", transactionsPerPage, a)
		}
	}

	// assert we got the expected total count
	if actual := len(all); actual != expectedCount {
		t.Errorf("Expected count %d, Got %d", expectedCount, actual)
	}
}

// Test fetchAllPages with a mock server
func TestFetchAllPages(t *testing.T) {
	type testCase struct {
		status     int
		totalCount int
		payload    []byte
		shouldPass bool
	}
	tests := []testCase{
		// success cases
		{http.StatusOK, 10000, nil, true},
		{http.StatusOK, 0, emptyPageJSON, true},

		// error cases
		{http.StatusOK, 10, []byte(fmt.Sprintf(mockPageStr, 20, 1)), false},
		{http.StatusOK, -1, []byte(`Not JSON`), false},
		{http.StatusNotFound, -1, nil, false},
		{http.StatusInternalServerError, -1, nil, false},
	}

	var wg sync.WaitGroup
	for _, tc := range tests {
		// Start the mock server
		handler := restTestHandler{tc.status, tc.totalCount, tc.payload}
		mockServer := httptest.NewServer(&handler)

		ch := make(chan []Transaction)

		wg.Add(1)

		// if it's expected fail, then recover and make sure that it panicked
		if !tc.shouldPass {
			go func(ch chan []Transaction, url string) {
				defer mockServer.Close()
				defer wg.Done()
				defer assertPanic(t)

				go func(ch chan []Transaction) { <-ch }(ch)

				fetchAllTransactions(ch, url, DefaultConcurrency)
			}(ch, mockServer.URL+"/%d")
		} else {
			go fetchAllTransactions(ch, mockServer.URL+"/%d", DefaultConcurrency)

			go func(ch chan []Transaction, tc testCase, mockServer *httptest.Server) {
				defer mockServer.Close()
				defer wg.Done()

				// stores all fetched transaction so we can check their length later
				var all []Transaction

				for {
					actual, more := <-ch
					all = append(all, actual...)
					if !more {
						break
					}

					for i, a := range actual {
						if expected := mockPage.Transactions[i]; a.String() != expected.String() {
							t.Errorf("Expected transaction %v\nGot %v", expected, a)
						}
					}

				}

				if actual := len(all); actual != tc.totalCount {
					t.Errorf("Expected total count %d, Got %d", tc.totalCount, actual)
				}
			}(ch, tc, mockServer)
		}
	}
	wg.Wait()
}

func TestTransportString(t *testing.T) {
	tr := Transaction{
		Date:    newDate("2006-02-01"),
		Ledger:  "Ledger 1",
		Amount:  money.Amount(10049),
		Company: "Bench",
	}
	expected := fmt.Sprintf("{\n\tDate: %v,\n\tLedger: %s,\n\tAmount: %v,\n\tCompany: %s\n}",
		tr.Date.Format(dateTemplate), tr.Ledger, tr.Amount, tr.Company)
	actual := tr.String()

	if expected != actual {
		t.Errorf("Expected transaction string %s, Got %s", expected, actual)
	}
}

func TestPageString(t *testing.T) {
	p := Page{
		TotalCount: 1,
		Page:       1,
		Transactions: []Transaction{{
			Date:    newDate("2006-02-01"),
			Ledger:  "Ledger 1",
			Amount:  money.Amount(10049),
			Company: "Bench",
		}},
	}
	expected := fmt.Sprintf("{\n\tTotal Count: %d,\n\tPage: %d,\n\tTransactions: %v\n}",
		p.TotalCount, p.Page, p.Transactions)
	actual := p.String()

	if expected != actual {
		t.Errorf("Expected page string %s, Got %s", expected, actual)
	}
}

func assertPanic(t *testing.T) {
	r := recover()
	if r == nil {
		t.Error("Expected it to panic, but it didn't")
	}
}
