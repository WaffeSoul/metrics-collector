package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/WaffeSoul/metrics-collector/internal/storage"
	"github.com/go-chi/chi/v5"
)

func PostMetrics(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeM := chi.URLParam(r, "type")
		nameM := chi.URLParam(r, "name")
		valueStrM := chi.URLParam(r, "value")
		switch typeM {
		case "gauge":
			valueM, err := strconv.ParseFloat(valueStrM, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			db.StorageGause.Add(nameM, valueM)

		case "counter":
			valueM, err := strconv.ParseInt(valueStrM, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			db.StorageConter.Add(nameM, valueM)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}
func GetValue(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		typeM := chi.URLParam(r, "type")
		nameM := chi.URLParam(r, "name")
		switch typeM {
		case "gauge":
			valueM, err := db.StorageGause.Get(nameM)
			if !err {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%v", valueM)))
		case "counter":
			valueM, err := db.StorageConter.Get(nameM)
			if !err {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%v", valueM)))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}

func GetAll(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		data := db.StorageConter.GetAll()
		for name, value := range data {
			w.Write([]byte(fmt.Sprintf("%v: %v\n", name, value.Value)))
		}
		data = db.StorageGause.GetAll()
		for name, value := range data {
			w.Write([]byte(fmt.Sprintf("%v: %v\n", name, value.Value)))
		}
	}
}
