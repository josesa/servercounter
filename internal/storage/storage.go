package storage

type Storage interface {
	Read() (string, error)
	Write(string) error
}

type Store struct {
}
