package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/WaffeSoul/metrics-collector/pkg/constant"
)

var (
	errorNoConnect = errors.New("sendToServer: no connect to server")
)

type Collector struct {
	address        string
	counter        int64
	pollInterval   int64
	reportInterval int64
	fields         fields
	mutex          sync.Mutex
}

type fields struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
}

func NewCollector(address string, pollInterval int64, reportInterval int64) *Collector {
	return &Collector{
		address:        address,
		counter:        0,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		fields:         fields{},
		mutex:          sync.Mutex{},
	}
}

func (s *Collector) SendToServer(data []model.Metrics) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json Marshal %e", err)
	}
	if dataBytes == nil {
		return fmt.Errorf("data is nil")
	}
	postURL := "http://" + s.address + "/updates/"
	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	if _, err = g.Write(dataBytes); err != nil {
		return fmt.Errorf("gzip compress %e", err)
	}
	if err = g.Close(); err != nil {
		return fmt.Errorf("gzip compress %e", err)
	}
	newClient := &http.Client{}
	req, err := http.NewRequest("POST", postURL, &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	var resp *http.Response
	for i := 0; i < 4; i++ {
		resp, err = newClient.Do(req)
		if err != nil {
			if i == 3 {
				err = errors.Join(err, errorNoConnect)
				return err
			}
			fmt.Printf("error: connection refused %e\n", err)
			time.Sleep(time.Duration(constant.RetriTimmer[i]) * time.Second)
		} else {
			defer resp.Body.Close()
			break
		}
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}
	return nil
}

func (s *Collector) UpdateMetricToServer() {
	for {
		var err error
		s.mutex.Lock()
		jsonMetric := s.fields.prepareSend()
		jsonMetric = append(jsonMetric, model.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &s.counter,
		})
		tempRand := rand.Float64()
		jsonMetric = append(jsonMetric, model.Metrics{
			ID:    "RandomValue",
			MType: "gauge",
			Value: &tempRand,
		})
		err = s.SendToServer(jsonMetric)
		if err == nil {
			s.counter = 0
		}
		s.mutex.Unlock()
		time.Sleep(time.Second * time.Duration(s.reportInterval))
	}
}

func (s *Collector) UpdateMetrict() {
	m := &runtime.MemStats{}
	for {
		s.mutex.Lock()
		runtime.ReadMemStats(m)
		s.fields.Alloc = float64(m.Alloc)
		s.fields.BuckHashSys = float64(m.BuckHashSys)
		s.fields.Frees = float64(m.Frees)
		s.fields.GCCPUFraction = float64(m.GCCPUFraction)
		s.fields.GCSys = float64(m.GCSys)
		s.fields.HeapAlloc = float64(m.HeapAlloc)
		s.fields.HeapIdle = float64(m.HeapIdle)
		s.fields.HeapInuse = float64(m.HeapInuse)
		s.fields.HeapObjects = float64(m.HeapObjects)
		s.fields.HeapReleased = float64(m.HeapReleased)
		s.fields.HeapSys = float64(m.HeapSys)
		s.fields.LastGC = float64(m.LastGC)
		s.fields.Lookups = float64(m.Lookups)
		s.fields.MCacheInuse = float64(m.MCacheInuse)
		s.fields.MCacheSys = float64(m.MCacheSys)
		s.fields.MSpanInuse = float64(m.MSpanInuse)
		s.fields.MSpanSys = float64(m.MSpanSys)
		s.fields.Mallocs = float64(m.Mallocs)
		s.fields.NextGC = float64(m.NextGC)
		s.fields.NumForcedGC = float64(m.NumForcedGC)
		s.fields.NumGC = float64(m.NumGC)
		s.fields.OtherSys = float64(m.OtherSys)
		s.fields.PauseTotalNs = float64(m.PauseTotalNs)
		s.fields.StackInuse = float64(m.StackInuse)
		s.fields.StackSys = float64(m.StackSys)
		s.fields.Sys = float64(m.Sys)
		s.fields.TotalAlloc = float64(m.TotalAlloc)
		s.counter += 1
		s.mutex.Unlock()
		time.Sleep(time.Duration(s.pollInterval) * time.Second)
	}
}

func (f fields) prepareSend() (updates []model.Metrics) {
	values := reflect.ValueOf(f)
	typesOf := values.Type()
	for i := 0; i < values.NumField(); i++ {
		temp := values.Field(i).Float()
		updates = append(updates, model.Metrics{
			ID:    typesOf.Field(i).Name,
			MType: "gauge",
			Value: &temp,
		})
	}
	return
}
