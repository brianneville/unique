package unique

import "sync/atomic"

type Uint64Counter struct {
	counter atomic.Uint64
}

// Get returns a uint64 which will be larger than any
// uint64 returned from a previous call to Get.
// Specifically, the value will increase by 1 on every call (until rollover).
// Get can be called from multiple goroutines concurrently, and will
// not return duplicate values.
// Note however, this counter will roll over to 0 at uint64 max value
func (u *Uint64Counter) Get() uint64 {
	var current uint64
	for {
		current = u.counter.Load()
		next := current + 1
		if u.counter.CompareAndSwap(current, next) {
			break
		}
	}
	return current
}
