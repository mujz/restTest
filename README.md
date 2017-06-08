[![codecov](https://codecov.io/gh/mujz/restTest/branch/master/graph/badge.svg)](https://codecov.io/gh/mujz/restTest)
[![Build Status](https://travis-ci.org/mujz/restTest.svg?branch=master)](https://travis-ci.org/mujz/restTest)
[![Go Report Card](https://goreportcard.com/badge/github.com/mujz/Resttest)](https://goreportcard.com/report/github.com/mujz/Resttest)

# restTest

This is a backend coding assignment for Bench's full-stack software developer application. Scroll down to the [getting started section](#getting-started) for the project description or visit [resttest.bench.co](http://resttest.bench.co/).

## Run it yourself

### From compiled binaries

Download and run the right binary for your system from the [`/bin`](bin/) directory.

### Build and install binaries from source

You need to have golang installed, then run:

```bash
go get -v github.com/mujz/restTest
cd $GOPATH/src/github.com/mujz/restTest
make install
```

## Usage

Once you install the binary, just run:

```bash
restTest
```

You can also set the flag `-concurrency`, which is the number of go routines that run conncurently to fetch transaction pages.

## Implementation

Since we need to execute multiple operations concurrently (ex. fetching pages, calculating the balance), it's preferrable to use a language that has support for coroutines (or lightweight threads) such as Go or Kotlin. Thus, I'm choosing to use Go. Here's how this works:

1. We use the standard `net/http` package to make a simple GET request. This function will be used by the next 2 operations.
1. We fetch the first page and read the `totalCount` to know how many pages to fetch. We then start as many coroutines as the number of pages minus 1 (since we don't need to re-fetch the first page). However, we must limit the number of concurrent coroutines to avoid running into too many http requests errors. Each coroutine fetches a page and sends its transactions over a channel to another coroutine that processes them.
1. The coroutine that receives the transactions calculates each date's transactions and dies. We start a new coroutine for each page of transactions until the channel closes. Once the channel closes, we sort the daily balances and add them to each other to find the running daily balances.
1. After all coroutines finish, we print the running daily balances and overall balance to the console.

## Known Limitations

### Monetary Amounts Data Structure

There are multiple ways to represent fractioned monetary amounts. One way is to store dollars and cents as separate integers (or only count using cents). Another is to use decimals (which Go lacks). My choice was to store them as cents. If the number has more than 2 decimal places, I round it to the nearest cent.

### Too many loops

With the current implementation, the app makes these loops:

1. Fetch pages and increment daily balances.
1. Sort daily balances.
1. Calculate running daily balances.
1. Print daily balances.

Since the number of transactions is not that big, this implementation is clean and works well. However, if the number of transaction grows, some of those loops should be consolidated.

The last two loops could merge into one. We can prepare the daily balances string as we calculate the running daily balances and have it ready when we call print.

Consolidating the first 2 loops requires using an insertion sort algorithm, which is less efficient than the one I'm using (which is implemented by the Go standard package).

Merging the second and third don't work either since we need the daily transactions to be sorted by day before we can calculate the daily balances.

Therefore, only the last two loops can be joined, but since they make the code less clear and the cost is not that big (O(2n) instead of O(n)), I chose to go with this implementation.

Additionally, I've optimized the sorting by saving the days (the daily balances map has `day: balance`) into a separate slice since sorting a slice is much faster than sorting a map.

---

## Getting Started

Welcome to the Bench Rest Test. The purpose of this exercise is to demonstrate your ability to reason about rudimentary APIs and data transformation. You can use any language you feel comfortable with.

We would like you to write an app that we can run from the command line, which:

1. Connects to a REST API documented below and fetches all pages of financial transactions.
1. Calculates total balance and prints it to the console, where balance is the sum of all amounts in all transactions. For example, if I have 3 transactions each for $4, then the total balance would be $12.
1. Calculates running daily balances and prints them to the console. For example, if I have 3 transactions for the 5th 6th 7th, each for $5, then the running daily balance on the 4th would be $0, on the 5th would be $5, on the 6th would be $10, on the 7th it would be $15.

Include unit tests for this application.

To submit your work, create a git repository on Github and send send us the link in an email

## API Documentation

This API provides access to all the transactions in an imaginary bank account.

This is a REST API, providing JSON-formatted data over HTTP.

There is a limit to how many transactions that can be returned in a single request, so the transactions are split into "pages". You will have to download all the pages to get all the data.

**GET** http://resttest.bench.co/transactions/{page}.json

Responses

**200 OK**

```js
{
  "totalCount": 32, // Integer, total number of transactions across all pages
  "page": 1, // Integer, current page
  "transactions": [
    {
      "Date": "2013-12-22", // String, date of transaction
      "Ledger": "Phone & Internet Expense", // String, ledger name
      "Amount": "-110.71", // String, amount
      "Company": "SHAW CABLESYSTEMS CALGARY AB" // String, company name
    },
    ...
  ]
}
```

**404 NOT FOUND**

No response body

## Considerations

1. Approach this problem as you would in the real-world. Consider errors that may occur when fetching and transforming data from the API, such as non-200 http responses.
1. Consider scalability when picking data abstractions and algorithms; what would happen if the transaction list was considerably larger?
1. Coding style matters. Ensure your code is consistent and easy to follow. Leave comments where appropriate and use meaningful methods and variables.
1. Avoid overly complex code. The complexity of the solution should make sense for the problems you're solving.
1. Document limitations and trade-offs of your code if appropriate.
1. Include a README explaining how to install and/or run your software.

