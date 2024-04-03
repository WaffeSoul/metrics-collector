package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/WaffeSoul/metrics-collector/internal/storage"
)

func PostMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	data := strings.Split(r.URL.Path, "/")
	if len(data) != 5 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	typeM := data[2]
	nameM := data[3]
	valueStrM := data[4]
	switch typeM {
	case "gauge":
		valueM, err := strconv.ParseFloat(valueStrM, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.StorageGause.Add(nameM, valueM)

	case "counter":
		valueM, err := strconv.ParseInt(valueStrM, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		storage.StorageConter.Add(nameM, valueM)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func InitMux() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/update/", PostMetrics)
	return
}
