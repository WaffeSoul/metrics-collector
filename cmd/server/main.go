package main

import (
	"net/http"

	"github.com/WaffeSoul/metrics-collector/internal/server"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func main() {
	storage.StorageGause = storage.Init()
	storage.StorageConter = storage.Init()
	serve := server.InitMux()
	err := http.ListenAndServe(`:8080`, serve)
	if err != nil {
		panic(err)
	}
}
