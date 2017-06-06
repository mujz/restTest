package main

import (
	"testing"
)

func TestCustomError(t *testing.T) {
	err := CustomError{"test error", 123}

	if expected := "test error 123"; err.Error() != expected {
		t.Errorf("Expected error string %s, got %s", expected, err.Error())
	}
}
