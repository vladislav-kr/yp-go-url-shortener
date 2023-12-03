package main

import "flag"

func parseFlags(host, redirectHost *string) {

	if len(*host) == 0 {
		flag.StringVar(host, "a", ":8080", "address and port to run server")
	}

	if len(*redirectHost) == 0 {
		flag.StringVar(redirectHost, "b", "http://localhost:8080", "redirect address")
	}

	flag.Parse()
}
