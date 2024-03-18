package counter

import (
	"sync"
	"time"
)

const DEFAULTSIZE = 60
const DEFAULTFLUSHINTERVAL = DEFAULTSIZE * 1.1 * time.Second // 10% higher than Default Size

type mwcounter struct {
	mu sync.RWMutex
	// key holds unix time seconds, value holds counter for that second
	data map[int64]uint64

	windowSizeSeconds int
	flushInterval     time.Duration

	timeHandler timeHandler
}

type Counter interface {
	Inc()
	GetCount() uint64

	Dump() map[int64]uint64
	Load(map[int64]uint64) error
}

func New(options ...func(*mwcounter)) *mwcounter {
	// Initializes with sane defaults
	counter := &mwcounter{
		data:              make(map[int64]uint64, DEFAULTSIZE),
		windowSizeSeconds: DEFAULTSIZE,
		flushInterval:     DEFAULTFLUSHINTERVAL,
		timeHandler:       realClock{},
	}
	// Apply user options
	for _, o := range options {
		o(counter)
	}

	go func() {
		// Periodically clean the old entries
		cleanerTicker := time.NewTicker(counter.flushInterval)
		for range cleanerTicker.C {
			counter.mu.Lock()
			counter.clean()
			counter.mu.Unlock()
		}
	}()

	return counter
}

type timeHandler interface {
	Now() time.Time
}

// realClock holds a placeholder structure for time handling. It implements the native time handler time.Now()
type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

func WithWindowSize(size int) func(*mwcounter) {
	return func(mw *mwcounter) {
		mw.windowSizeSeconds = size
		mw.data = make(map[int64]uint64, size)
	}
}

// WithTime optionally sets a different time handler instead of the default time.Now().
func WithTime(t timeHandler) func(*mwcounter) {
	return func(mw *mwcounter) {
		mw.timeHandler = t
	}
}

// WithFlushInterval sets the duration between clean operation. It should be larger than window size to ensure that there is no data loss.
func WithFlushInterval(t time.Duration) func(*mwcounter) {
	return func(mw *mwcounter) {
		mw.flushInterval = t
	}
}

// Inc increases the counter
func (c *mwcounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Get current second
	now := c.timeHandler.Now().Unix()
	c.data[now]++
}

// GetCount returns the total count stored in the counter within the defined Window size.
func (c *mwcounter) GetCount() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var count uint64

	earliest := c.timeHandler.Now().Add(time.Duration(-c.windowSizeSeconds) * time.Second).Unix()
	// Check earliest time and only count within window
	for t, v := range c.data {
		if t > earliest {
			count = count + v
		}
	}
	return count
}

// Dump returns the current dump of stored values
func (c *mwcounter) Dump() map[int64]uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.clean()
	return c.data
}

// Load initializes it internal storage to the supplied data
func (c *mwcounter) Load(m map[int64]uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = m
	c.clean() // Deletes no longer valid entries
	return nil
}

// clean should be triggered periodically to delete old entries that are no longer valid for the evaluation window
func (c *mwcounter) clean() {
	earliest := c.timeHandler.Now().Add(time.Duration(-c.windowSizeSeconds) * time.Second).Unix()
	for k := range c.data {
		if k < earliest {
			delete(c.data, k)
		}
	}
}
