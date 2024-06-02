package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
	"github.com/go-chi/chi/v5"
)

func PostMetricsJSON(db *storage.MemStorage) http.HandlerFunc {
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
			fmt.Println(resJSON)
			switch resJSON.MType {
			case "gauge":
				if resJSON.Value == nil {
					fmt.Println("error value")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				db.StorageGauge.Add(resJSON.ID, *resJSON.Value)

			case "counter":
				if resJSON.Delta == nil {
					fmt.Println("error value")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				valueOldM, ok := db.StorageCounter.Get(resJSON.ID)
				if ok {
					*resJSON.Delta += valueOldM.(int64)
				}
				db.StorageCounter.Add(resJSON.ID, *resJSON.Delta)
			default:
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if db.LastSave.Add(time.Duration(db.InterlvalSave)*time.Second).Before(time.Now()) || db.InterlvalSave == 0 {
				db.SaveStorage()
			}
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	}
}

func PostMetrics(db *storage.MemStorage) http.HandlerFunc {
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
		if db.LastSave.Add(time.Duration(db.InterlvalSave)*time.Second).Before(time.Now()) || db.InterlvalSave == 0 {
			db.SaveStorage()
		}
		w.WriteHeader(http.StatusOK)

	}
}

func GetValue(db *storage.MemStorage) http.HandlerFunc {
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

func GetValueJSON(db *storage.MemStorage) http.HandlerFunc {
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
		switch resJSON.MType {
		case "gauge":
			valueM, ok := db.StorageGauge.Get(resJSON.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			fmt.Println(valueM.(float64))
			temp := valueM.(float64)
			resJSON.Value = &temp
		case "counter":
			valueM, ok := db.StorageCounter.Get(resJSON.ID)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			temp := valueM.(int64)
			resJSON.Delta = &temp
		default:
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

func GetAll(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptData := r.Header.Get("Accept")
		if acceptData == "html/text" {
			w.Header().Add("Content-Type", "text/html")
		} else {
			w.Header().Add("Content-Type", "text/plain")
		}
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

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
