package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
	Addr           string `env:"ADDRESS"`
}

var cfg Config

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
	if cfg.Addr == "" {
		cfg.Addr = addr
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = reportInterval
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = pollInterval
	}
}
