package mem

import (
	"sync"
)

type MemStorage struct {
	mutex sync.Mutex
	items map[string]Item
}

type Item struct {
	Value interface{}
}

func Init() *MemStorage {
	items := make(map[string]Item)
	return &MemStorage{
		mutex: sync.Mutex{},
		items: items,
	}
}

func (s *MemStorage) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.items, key)
}

func (s *MemStorage) Add(key string, value interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = Item{
		Value: value,
	}
}

func (s *MemStorage) Get(key string) (interface{}, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	item, found := s.items[key]
	if !found {
		return nil, false
	}
	return item.Value, true
}

func (s *MemStorage) GetAll() map[string]Item {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tempItems := make(map[string]Item, len(s.items))
	for key, val := range s.items {
		tempItems[key] = val
	}
	return tempItems
}
