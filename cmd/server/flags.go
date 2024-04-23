package main

import (
	"flag"
	"os"
)

var (
	addr string
	flagLogLevel string
)

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		addr = envAddr
	}
}
