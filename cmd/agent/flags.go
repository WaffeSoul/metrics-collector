package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type config struct {
	reportInterval int64  `env:"REPORT_INTERVAL"`
	pollInterval   int64  `env:"POLL_INTERVAL"`
	addr           string `env:"ADDRESS"`
}

var cfg config

func parseFlags() {
	var addr string
	var reportInterval int64
	var pollInterval int64
	flag.StringVar(&addr, "a", "localhost:8080", "address and port to run server")
	flag.Int64Var(&reportInterval, "r", 10, "report interval in seconds")
	flag.Int64Var(&pollInterval, "p", 2, "poll interval in seconds")
	flag.Parse()
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.addr == "" {
		cfg.addr = addr
	}
	if cfg.reportInterval == 0 {
		cfg.reportInterval = reportInterval
	}
	if cfg.pollInterval == 0 {
		cfg.pollInterval = pollInterval
	}
}
