package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type MemStorage struct {
	StorageGauge   *Storage
	StorageCounter *Storage
	InterlvalSave  int
	LastSave       time.Time
	PathFile       string
}

type Storage struct {
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
	memStorage.LastSave = time.Time{}
	return &memStorage
}

func Init() *Storage {
	items := make(map[string]Item)
	return &Storage{
		items: items,
	}

}

func (s *Storage) Delete(key string) {
	delete(s.items, key)
}

func (s *Storage) Add(key string, value interface{}) {
	s.items[key] = Item{
		Value: value,
	}
}

func (s *Storage) Get(key string) (interface{}, bool) {
	item, found := s.items[key]
	if !found {
		return nil, false
	}
	return item.Value, true
}

func (s *Storage) GetAll() map[string]Item {
	return s.items
}

func (m *MemStorage) SaveStorage() {
	if m.LastSave.Add(time.Duration(m.InterlvalSave)*time.Second).After(time.Now()) || m.InterlvalSave == 0 {
		return
	}
	m.LastSave = time.Now()
	file, err := os.OpenFile(m.PathFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	var temp []byte
	_, err = file.Read(temp)
	if err != nil {
		return
	}
	preData := map[string]interface{}{
		"gauge":   m.StorageGauge.items,
		"counter": m.StorageCounter.items,
	}
	data, err := json.MarshalIndent(preData, "", "    ")
	if err != nil {
		panic(err)
	}
	file.Write(data)
}

func (m *MemStorage) LoadStorage() {
	if _, err := os.Stat(m.PathFile); errors.Is(err, os.ErrNotExist) {
		return
	}
	file, err := os.ReadFile(m.PathFile)
	if err != nil {
		return
	}
	fmt.Println(file)
	data := map[string]map[string]Item{}
	if err := json.Unmarshal(file, &data); err != nil {
		panic(err)
	}
	m.StorageGauge.items = data["gauge"]
	m.StorageCounter.items = data["counter"]
}
