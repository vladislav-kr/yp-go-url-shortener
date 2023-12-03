package main

import (
	"flag"
)

const (
	defHost         = ":8080"
	defRedirectHost = "http://localhost:8080"
)

func parseFlags(host, redirectHost *string) {

	if len(*host) == 0 {
		flag.StringVar(host, "a", defHost, "address and port to run server")
	}

	if len(*redirectHost) == 0 {
		flag.StringVar(redirectHost, "b", defRedirectHost, "redirect address")
	}

	flag.Parse()

	if len(*host) == 0 {
		*host = defHost
	}

	if len(*redirectHost) == 0 {
		*redirectHost = defRedirectHost
	}
}
