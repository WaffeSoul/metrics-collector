package main

import (
	"flag"
	"os"
)

var (
	addr string
)

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		addr = envAddr
	}
}
