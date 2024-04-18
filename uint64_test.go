package unique

import (
	"math"
	"sync"
	"testing"
)

func TestUint64(t *testing.T) {

	u := &Uint64Counter{}

	const goroutineCount = 500_000
	var wg sync.WaitGroup
	wg.Add(1)
	out := make(chan uint64, goroutineCount)
	expect := make(map[uint64]struct{}, goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		expect[uint64(i)] = struct{}{}
		go func(i int) {
			wg.Wait()
			uniqueNumber := u.Get()
			out <- uniqueNumber
		}(i)
	}
	wg.Done()
	for i := 0; i < goroutineCount; i++ {
		delete(expect, <-out)
	}

	if len(expect) != 0 {
		t.Fatalf("expected every number to be unique, but there were %d collisions",
			len(expect))
	}
}

func TestUint64Rollover(t *testing.T) {

	u := &Uint64Counter{}
	u.counter.Store(math.MaxUint64 - 50)

	const goroutineCount = 500_000
	var wg sync.WaitGroup
	wg.Add(1)
	out := make(chan uint64, goroutineCount)
	expect := make(map[uint64]struct{}, goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		expect[uint64(i)+math.MaxUint64-50] = struct{}{}
		go func(i int) {
			wg.Wait()
			uniqueNumber := u.Get()
			out <- uniqueNumber
		}(i)
	}
	wg.Done()
	for i := 0; i < goroutineCount; i++ {
		delete(expect, <-out)
	}

	if len(expect) != 0 {
		t.Fatalf("expected every number to be unique, but there were %d collisions",
			len(expect))
	}
}
