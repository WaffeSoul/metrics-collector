package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/WaffeSoul/metrics-collector/internal/handlers"
	"github.com/WaffeSoul/metrics-collector/internal/logger"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func main() {
	parseFlags()
	logger.Initialize()
	db := storage.InitMem(storeInterval, fileStoragePath)
	if true {
		db.LoadStorage()
	}
	// Пока так лучше способа не нашел
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		db.SaveStorage()
		os.Exit(0)
	}()
	if db.InterlvalSave > 0 {
		go db.AutoSaveStorage()
	}

	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Use(handlers.GzipMiddleware)
	r.Route("/", func(r chi.Router) {
		r.Get("/", handlers.GetAll(db))
		r.Post("/update/{type}/{name}/{value}", handlers.PostMetrics(db))
		r.Post("/update/", handlers.PostMetricsJSON(db))
		r.Get("/value/{type}/{name}", handlers.GetValue(db))
		r.Post("/value/", handlers.GetValueJSON(db))

	})
	log.Fatal(http.ListenAndServe(addr, r))

}
