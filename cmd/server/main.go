package main

import (
	"net/http"

	"github.com/WaffeSoul/metrics-collector/internal/handlers"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func main() {
	storage.StorageGause = storage.Init()
	storage.StorageConter = storage.Init()
	serve := handlers.InitMux()
	err := http.ListenAndServe(`localhost:8080`, serve)
	if err != nil {
		panic(err)
	}
}
