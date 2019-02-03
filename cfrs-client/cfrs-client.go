package main

import (
	"flag"
	"log"
	"os"
)

var (
	// filled by go build -ldflags="-X main.versionFlag=1.0" or goreleaser
	version string = "snapshot"
	showVersion = *flag.Bool("version", false, "show version info only")
)


func main() {
	if showVersion {
		log.Printf("version: %v", version)
		os.Exit(0)
	}

	log.Printf("cfrs-client running")
	log.Printf("version: %v", version)

}
