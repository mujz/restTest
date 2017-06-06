package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	mockPage = []byte(`{
	"totalCount": 38,
	"page": 4,
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
}`)
)

type restTestHandler struct {
	status  int
	payload []byte
}

func (h *restTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.status != http.StatusOK {
		w.WriteHeader(h.status)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(h.payload)
	}
}

func TestFetchPage(t *testing.T) {
	// Test without a mock server
	page, err := FetchPage(1)
	if err != nil {
		t.Fatal(err)
	}
	if n := page.Page; n != 1 {
		t.Errorf("Expected page number to be %d, got %d", 1, n)
	}

	// start a mock server
	handler := restTestHandler{http.StatusOK, mockPage}
	mockServer := httptest.NewServer(&handler)
	defer mockServer.Close()

	// Test success case
	page, err = fetchPage(mockServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Assertions
	if n := page.Page; n != 4 {
		t.Errorf("Expected page number %d, got %d", 4, n)
	}
	if c := page.TotalCount; c != 38 {
		t.Errorf("Expected total count %d, got %d", 38, c)
	}
	if n := len(page.Transactions); n != 8 {
		t.Errorf("Expected number of transactions %d, got %d", 8, n)
	}

	tr := page.Transactions[0]
	if expected, _ := time.Parse(dateTemplate, "2013-12-13"); !tr.Date.Equal(expected) {
		t.Errorf("Expected transaction date %v, got %v", expected, tr.Date)
	}
	if expected := Amount(-117.81); tr.Amount != expected {
		t.Errorf("Expected transaction amount %s, got %s", expected, tr.Amount)
	}
	if tr.Ledger == "" {
		t.Errorf("Transaction ledger is empty.")
	}
	if tr.Company == "" {
		t.Errorf("Transaction Company is empty.")
	}
	// ---

	// Test 404 case
	handler.status = http.StatusNotFound
	page, err = fetchPage(mockServer.URL)
	if e, ok := err.(CustomError); !ok || e.Status != http.StatusNotFound {
		t.Error(err)
	}

	// Test 500 case
	handler.status = http.StatusInternalServerError
	page, err = fetchPage(mockServer.URL)
	if e, ok := err.(CustomError); !ok || e.Status != http.StatusInternalServerError {
		t.Error(err)
	}

	// Test invalid URL case
	if _, err = fetchPage("invalid url"); err == nil {
		t.Error("Expected fetch page to fail because passed url is invalid")
	}

	// Test invalid json body
	handler.payload = []byte("Not JSON")
	handler.status = http.StatusOK
	if _, err := fetchPage(mockServer.URL); err == nil {
		t.Error("Expected fetch page to fail because response payload is not of the expected format")
	}

}
