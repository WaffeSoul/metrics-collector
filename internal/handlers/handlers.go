package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
	"github.com/go-chi/chi/v5"
)

func PostMetrics(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Content-Type") {
		case "application/json":
			var resJSON model.Metrics
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&resJSON)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Header().Add("Content-Type", "text/plain")
			switch resJSON.MType {
			case "gauge":
				db.StorageGauge.Add(resJSON.ID, resJSON.Value)

			case "counter":
				valueOldM, ok := db.StorageCounter.Get(resJSON.ID)
				if ok {
					*resJSON.Delta += valueOldM.(int64)
				}
				db.StorageCounter.Add(resJSON.ID, *resJSON.Delta)
			default:
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func PostMetricsOLD(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		typeM := chi.URLParam(r, "type")
		if typeM == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		nameM := chi.URLParam(r, "name")
		if nameM == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		valueStrM := chi.URLParam(r, "value")
		if nameM == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch typeM {
		case "gauge":
			valueM, err := strconv.ParseFloat(valueStrM, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			db.StorageGauge.Add(nameM, valueM)

		case "counter":
			valueM, err := strconv.ParseInt(valueStrM, 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			valueOldM, ok := db.StorageCounter.Get(nameM)
			if ok {
				valueM += valueOldM.(int64)
			}
			db.StorageCounter.Add(nameM, valueM)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetValueOLD(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		typeM := chi.URLParam(r, "type")
		nameM := chi.URLParam(r, "name")
		switch typeM {
		case "gauge":
			valueM, err := db.StorageGauge.Get(nameM)
			if !err {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%v", valueM)))
		case "counter":
			valueM, err := db.StorageCounter.Get(nameM)
			if !err {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("%v", valueM)))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}

func GetValue(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var resJSON model.Metrics
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&resJSON)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch resJSON.MType {
		case "gauge":
			valueM, ok := db.StorageGauge.Get(resJSON.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			*resJSON.Value = valueM.(float64)
		case "counter":
			valueM, ok := db.StorageCounter.Get(resJSON.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			*resJSON.Delta = valueM.(int64)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jsonResp, err := json.Marshal(resJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

func GetAll(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		data := db.StorageCounter.GetAll()
		for name, value := range data {
			w.Write([]byte(fmt.Sprintf("%v: %v\n", name, value.Value)))
		}
		data = db.StorageGauge.GetAll()
		for name, value := range data {
			w.Write([]byte(fmt.Sprintf("%v: %v\n", name, value.Value)))
		}
	}
}
