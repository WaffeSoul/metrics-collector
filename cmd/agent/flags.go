package main

import (
	"flag"
)

var (
	reportInterval int64
	pollInterval   int64
	addr           string
)

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&reportInterval, "r", 10, "report interval in seconds")
	flag.Int64Var(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
}
