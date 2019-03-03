package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sklevenz/cf-routing-suite/server/mongo"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type responseDataStruct struct {
	Count           uint64 `json:"count"`
	Method          string `json:"method"`
	CfInstanceIndex string `json:"cf-instanceindex"`
	CfInstanceId    string `json:"cf-instanceid"`
	CfApplicationId string `json:"cf-applicationid"`
	data            map[string]interface{}
}

const (
	simulator = "simulator"
	mongodb   = "mongodb"

	envDbMode = "MODE"
	envPort   = "PORT"

	waitDefault = 100 * time.Millisecond
	waitQuery   = "wait"
)

var (
	// filled by go build -ldflags="-X main.version=1.0" or goreleaser
	version = "snapshot"
	port    string
	mode    string
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
		Mode    string
	}{version, mode}

	if err := rootTemplate.ExecuteTemplate(w, "layout", data); err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
	w.Header().Set("Content-Type", "text/html")
}

func waitHandler(w http.ResponseWriter, r *http.Request) {
	query := createQuery()
	query.IncRequestCounter()
	count := query.GetRequestCounter()
	responseData := responseDataStruct{
		count,
		"/wait",
		r.Header.Get("x-cf-instanceindex"),
		r.Header.Get("x-cf-instanceid"),
		r.Header.Get("x-cf-applicationid"),
		nil}

	if wait := r.URL.Query().Get(waitQuery); wait != "" {
		sleep, err := time.ParseDuration(wait)
		if err != nil {
			log.Printf("Error in waitHandler, can't parse wait parameter value: %v, error = %v", wait, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		time.Sleep(sleep)
	} else {
		time.Sleep(waitDefault)
	}

	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error in incHandler: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Printf("responseData: %v", responseData)
}

func incHandler(w http.ResponseWriter, r *http.Request) {
	query := createQuery()
	query.IncRequestCounter()
	count := query.GetRequestCounter()
	responseData := responseDataStruct{
		count,
		"/inc",
		r.Header.Get("x-cf-instanceindex"),
		r.Header.Get("x-cf-instanceid"),
		r.Header.Get("x-cf-applicationid"),
		nil}

	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error in incHandler: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Printf("responseData: %v", responseData)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	query := createQuery()
	query.ResetRequestCounter()
	responseData := responseDataStruct{
		0,
		"/reset",
		r.Header.Get("x-cf-instanceindex"),
		r.Header.Get("x-cf-instanceid"),
		r.Header.Get("x-cf-applicationid"),
		nil}

	err := json.NewEncoder(w).Encode(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error in resetHandler: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("responseData: %v", responseData)
}

func createQuery() mongo.MongoDBQuery {
	if mode == simulator {
		return mongo.CreateSimulator()
	}
	if mode == mongodb {
		return mongo.Create()
	}
	log.Fatalf("Unsupported mode: %v, check environment variable %v", mode, envDbMode)
	os.Exit(-1)
	return nil
}

func main() {
	log.Print("CF-Routing-Suite Server")

	parseFlags()
	readEnvironment()

	log.Printf("Server running on http://localhost:%s ...\n", port)
	log.Printf("version: %v", version)

	http.HandleFunc("/wait", waitHandler)
	http.HandleFunc("/inc", incHandler)
	http.HandleFunc("/reset", resetHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", rootHandler)

	err := http.ListenAndServe(fmt.Sprintf(":"+port), nil)
	log.Fatal(err)
}

func readEnvironment() {
	var ok bool

	if port, ok = os.LookupEnv(envPort); !ok {
		log.Printf("Environment variable %v not set.", envPort)
		os.Exit(1)
	}
	if mode, ok = os.LookupEnv(envDbMode); !ok {
		log.Printf("Environment variable %v not set.", envDbMode)
		os.Exit(1)
	}

	log.Printf("Server mode: %v", mode)
	log.Printf("Server port: %v", port)
}

func parseFlags() {
	showVersionPtr := flag.Bool("v", false, "show version info only")
	showHelpPtr := flag.Bool("help", false, "show help")

	flag.Parse()

	if *showHelpPtr {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *showVersionPtr {
		log.Printf("version: %v", version)
		os.Exit(0)
	}
}
