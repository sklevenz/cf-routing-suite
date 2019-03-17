package main

import (
	"fmt"
	"net/http"
)

type httpError struct {
	status  int
	message string
}

func newHttpError(message string, status int) *httpError {
	return &httpError{
		message: message,
		status:  status,
	}
}

func (e *httpError) Error() string {
	return fmt.Sprintf("%v %v - %v", e.status, http.StatusText(e.status), e.message)
}
