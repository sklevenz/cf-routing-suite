package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type responseDataStruct struct {
	Count           uint64 `json:"count"`
	Method          string `json:"method"`
	CfInstanceIndex string `json:"cf-instanceindex"`
	CfInstanceId    string `json:"cf-instanceid"`
	CfApplicationId string `json:"cf-applicationid"`
}

type versionStruct struct {
	Version string
	Commit  string
	Date    string
}

var (
	count       uint64 = 0
	count_mutex sync.Mutex

	// filled by go build -ldflags="-X main.versionFlag=1.0 ..." or goreleaser
	versionFlag string = "snapshot"
	commitFlag  string = "n/a"
	dateFlag    string = "n/a"

	version versionStruct
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("template", "layout.html")
	fp := filepath.Join("template", filepath.Clean(r.URL.Path))

	if fp == "template" {
		fp = "template/index.html"
	}

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", version); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func incHandler(w http.ResponseWriter, r *http.Request) {
	incCounter()
	responseData := responseDataStruct{
		count,
		"/inc",
		r.Header.Get("x-cf-instanceindex"),
		r.Header.Get("x-cf-instanceid"),
		r.Header.Get("x-cf-applicationid")}

	json.NewEncoder(w).Encode(responseData)

	w.Header().Set("Content-Type", "application/json")
	log.Printf("responseData: %v", responseData)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	resetCounter()
	responseData := responseDataStruct{
		count,
		"/reset",
		r.Header.Get("x-cf-instanceindex"),
		r.Header.Get("x-cf-instanceid"),
		r.Header.Get("x-cf-applicationid")}

	json.NewEncoder(w).Encode(responseData)
	w.Header().Set("Content-Type", "application/json")
	log.Printf("responseData: %v", responseData)
}

func incCounter() {
	count_mutex.Lock()
	count++
	count_mutex.Unlock()
}

func resetCounter() {
	count_mutex.Lock()
	count = 0
	count_mutex.Unlock()
}

func main() {
	version = versionStruct{versionFlag, commitFlag, dateFlag}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/inc", incHandler)
	http.HandleFunc("/reset", resetHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", rootHandler)

	log.Printf("Server running on http://localhost:%s ...\n", port)
	log.Printf("version: %v", version)
	err := http.ListenAndServe(fmt.Sprintf(":"+port), nil)
	log.Fatal(err)
}
