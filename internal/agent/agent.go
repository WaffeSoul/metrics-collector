package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
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

func (s *Collector) SendToServer(data string) error {
	postURL := s.address + "/update/" + data
	resp, err := http.Post(postURL, "text/plain", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("error: status code %d", resp.StatusCode)
	}
	return nil
}

func (s *Collector) UpdateMetricToServer() {
	for {
		s.mutex.Lock()
		stringsMetric := s.fields.prepareSend()
		counterStr := fmt.Sprintf("counter/PollCount/%d", s.counter)
		randomStr := fmt.Sprintf("gauge/RandomValue/%f", rand.Float64())
		s.mutex.Unlock()
		for _, metric := range stringsMetric {
			err := s.SendToServer(metric)
			if err != nil {
				continue
			}
		}
		err := s.SendToServer(counterStr)
		if err != nil {
			continue
		}
		err = s.SendToServer(randomStr)
		if err != nil {
			continue
		}
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

func (f *fields) prepareSend() (updates []string) {
	updates = append(updates, fmt.Sprintf("gauge/Alloc/%f", f.Alloc))
	updates = append(updates, fmt.Sprintf("gauge/BuckHashSys/%f", f.BuckHashSys))
	updates = append(updates, fmt.Sprintf("gauge/Frees/%f", f.Frees))
	updates = append(updates, fmt.Sprintf("gauge/GCCPUFraction/%f", f.GCCPUFraction))
	updates = append(updates, fmt.Sprintf("gauge/GCSys/%f", f.GCSys))
	updates = append(updates, fmt.Sprintf("gauge/HeapAlloc/%f", f.HeapAlloc))
	updates = append(updates, fmt.Sprintf("gauge/HeapIdle/%f", f.HeapIdle))
	updates = append(updates, fmt.Sprintf("gauge/HeapInuse/%f", f.HeapInuse))
	updates = append(updates, fmt.Sprintf("gauge/HeapObjects/%f", f.HeapObjects))
	updates = append(updates, fmt.Sprintf("gauge/HeapReleased/%f", f.HeapReleased))
	updates = append(updates, fmt.Sprintf("gauge/HeapSys/%f", f.HeapSys))
	updates = append(updates, fmt.Sprintf("gauge/LastGC/%f", f.LastGC))
	updates = append(updates, fmt.Sprintf("gauge/Lookups/%f", f.Lookups))
	updates = append(updates, fmt.Sprintf("gauge/MCacheInuse/%f", f.MCacheInuse))
	updates = append(updates, fmt.Sprintf("gauge/MCacheSys/%f", f.MCacheSys))
	updates = append(updates, fmt.Sprintf("gauge/MSpanInuse/%f", f.MSpanInuse))
	updates = append(updates, fmt.Sprintf("gauge/MSpanSys/%f", f.MSpanSys))
	updates = append(updates, fmt.Sprintf("gauge/Mallocs/%f", f.Mallocs))
	updates = append(updates, fmt.Sprintf("gauge/NextGC/%f", f.NextGC))
	updates = append(updates, fmt.Sprintf("gauge/NumForcedGC/%f", f.NumForcedGC))
	updates = append(updates, fmt.Sprintf("gauge/NumGC/%f", f.NumGC))
	updates = append(updates, fmt.Sprintf("gauge/OtherSys/%f", f.OtherSys))
	updates = append(updates, fmt.Sprintf("gauge/PauseTotalNs/%f", f.PauseTotalNs))
	updates = append(updates, fmt.Sprintf("gauge/StackInuse/%f", f.StackInuse))
	updates = append(updates, fmt.Sprintf("gauge/StackSys/%f", f.StackSys))
	updates = append(updates, fmt.Sprintf("gauge/Sys/%f", f.Sys))
	updates = append(updates, fmt.Sprintf("gauge/TotalAlloc/%f", f.TotalAlloc))
	return
}
