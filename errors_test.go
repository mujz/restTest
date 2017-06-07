package main

import "testing"

func TestHTTPError(t *testing.T) {
	err := HTTPError{"404 Not Found", 404}
	expected := "Remote server responded with status: 404 Not Found"
	if actual := err.Error(); actual != expected {
		t.Errorf("Expected error %s, Got %s", expected, actual)
	}
}
