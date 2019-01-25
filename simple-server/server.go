package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, is a simple server.")
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
      port = "8080"
    }
    http.HandleFunc("/", handler)
    log.Printf("Server running on http://localhost:%s ...\n", port)
    http.ListenAndServe(fmt.Sprintf(":" + port), nil)
}
