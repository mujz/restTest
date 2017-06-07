package restTest

import "fmt"

// Returned when a remote server responds with a non-200 status code.
type HTTPError struct {
	// Error status string. Ex. 404 Not Found. Matches http.Response.Status.
	Status string
	// Error status code. Ex. 404. Matches http.Response.StatusCode.
	StatusCode int
}

// Implements error
func (err HTTPError) Error() string {
	return fmt.Sprintf("Remote server responded with status: %s", err.Status)
}
