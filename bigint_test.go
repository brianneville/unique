package unique

import (
	"math"
	"math/big"
	"sync"
	"testing"
)

func TestInfiniteNumber(t *testing.T) {

	b := &BigIntCounter{}

	startCounter := big.NewInt(0)
	const goroutineCount = 500_000
	var wg sync.WaitGroup
	wg.Add(1)
	out := make(chan *big.Int, goroutineCount)
	expect := make(map[string]struct{}, goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		expect[startCounter.String()] = struct{}{}
		go func(i int) {
			wg.Wait()
			uniqueNumber := b.Get()
			out <- uniqueNumber
		}(i)
	}
	wg.Done()
	for i := 0; i < goroutineCount; i++ {
		got := <-out
		delete(expect, got.String())
	}

	if len(expect) != 0 {
		t.Fatalf("expected every number to be unique, but there were %d collisions",
			len(expect))
	}
}

func TestInfiniteNumberRollover(t *testing.T) {

	b := &BigIntCounter{}

	const startPoint = math.MaxInt64 - 5
	b.infCounterBase.Store(startPoint)
	startCounter := big.NewInt(startPoint)
	const goroutineCount = 500_000
	var wg sync.WaitGroup
	wg.Add(1)
	out := make(chan *big.Int, goroutineCount)
	expect := make(map[string]struct{}, goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		expect[startCounter.String()] = struct{}{}
		go func(i int) {
			wg.Wait()
			uniqueNumber := b.Get()
			out <- uniqueNumber
		}(i)
	}
	wg.Done()
	for i := 0; i < goroutineCount; i++ {
		got := <-out
		delete(expect, got.String())
	}

	if len(expect) != 0 {
		t.Fatalf("expected every number to be unique, but there were %d collisions",
			len(expect))
	}
}
