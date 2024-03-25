package main

import (
	"flag"

)

const (
	defHost         = ":8080"
	defRedirectHost = "http://localhost:8080"
	defFilePath     = "\\tmp\\short-url-db.json"
)

func parseFlags(host, redirectHost, filePath, pgDNS *string) {
	var fHost, fRedirectHost, fFilePath string

	flag.StringVar(&fHost, "a", defHost, "address and port to run server")
	flag.StringVar(&fRedirectHost, "b", defRedirectHost, "redirect address")
	flag.StringVar(&fFilePath, "f", defFilePath, "redirect address")
	flag.StringVar(pgDNS, "d", "", "database connection address")
	flag.Parse()

	if len(*host) == 0 {
		*host = fHost
	}

	if len(*redirectHost) == 0 {
		*redirectHost = fRedirectHost
	}

	if len(*filePath) == 0 {
		*filePath = fFilePath
	}
}
