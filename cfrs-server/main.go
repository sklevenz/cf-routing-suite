package main

import (
	"encoding/json"
	"flag"
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

var (
	count       uint64 = 0
	count_mutex sync.Mutex

	// filled by go build -ldflags="-X main.versionFlag=1.0" or goreleaser
	version string = "snapshot"
	port    string = os.Getenv("PORT")
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

	rootTemplate, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Println(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := struct { // add data to the template
		Version string
	}{version}

	if err := rootTemplate.ExecuteTemplate(w, "layout", data); err != nil {
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
	handleFlags()

	http.HandleFunc("/inc", incHandler)
	http.HandleFunc("/reset", resetHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", rootHandler)

	log.Printf("Server running on http://localhost:%s ...\n", port)
	log.Printf("version: %v", version)
	err := http.ListenAndServe(fmt.Sprintf(":"+port), nil)
	log.Fatal(err)
}

func handleFlags() {
	showVersionPtr := flag.Bool("version", false, "show version info only")
	showHelpPtr := flag.Bool("help", false, "show help")
	showHelp2Ptr := flag.Bool("?", false, "show help")
	portPtr := flag.String("port", "8080", "set server port")

	flag.Parse()

	if *showHelpPtr || *showHelp2Ptr {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *showVersionPtr {
		log.Printf("version: %v", version)
		os.Exit(0)
	}

	if port == "" {
		port = *portPtr
	}
}
