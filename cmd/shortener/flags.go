package main

import (
	"flag"
)

const (
	defHost         = ":8080"
	defRedirectHost = "http://localhost:8080"
)

func parseFlags(host, redirectHost *string) {
	var fHost, fRedirectHost string

	flag.StringVar(&fHost, "a", defHost, "address and port to run server")
	flag.StringVar(&fRedirectHost, "b", defRedirectHost, "redirect address")

	flag.Parse()

	if len(*host) == 0 {
		*host = fHost
	}

	if len(*redirectHost) == 0 {
		*redirectHost = fRedirectHost
	}
}
