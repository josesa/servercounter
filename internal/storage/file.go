package storage

import (
	"os"
)

type FileStore struct {
	path string
}

func NewFileStorage(path string) *FileStore {
	fs := &FileStore{
		path: path,
	}

	return fs
}

func (s *FileStore) Read() (string, error) {

	content, err := os.ReadFile(s.path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (fs *FileStore) Write(s string) error {
	return os.WriteFile(fs.path, []byte(s), 0666)
}
