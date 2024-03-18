package service

import (
	"strings"
	"testing"

	"github.com/josesa/servercounter/internal/storage"
)

type mockCounter struct {
	fakeTime int64
	data     map[int64]uint64
}

func (mc *mockCounter) Inc() {
	mc.data[mc.fakeTime]++
}

func (mc *mockCounter) GetCount() uint64 {
	var count uint64
	for _, v := range mc.data {
		count = count + v
	}

	return count
}

func (mc *mockCounter) Dump() map[int64]uint64 {
	return mc.data
}

func (mc *mockCounter) Load(d map[int64]uint64) error {
	mc.data = d
	return nil
}

func Test_CounterServiceLoad(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	fakeCounter := mockCounter{
		fakeTime: 915930127, //Fixes time
	}

	seed := map[int64]uint64{
		915930127: 1,
	}

	saveToStorage(seed, memoryStorage)

	uat, err := New(&fakeCounter, memoryStorage)
	if err != nil {
		t.Error(err)
	}

	count := uat.IncAndGetCount()
	if count != 2 {
		t.Errorf("Expected %d, Got %d", 2, count)
	}
}

type testCase struct {
	seed   string
	expect string
	action func(*HitCounter)
}

func Test_CounterServiceSave(t *testing.T) {
	memoryStorage := storage.NewMemoryStorage()
	fakeCounter := mockCounter{
		fakeTime: 915930127, //Fixes time
	}

	testCases := map[string]testCase{
		"single increase": {
			seed:   "915930127:1\n",
			expect: "915930127:2\n",
			action: func(uat *HitCounter) {
				uat.IncAndGetCount()
				uat.Flush()
			},
		},
		"increase different time": {
			seed:   "915930127:1\n",
			expect: "915930127:1\n915930128:1\n",
			action: func(uat *HitCounter) {
				fakeCounter.fakeTime = 915930128
				uat.IncAndGetCount()
				uat.Flush()
			},
		},
	}

	for title, c := range testCases {
		t.Run(title, func(t *testing.T) {
			memoryStorage.Write(c.seed)
			uat, err := New(&fakeCounter, memoryStorage)
			if err != nil {
				t.Error(err)
			}
			c.action(uat)

			content, _ := memoryStorage.Read()
			if strings.Compare(content, c.expect) != 0 {
				t.Errorf("Case %s failed. Expected %v, Got %v", title, c.expect, content)
			}
		})

	}
}
