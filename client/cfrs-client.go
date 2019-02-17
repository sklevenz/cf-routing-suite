package main

import (
	"flag"
	"log"
	"os"
)

var (
	// filled by go build -ldflags="-X main.versionFlag=1.0" or goreleaser
	version string = "snapshot"
)

func main() {
	log.Print("CF-Routing-Suite Client")
	handleFlags()

	log.Printf("client running")
	log.Printf("version: %v", version)

}

func handleFlags() {
	showVersionPtr := flag.Bool("v", false, "show version info only")
	showHelpPtr := flag.Bool("help", false, "show help")
	showHelp2Ptr := flag.Bool("?", false, "show help")

	flag.Parse()

	if *showHelpPtr || *showHelp2Ptr {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *showVersionPtr {
		log.Printf("version: %v", version)
		os.Exit(0)
	}
}
