package main

import "fmt"

type HTTPError struct {
	Status     string
	StatusCode int
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("Remote server responded with status: %s", err.Status)
}
