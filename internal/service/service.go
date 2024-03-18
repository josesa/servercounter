package service

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/josesa/servercounter/internal/counter"
)

// HitCounter provides a persistance hit counter service
type HitCounter struct {
	counter counter.Counter

	dataStoragePath string
}

func New(c counter.Counter, path string) (*HitCounter, error) {
	// If a dataStoragePath is supplied, then it is used to load the initial values to the counter
	if len(path) > 0 {
		data, err := readFromFile(path)
		if err != nil {
			return nil, err
		}

		if len(data) > 0 {
			c.Load(data)
		}
	}

	wc := &HitCounter{
		counter:         c,
		dataStoragePath: path,
	}

	return wc, nil
}

// Shutdown should be called when terminating the service if the counter values should be persisted
func (hc *HitCounter) Shutdown() error {
	// On shutdown, save the current counter values into permanent storage

	if len(hc.dataStoragePath) > 0 {
		data := hc.counter.Dump()
		saveToFile(data, hc.dataStoragePath)
	}

	return nil
}

// IncAndGetCount increases the counter value and returns current count
func (hc *HitCounter) IncAndGetCount() uint64 {
	hc.counter.Inc()
	return hc.counter.GetCount()
}

// loadFromFile tries to read counter information from a text file
func readFromFile(path string) (map[int64]uint64, error) {
	data := make(map[int64]uint64)

	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer file.Close()

	buffer := bufio.NewScanner(file)
	buffer.Split(bufio.ScanLines)
	for buffer.Scan() {
		parts := strings.Split(buffer.Text(), ":")
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

// saveToFile extracts current counter values and dumps them into the text file
func saveToFile(data map[int64]uint64, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := bufio.NewWriter(file)
	for k, v := range data {
		fmt.Fprintf(buffer, "%d:%d\n", k, v)
	}
	buffer.Flush()

	return nil
}
