package main

import (
	"flag"
)

var (
	addr string
)

func parseFlags() {
	flag.StringVar(&addr, "a", ":8080", "address and port to run server")
	flag.Parse()
}
