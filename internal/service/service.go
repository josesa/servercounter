package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/josesa/servercounter/internal/counter"
	"github.com/josesa/servercounter/internal/storage"
)

// HitCounter provides a persistance hit counter service
type HitCounter struct {
	counter counter.Counter

	store storage.Storage
}

// New creates a HitCounter service
func New(c counter.Counter, storage storage.Storage) (*HitCounter, error) {

	data, err := readFromStorage(storage)
	if err != nil {
		// Is not possible to load the data from storage, in this case, proceed with empty data instead of failing.
		// Alternative, the service could abort if that information is considered critical
		data = make(map[int64]uint64)
	}

	if len(data) > 0 {
		c.Load(data)
	}

	wc := &HitCounter{
		counter: c,
		store:   storage,
	}

	return wc, nil
}

// Flush should be called when terminating the service if the counter values should be persisted
func (hc *HitCounter) Flush() error {
	data := hc.counter.Dump()
	return saveToStorage(data, hc.store)
}

// IncAndGetCount increases the counter value and returns current count
func (hc *HitCounter) IncAndGetCount() uint64 {
	hc.counter.Inc()

	return hc.counter.GetCount()
}

func readFromStorage(s storage.Storage) (map[int64]uint64, error) {
	data := make(map[int64]uint64)
	content, err := s.Read()
	if err != nil {
		return data, err
	}

	lines := strings.Split(content, "\n")
	for _, l := range lines {
		parts := strings.Split(l, ":")
		if len(parts) == 2 {
			k, errk := strconv.ParseInt(parts[0], 10, 64)
			v, errv := strconv.ParseUint(parts[1], 10, 64)
			if errk == nil && errv == nil {
				data[k] = v
			}
		}
	}

	return data, nil
}

func saveToStorage(data map[int64]uint64, s storage.Storage) error {
	var sb strings.Builder
	for k, v := range data {
		sb.WriteString(fmt.Sprintf("%d:%d\n", k, v))
	}

	return s.Write(sb.String())
}
