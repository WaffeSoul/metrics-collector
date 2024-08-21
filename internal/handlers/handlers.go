package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
	"github.com/go-chi/chi/v5"
)

func PostMetricsJSON(db *storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New")
		switch r.Header.Get("Content-Type") {
		case "application/json":
			w.Header().Add("Content-Type", "text/plain")
			var resJSON model.Metrics
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&resJSON)
			if err != nil {
				fmt.Println(r.Body)
				fmt.Println("error json")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if len(resJSON.ID) == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			err = db.DB.AddJSON(resJSON)
			if err != nil {
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

func PostMetrics(db *storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Old")
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
		err := db.DB.Add(typeM, nameM, valueStrM)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	}
}

func GetValue(db *storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		typeM := chi.URLParam(r, "type")
		nameM := chi.URLParam(r, "name")
		valueM, err := db.DB.Get(typeM, nameM)
		if err != nil && err.Error() == "NotFound" {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", valueM)))
	}
}

func GetValueJSON(db *storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		w.Header().Add("Content-Type", "application/json")
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
		if len(resJSON.ID) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		resJSON, err = db.DB.GetJSON(resJSON)
		if err != nil && err.Error() == "NotFound" {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jsonResp, err := json.Marshal(resJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

func GetAll(db *storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptData := r.Header.Get("Accept")
		if acceptData == "html/text" {
			w.Header().Add("Content-Type", "text/html")
		} else {
			w.Header().Add("Content-Type", "text/plain")
		}
		w.WriteHeader(http.StatusOK)
		data := db.DB.GetAll()
		w.Write(data)
	}
}

func PingDB(db *storage.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.DB.Ping()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

}
