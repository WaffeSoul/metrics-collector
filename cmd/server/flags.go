package main

import (
	"flag"
	"os"
	"strconv"
)

var (
	addr            string
	flagLogLevel    string
	storeInterval   int
	fileStoragePath string
	restore         bool
	addrDB          string
)

func parseFlags() {
	flag.StringVar(&addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&addrDB, "d", "", "address and port to connect db")
	flag.StringVar(&flagLogLevel, "l", "info", "log level")
	flag.IntVar(&storeInterval, "i", 300, "interval save store")
	flag.StringVar(&fileStoragePath, "f", "/tmp/metrics-db.json", "path file storage")
	flag.BoolVar(&restore, "r", true, "restore file storage")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		addr = envAddr
	}
	if envaddrDB := os.Getenv("DATABASE_DSN"); envaddrDB != "" {
		addrDB = envaddrDB
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		tempStoreInterval, err := strconv.Atoi(envStoreInterval)
		if err == nil {
			storeInterval = tempStoreInterval
		}
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}

	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		if envRestore == "true" {
			restore = true
		} else if envRestore == "false" {
			restore = false
		}
	}
}
