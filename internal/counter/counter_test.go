package counter

import (
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

type fakeTime struct {
	sync.Mutex
	fake time.Time
}

func createFakeTime() *fakeTime {
	ft := &fakeTime{
		fake: time.Date(1999, time.January, 10, 1, 2, 3, 4, time.UTC),
	}

	return ft
}
func (ft *fakeTime) Now() time.Time {
	ft.Lock()
	defer ft.Unlock()
	return ft.fake
}

func (ft *fakeTime) IncreaseSecond(s int) time.Time {
	ft.Lock()
	defer ft.Unlock()
	ft.fake = ft.fake.Add(time.Duration(s) * time.Second)
	return ft.fake
}

func Test_CounterWindow(t *testing.T) {

	fakeTime := createFakeTime()
	uat := New(WithTime(fakeTime), WithWindowSize(5))

	uat.Inc() // Second 0
	assertCount(uat, 1, "Invalid count returned in first interval", t)

	uat.Inc() // Second 0
	assertCount(uat, 2, "Invalid count returned in second interval", t)

	fakeTime.IncreaseSecond(1) // Step time, Second 1

	uat.Inc() // Second 1
	assertCount(uat, 3, "Invalid count returned in third interval", t)

	fakeTime.IncreaseSecond(4) // Step time, first Increases on time 0 are now outside of evaluation window
	assertCount(uat, 1, "Invalid count returned after window expiration", t)

	fakeTime.IncreaseSecond(1) // Step time, all values should be outside evaluation window
	assertCount(uat, 0, "Invalid count after full window expiration", t)
}

// Important to note that these operations are concurrent, not parallel.
func Test_CounterCount(t *testing.T) {

	windowSize := 5
	fakeTime := createFakeTime()
	uat := New(WithTime(fakeTime), WithWindowSize(windowSize))
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		for i := 0; i < 50; i++ {
			uat.Inc()
			if i%10 == 0 {
				runtime.Gosched() // Yields processing to other goroutines
			}

		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 50; i++ {
			uat.Inc()
			if i%10 == 0 {
				runtime.Gosched() // Yields processing to other goroutines
			}
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < windowSize-1; i++ {
			fakeTime.IncreaseSecond(1)
			time.Sleep(2 * time.Microsecond) // To prevent processing hogging on the current routine. Yields to other goroutines
		}
		wg.Done()
	}()
	wg.Wait() // Wait until all loops are completed
	t.Log(uat.Dump())
	assertCount(uat, 100, "", t)
}

func Test_CounterLoad(t *testing.T) {
	windowSize := 5
	fakeTime := createFakeTime()
	uat := New(WithTime(fakeTime), WithWindowSize(windowSize))

	seed := map[int64]uint64{
		915930124: 1,
		915930125: 2,
		915930126: 1,
		915930127: 2,
	}

	uat.Load(seed)
	assertCount(uat, 6, "Invalid loaded count", t)

	seed = map[int64]uint64{
		915930124: 1,
		915930125: 2,
		915930126: 1,
		915930127: 2,
	}
	uat.Load(seed)

	fakeTime.IncreaseSecond(windowSize + 2) // Exclude first 2 seconds entries

	assertCount(uat, 3, "Invalid loaded count", t)
}

func Test_CounterDump(t *testing.T) {
	windowSize := 5
	fakeTime := createFakeTime()
	uat := New(WithTime(fakeTime), WithWindowSize(windowSize))

	seed := map[int64]uint64{
		915930124: 1,
		915930125: 2,
		915930126: 1,
		915930127: 2,
	}

	uat.Load(seed)
	dump := uat.Dump()

	if !reflect.DeepEqual(dump, seed) {
		t.Errorf("Non matching dump from seed. Expected %v, Got %v", seed, dump)
	}
}

func assertCount(uat Counter, expected uint64, msg string, t *testing.T) {
	count := uat.GetCount()
	if count != expected {
		t.Errorf("%s: Expected %d, Got %d", msg, expected, count)
	}
}
