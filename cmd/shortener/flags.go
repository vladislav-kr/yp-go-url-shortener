package main

import "flag"

func parseFlags(host, redirectHost *string) {

	flag.StringVar(host, "a", ":8080", "address and port to run server")
	flag.StringVar(redirectHost, "b", "http://localhost:8080", "redirect address")
	
	flag.Parse()
}
