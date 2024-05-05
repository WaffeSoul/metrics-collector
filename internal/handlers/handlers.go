package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/WaffeSoul/metrics-collector/internal/storage"
	"github.com/go-chi/chi/v5"
)

func PostMetricsJSON(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("New")
		switch r.Header.Get("Content-Type") {
		case "application/json":
			w.Header().Add("Content-Type", "text/plain; text/html")
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
		w.Header().Add("Content-Type", "text/plain; text/html")
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

func GetValue(db *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; text/html")
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
		w.Header().Add("Content-Type", "text/plain; text/html")
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

// type GzipWriter struct {
// 	http.ResponseWriter
// 	Writer io.Writer
// }

// func (w GzipWriter) Write(b []byte) (int, error) {
// 	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
// 	return w.Writer.Write(b)
// }

// func MiddlewareGzip(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		// создаём gzip.Writer поверх текущего w
// 		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
// 		if err != nil {
// 			io.WriteString(w, err.Error())
// 			return
// 		}
// 		defer gz.Close()

// 		w.Header().Set("Content-Encoding", "gzip")
// 		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
// 		next.ServeHTTP(GzipWriter{ResponseWriter: w, Writer: gz}, r)
// 	})
// }

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

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

// compressReader реализует интерфейс io.ReadCloser и позволяет прозрачно для сервера
// декомпрессировать получаемые от клиента данные
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
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := newCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(ow, r)
	})
}
