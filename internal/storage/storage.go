package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type MemStorage struct {
	StorageGauge   *Storage
	StorageCounter *Storage
	InterlvalSave  int
	PathFile       string
}

type Storage struct {
	mutex sync.Mutex
	items map[string]Item
}

type Item struct {
	Value interface{}
}

func InitMem(interlval int, path string) *MemStorage {
	var memStorage MemStorage
	memStorage.StorageGauge = Init()
	memStorage.StorageCounter = Init()
	memStorage.InterlvalSave = interlval
	memStorage.PathFile = path
	return &memStorage
}

func Init() *Storage {
	items := make(map[string]Item)
	return &Storage{
		mutex: sync.Mutex{},
		items: items,
	}

}

func (s *Storage) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.items, key)
}

func (s *Storage) Add(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = Item{
		Value: value,
	}
}

func (s *Storage) Get(key string) (interface{}, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	item, found := s.items[key]
	if !found {
		return nil, false
	}
	return item.Value, true
}

func (s *Storage) GetAll() map[string]Item {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tempItems := make(map[string]Item, len(s.items))
	for key, val := range s.items {
		tempItems[key] = val
	}
	return tempItems
}

func (m *MemStorage) AutoSaveStorage() {
	for {
		time.Sleep(time.Duration(m.InterlvalSave) * time.Second)
		m.SaveStorage() // add error handler
	}

}

func (m *MemStorage) SaveStorage() error {
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

func (m *MemStorage) LoadStorage() error {
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
