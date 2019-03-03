package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	mode = simulator
}

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
	assert.Equal(t, "text/html", rr.Header().Get("Content-Type"))

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
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
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

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWaitDefaultHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/wait", nil)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(waitHandler)
	handler.ServeHTTP(rr, req)
	duration := time.Now().Sub(start)
	assert.True(t, duration >= waitDefault, "Duration was: ", duration)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWaitHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/wait?wait=100ms", nil)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(waitHandler)
	handler.ServeHTTP(rr, req)
	duration := time.Now().Sub(start)

	assert.True(t, duration >= 100 * time.Millisecond, "Duration was: ", duration)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestWaitHandlerWithError(t *testing.T) {
	req, err := http.NewRequest("GET", "/wait?wait=error", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(waitHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
