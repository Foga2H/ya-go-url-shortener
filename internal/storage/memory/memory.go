package storage

type MemStorage struct {
	links map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		links: make(map[string]string),
	}
}

func (m *MemStorage) Set(key string, value string) {
	m.links[key] = value
}

func (m *MemStorage) Get(key string) (string, bool) {
	val, ok := m.links[key]
	return val, ok
}
