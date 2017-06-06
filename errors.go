package main

import "fmt"

type CustomError struct {
	ErrorString string
	Status      int
}

var RemoteServerErr = CustomError{ErrorString: "Remote server responded with status"}

func (err CustomError) Error() string {
	return fmt.Sprintf("%s %d", err.ErrorString, err.Status)
}
