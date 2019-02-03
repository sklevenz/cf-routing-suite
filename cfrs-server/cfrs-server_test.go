package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "<title>CF Routing Suite</title>"
	if strings.LastIndex(rr.Body.String(), expected) < 0 {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestIncHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/inc", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(incHandler)
	handler.ServeHTTP(rr, req)

	var jsonResponse responseDataStruct
	json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	if jsonResponse.Count != 1 {
		t.Errorf("handler returned wrong count: got %v want %v",
			jsonResponse, 1)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestResetHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/reset", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(resetHandler)
	handler.ServeHTTP(rr, req)

	var jsonResponse responseDataStruct
	json.Unmarshal(rr.Body.Bytes(), &jsonResponse)
	if jsonResponse.Count != 0 {
		t.Errorf("handler returned wrong count: got %v want %v",
			jsonResponse, 0)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
