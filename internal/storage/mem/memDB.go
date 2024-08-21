package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/WaffeSoul/metrics-collector/internal/model"
)

type Repository struct {
	StorageGauge   *MemStorage
	StorageCounter *MemStorage
	InterlvalSave  int
	PathFile       string
}

func InitMem(interlval int, path string) *Repository {
	var memStorage Repository
	memStorage.StorageGauge = Init()
	memStorage.StorageCounter = Init()
	memStorage.InterlvalSave = interlval
	memStorage.PathFile = path
	memStorage.loadStorage()
	return &memStorage
}

func (s *Repository) Delete(typeMetric string, key string) error {
	switch typeMetric {
	case "gauge":
		s.StorageGauge.Delete(key)
	case "counter":
		s.StorageCounter.Delete(key)
	default:
		return errors.New("type metric error")
	}
	return nil
}

func (s *Repository) Add(typeMetric string, key string, value string) error {
	switch typeMetric {
	case "gauge":
		valueFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		s.StorageGauge.Add(key, valueFloat)

	case "counter":
		valueInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		valueOldM, ok := s.StorageCounter.Get(key)
		if ok {
			valueInt += valueOldM.(int64)
		}

		s.StorageCounter.Add(key, valueInt)
	default:
		return errors.New("type metric error")
	}
	if s.InterlvalSave == 0 {
		s.SaveStorage()
	}
	return nil
}

func (s *Repository) AddJSON(data model.Metrics) error {
	switch data.MType {
	case "gauge":
		s.StorageGauge.Add(data.ID, *data.Value)
	case "counter":
		s.StorageCounter.Add(data.ID, *data.Delta)
	default:
		return errors.New("type metric error")
	}
	if s.InterlvalSave == 0 {
		s.SaveStorage()
	}
	return nil
}

func (s *Repository) GetJSON(data model.Metrics) (model.Metrics, error) {
	switch data.MType {
	case "gauge":
		valueM, ok := s.StorageGauge.Get(data.ID)
		if !ok {
			return data, errors.New("NotFound")
		}
		temp := valueM.(float64)
		data.Value = &temp
	case "counter":
		valueM, ok := s.StorageCounter.Get(data.ID)
		if !ok {
			return data, errors.New("NotFound")
		}
		temp := valueM.(int64)
		data.Delta = &temp
	default:
		return data, errors.New("type metric error")
	}
	return data, nil
}

func (s *Repository) Get(typeMetric string, key string) (interface{}, error) {
	switch typeMetric {
	case "gauge":
		valueM, err := s.StorageGauge.Get(key)
		if !err {
			return nil, errors.New("NotFound")
		}
		return valueM, nil

	case "counter":
		valueM, err := s.StorageCounter.Get(key)
		if !err {
			return nil, errors.New("NotFound")
		}
		return valueM, nil
	default:
		return nil, errors.New("type metric error")
	}
}

func (s *Repository) GetAll() []byte {
	resultData := []byte{}
	data := s.StorageCounter.GetAll()
	for name, value := range data {
		resultData = append(resultData, []byte(fmt.Sprintf("%v: %v\n", name, value.Value))...)
	}
	data = s.StorageGauge.GetAll()
	for name, value := range data {
		resultData = append(resultData, []byte(fmt.Sprintf("%v: %v\n", name, value.Value))...)
	}
	return resultData
}

func (m *Repository) AutoSaveStorage() {
	fmt.Print("Ale")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		fmt.Print("Ale")
		<-sigChan
		fmt.Print("Ale")
		m.SaveStorage()
		os.Exit(0)
	}()
	if m.InterlvalSave > 0 {
		for {
			time.Sleep(time.Duration(m.InterlvalSave) * time.Second)
			m.SaveStorage() // add error handler
		}
	}
}

func (m *Repository) SaveStorage() error {
	if m.PathFile == "" {
		return nil
	}
	m.StorageGauge.mutex.Lock()
	m.StorageCounter.mutex.Lock()
	defer m.StorageGauge.mutex.Unlock()
	defer m.StorageCounter.mutex.Unlock()
	preData := map[string]interface{}{
		"gauge":   m.StorageGauge.items,
		"counter": m.StorageCounter.items,
	}
	data, err := json.MarshalIndent(preData, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(m.PathFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (m *Repository) loadStorage() error {
	if _, err := os.Stat(m.PathFile); errors.Is(err, os.ErrNotExist) {
		return err
	}
	file, err := os.ReadFile(m.PathFile)
	if err != nil {
		return err
	}
	data := map[string]map[string]Item{}
	if err := json.Unmarshal(file, &data); err != nil {
		return err
	}
	m.StorageGauge.items = data["gauge"]
	m.StorageCounter.items = data["counter"]
	for i := range m.StorageCounter.items {
		temp := data["counter"][i].Value.(float64)
		m.StorageCounter.items[i] = Item{
			Value: int64(temp),
		}

	}
	return nil
}

func (m *Repository) Ping() error {
	return nil
}
