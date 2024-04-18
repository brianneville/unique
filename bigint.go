package unique

import (
	"math"
	"math/big"
	"sync/atomic"
)

type BigIntCounter struct {
	infCounterBase atomic.Uint64
	infCounter     big.Int
}

// Get returns a *big.Int which will always be larger than the
// any *big.Int returned from a previous call to Get.
// Specifically, the value will increase by 1 on every call.
// Get can be called from multiple goroutines concurrently, and will
// not return duplicate values.
// This counter will never roll over to 0.
func (b *BigIntCounter) Get() *big.Int {
	var current uint64
	var snapshot []big.Word
	for {
		current = b.infCounterBase.Load()
		snapshot = b.infCounter.Bits()[:] // take a view, in case this increases in size after
		next := current + 1
		if current > math.MaxInt64 {
			// spin in the loop until the one goroutine which had got the maxInt64
			// reaches the end of the loop
			continue
		}
		if b.infCounterBase.CompareAndSwap(current, next) {
			break
		}
	}

	// When this happens, we've hit the int64 max value,
	// so add that max value into the big int
	if current == math.MaxInt64 {
		b.infCounter.Add(&b.infCounter, big.NewInt(math.MaxInt64))
		bigResult := (&big.Int{}).SetBits(big.NewInt(math.MaxInt64).Bits())
		b.infCounterBase.Store(1)
		return bigResult
	}

	bigResult := (&big.Int{}).SetBits(snapshot)
	bigResult.Add(bigResult, big.NewInt(int64(current)))
	return bigResult
}
