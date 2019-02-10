package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "<title>CF Routing Suite</title>")

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

	assert.Equal(t, uint64(1), jsonResponse.Count)
	assert.Equal(t, http.StatusOK, rr.Code)
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
	assert.Equal(t, uint64(0), jsonResponse.Count)

	assert.Equal(t, http.StatusOK, rr.Code)
}
