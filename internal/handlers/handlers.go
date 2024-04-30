package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func PostMetrics(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var resJson model.Metrics
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&resJson)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		switch resJson.MType {
		case "gauge":
			db.StorageGauge.Add(resJson.ID, resJson.Value)

		case "counter":
			valueOldM, ok := db.StorageCounter.Get(resJson.ID)
			if ok {
				*resJson.Delta += valueOldM.(int64)
			}
			db.StorageCounter.Add(resJson.ID, *resJson.Delta)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetValue(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var resJson model.Metrics
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&resJson)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch resJson.MType {
		case "gauge":
			valueM, ok := db.StorageGauge.Get(resJson.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			*resJson.Value = valueM.(float64)
		case "counter":
			valueM, ok := db.StorageCounter.Get(resJson.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			*resJson.Delta = valueM.(int64)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jsonResp, err := json.Marshal(resJson)
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
