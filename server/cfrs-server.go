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

const (
	envDbMode = "MODE"
	envPort   = "PORT"

	waitDefault = 0 * time.Millisecond
	waitQuery   = "wait"

	tagQuery   = "tag"
	tagDefault = "default"
)

var (
	// filled by go build -ldflags="-X main.version=1.0" or goreleaser
	version = "snapshot"
	port    string
	mode    string
)

func main() {
	log.Print("CF-Routing-Suite Server")

	parseFlags()
	readEnvironment()

	log.Printf("Server running on http://localhost:%s ...\n", port)
	log.Printf("version: %v", version)

	http.HandleFunc("/api/probe", probeHandler)
	http.HandleFunc("/api/reset", resetHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", rootHandler)

	err := http.ListenAndServe(fmt.Sprintf(":"+port), nil)
	log.Fatal(err)
}

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
		log.Println(err)
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

func probeHandler(w http.ResponseWriter, r *http.Request) {
	if waitValue := r.URL.Query().Get(waitQuery); waitValue != "" {
		wait, err := time.ParseDuration(waitValue)
		if err != nil {
			errorHandler(w, r, fmt.Sprintf("Time format error: %v", err), http.StatusInternalServerError)
			return
		}
		time.Sleep(wait)
	} else {
		time.Sleep(waitDefault)
	}

	tag := r.URL.Query().Get(tagQuery)
	if tag == "" {
		tag = tagDefault
	}

	query, err := mongo.Dial(mode)
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("MongoDB query error: %v", err), http.StatusInternalServerError)
		return
	}

	requestData := mongo.RequestData{
		Method:          r.Method,
		Remote:          r.RemoteAddr,
		Timestamp:       time.Now(),
		Url:             r.URL.String(),
		XB3ParentSpanId: r.Header.Get("x-b3-parentspanid"),
		XB3SpanId:       r.Header.Get("x-b3-spanid"),
		XB3TraceId:      r.Header.Get("x-b3-traceid"),
		XForwardedFor:   r.Header.Get("x-forwarded-for"),
		Tag:             tag,
	}

	data := query.RecordRequest(&requestData)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("Json format error: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("responseData: %v", data)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	query, err := mongo.Dial(mode)
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("MongoDB query error: %v", err), http.StatusInternalServerError)
		return
	}

	data := query.ResetAll()

	js, err := json.Marshal(data)
	if err != nil {
		errorHandler(w, r, fmt.Sprintf("Json format error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	log.Printf("responseData: %v", data)
}

func errorHandler(w http.ResponseWriter, r *http.Request, msg string, status int) {
	log.Printf("Internal Server Error 500: %v", msg)
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Internal Server Error 500 - %v", msg)

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
