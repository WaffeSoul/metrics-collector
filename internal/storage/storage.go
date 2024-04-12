package storage

type MemStorage struct {
	StorageGause  *Storage
	StorageConter *Storage
}

type Storage struct {
	items map[string]Item
}

type Item struct {
	Value interface{}
}

func InitMem() *MemStorage {
	var memStorage MemStorage
	memStorage.StorageGause = Init()
	memStorage.StorageConter = Init()
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
	for i, a := range s.items {
		println(i, a.Value)
	}
	return s.items
}
