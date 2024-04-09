package storage

var (
	StorageGause  *Storage
	StorageConter *Storage
)

type Storage struct {
	items map[string]Item
}

type Item struct {
	Value interface{}
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
	return item, true
}
