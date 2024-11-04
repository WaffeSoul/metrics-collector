package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	gomem "github.com/shirou/gopsutil/v4/mem"

	"github.com/WaffeSoul/metrics-collector/internal/crypto"
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
	keyHash        string
	Rate           int64
}

type Fields struct {
	Alloc           float64
	BuckHashSys     float64
	Frees           float64
	GCCPUFraction   float64
	GCSys           float64
	HeapAlloc       float64
	HeapIdle        float64
	HeapInuse       float64
	HeapObjects     float64
	HeapReleased    float64
	HeapSys         float64
	LastGC          float64
	Lookups         float64
	MCacheInuse     float64
	MCacheSys       float64
	MSpanInuse      float64
	MSpanSys        float64
	Mallocs         float64
	NextGC          float64
	NumForcedGC     float64
	NumGC           float64
	OtherSys        float64
	PauseTotalNs    float64
	StackInuse      float64
	StackSys        float64
	Sys             float64
	TotalAlloc      float64
	TotalMemory     float64
	FreeMemory      float64
	CPUutilization1 float64
}

func NewCollector(address string, pollInterval int64, reportInterval int64, key string, rate int64) *Collector {
	return &Collector{
		address:        address,
		counter:        0,
		keyHash:        key,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		Rate:           rate,
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
	hash := crypto.HashWithKey(dataBytes, s.keyHash)
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
	req.Header.Set("HashSHA256", hash)
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

func (s *Collector) UpdateMetricToServer(inCh chan Fields) {
	for {
		var err error
		fields := <-inCh
		jsonMetric := fields.prepareSend()
		counter := atomic.LoadInt64(&s.counter)
		atomic.AddInt64(&s.counter, -counter)
		jsonMetric = append(jsonMetric, model.Metrics{
			ID:    "PollCount",
			MType: "counter",
			Delta: &counter,
		})
		tempRand := rand.Float64()
		jsonMetric = append(jsonMetric, model.Metrics{
			ID:    "RandomValue",
			MType: "gauge",
			Value: &tempRand,
		})
		err = s.SendToServer(jsonMetric)
		if err != nil {
			atomic.AddInt64(&s.counter, counter)
			s.counter = 0
		}
		time.Sleep(time.Second * time.Duration(s.reportInterval))
	}
}

func (s *Collector) UpdateMetrict(outCh chan Fields) {
	m := &runtime.MemStats{}
	var fields Fields
	for {
		runtime.ReadMemStats(m)
		fields.Alloc = float64(m.Alloc)
		fields.BuckHashSys = float64(m.BuckHashSys)
		fields.Frees = float64(m.Frees)
		fields.GCCPUFraction = float64(m.GCCPUFraction)
		fields.GCSys = float64(m.GCSys)
		fields.HeapAlloc = float64(m.HeapAlloc)
		fields.HeapIdle = float64(m.HeapIdle)
		fields.HeapInuse = float64(m.HeapInuse)
		fields.HeapObjects = float64(m.HeapObjects)
		fields.HeapReleased = float64(m.HeapReleased)
		fields.HeapSys = float64(m.HeapSys)
		fields.LastGC = float64(m.LastGC)
		fields.Lookups = float64(m.Lookups)
		fields.MCacheInuse = float64(m.MCacheInuse)
		fields.MCacheSys = float64(m.MCacheSys)
		fields.MSpanInuse = float64(m.MSpanInuse)
		fields.MSpanSys = float64(m.MSpanSys)
		fields.Mallocs = float64(m.Mallocs)
		fields.NextGC = float64(m.NextGC)
		fields.NumForcedGC = float64(m.NumForcedGC)
		fields.NumGC = float64(m.NumGC)
		fields.OtherSys = float64(m.OtherSys)
		fields.PauseTotalNs = float64(m.PauseTotalNs)
		fields.StackInuse = float64(m.StackInuse)
		fields.StackSys = float64(m.StackSys)
		fields.Sys = float64(m.Sys)
		fields.TotalAlloc = float64(m.TotalAlloc)
		atomic.AddInt64(&s.counter, 1)
		outCh <- fields
		time.Sleep(time.Duration(s.pollInterval) * time.Second)
	}
}

func (f Fields) prepareSend() (updates []model.Metrics) {
	values := reflect.ValueOf(f)
	typesOf := values.Type()
	for i := 0; i < values.NumField(); i++ {
		temp := values.Field(i).Float()
		if temp == 0 {
			continue
		}
		updates = append(updates, model.Metrics{
			ID:    typesOf.Field(i).Name,
			MType: "gauge",
			Value: &temp,
		})
	}
	return
}

func (s *Collector) UpdataGopsutil(outCh chan Fields) {
	var fields Fields
	for {
		time.Sleep(time.Duration(s.pollInterval) * time.Second)
		v, err := gomem.VirtualMemory()
		if err != nil {
			log.Fatal(err)
			continue
		}
		fields.TotalMemory = float64(v.Total)
		fields.FreeMemory = float64(v.Free)
		temp, err := cpu.Percent(0, true)
		if err != nil {
			log.Fatal(err)
			continue
		}
		fields.CPUutilization1 = temp[0]
		atomic.AddInt64(&s.counter, 1)
		outCh <- fields
		time.Sleep(time.Duration(s.pollInterval) * time.Second)
	}
}
