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
	defer errorHandler(w, r)

	lp := filepath.Join("static/html", "layout.html")
	fp := filepath.Join("static/html", filepath.Clean(r.URL.Path))

	if fp == "static/html" {
		fp = "static/html/index.html"
	}

	// Return a 404 if the html doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			err := newHttpError("resource not found", http.StatusNotFound)
			panic(err)
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		err := newHttpError("resource not found", http.StatusNotFound)
		panic(err)
	}

	rootTemplate, err := template.ParseFiles(lp, fp)
	if err != nil {
		err := newHttpError(fmt.Sprintf("Template parsering error - %v", err), http.StatusInternalServerError)
		panic(err)
	}

	data := struct { // add data to the html
		Version string
		Mode    string
	}{version, mode}

	if err := rootTemplate.ExecuteTemplate(w, "layout", data); err != nil {
		err := newHttpError(fmt.Sprintf("Template exeution error - %v", err), http.StatusInternalServerError)
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html")
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	defer errorHandler(w, r)

	if waitValue := r.URL.Query().Get(waitQuery); waitValue != "" {
		wait, err := time.ParseDuration(waitValue)
		if err != nil {
			err := newHttpError(fmt.Sprintf("Error in query parameter %v - %v", waitQuery, err), http.StatusBadRequest)
			panic(err)
		}
		time.Sleep(wait)
	} else {
		time.Sleep(waitDefault)
	}

	tag := r.URL.Query().Get(tagQuery)
	if tag == "" {
		tag = tagDefault
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
	data := mongo.Dial(mode).RecordRequest(&requestData)

	js, err := json.Marshal(data)
	if err != nil {
		err := newHttpError(fmt.Sprintf("Json format error: %v", err), http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	log.Printf("responseData: %v", data)
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	defer errorHandler(w, r)
	data := mongo.Dial(mode).ResetAll()

	js, err := json.Marshal(data)
	if err != nil {
		err := newHttpError(fmt.Sprintf("Json format error: %v", err), http.StatusInternalServerError)
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	log.Printf("responseData: %v", data)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	err := recover()

	if err != nil {
		w.Header().Set("Content-Type", "text/html")

		switch err.(type) {
		case *httpError:
			if httpError, ok := err.(httpError); ok {
				w.WriteHeader(httpError.status)
			}
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		log.Printf("%v", err)
		fmt.Fprintf(w, "%v", err)
	}
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
