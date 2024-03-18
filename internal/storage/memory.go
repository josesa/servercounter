package storage

type MemoryStore struct {
	content string
}

func NewMemoryStorage() *MemoryStore {
	ms := &MemoryStore{}

	return ms
}

func (ms *MemoryStore) Read() (string, error) {
	return string(ms.content), nil
}

func (ms *MemoryStore) Write(s string) error {
	ms.content = s

	return nil
}
