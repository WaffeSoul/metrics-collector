package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/WaffeSoul/metrics-collector/internal/handlers"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func main() {
	parseFlags()
	fmt.Println(addr)
	db := storage.InitMem()
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(db))
		r.Post("/update/{type}/{name}/{value}", handlers.PostMetrics(db))
		r.Get("/value/{type}/{name}", handlers.GetValue(db))
	})

	log.Fatal(http.ListenAndServe(addr, r))

}
