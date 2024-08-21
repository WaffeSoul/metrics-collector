package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/WaffeSoul/metrics-collector/internal/handlers"
	"github.com/WaffeSoul/metrics-collector/internal/logger"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func main() {
	var db *storage.Database
	var err error
	parseFlags()
	logger.Initialize()
	if addrDB == "" {
		db, err = storage.New("mem", storeInterval, fileStoragePath, "")
	} else {
		db, err = storage.New("postgresql", 0, "", addrDB)
	}
	if err != nil {
		log.Fatal(err)
	}
	go db.DB.AutoSaveStorage()
	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(handlers.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(db))
		r.Post("/update/{type}/{name}/{value}", handlers.PostMetrics(db))
		r.Post("/update/", handlers.PostMetricsJSON(db))
		r.Get("/ping", handlers.PingDB(db))
		r.Get("/value/{type}/{name}", handlers.GetValue(db))
		r.Post("/value/", handlers.GetValueJSON(db))

	})
	log.Fatal(http.ListenAndServe(addr, r))

}
