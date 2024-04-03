package main

import (
	"net/http"

	"github.com/WaffeSoul/metrics-collector/internal/app"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func main() {
	storage.StorageGause = storage.Init()
	storage.StorageConter = storage.Init()
	serve := app.InitMux()
	err := http.ListenAndServe(`:8080`, serve)
	if err != nil {
		panic(err)
	}
}
