package full

import (
	"errors"
	"sync/atomic"
	"time"
)

// ErrInsufficientFunds indicates the balance would drop below zero.
var ErrInsufficientFunds = errors.New("insufficient funds")

// AtomicCASFullBalance stores balance metadata while protecting every
// update via CAS loops.
type AtomicCASFullBalance struct {
	// value holds the running balance.
	value int64
	// trx counts successful mutations.
	trx int64
	// updated records the timestamp of the latest mutation.
	updated int64
}

// New creates a zeroed AtomicCASFullBalance.
func New() (*AtomicCASFullBalance, error) {
	return &AtomicCASFullBalance{}, nil
}

// Balance returns the current value.
func (b *AtomicCASFullBalance) Balance() int64 {
	return atomic.LoadInt64(&b.value)
}

// TransactionCount reports completed mutations.
func (b *AtomicCASFullBalance) TransactionCount() int64 {
	return atomic.LoadInt64(&b.trx)
}

// LastUpdated returns the timestamp for the latest mutation.
func (b *AtomicCASFullBalance) LastUpdated() int64 {
	return atomic.LoadInt64(&b.updated)
}

// Add increments the value and metadata.
func (b *AtomicCASFullBalance) Add(amount int64) {
	atomic.AddInt64(&b.value, amount)
	atomic.AddInt64(&b.trx, 1)
	atomic.StoreInt64(&b.updated, time.Now().UnixNano())
}

// Subtract decrements the value via CAS and records metadata updates.
func (b *AtomicCASFullBalance) Subtract(amount int64) error {
	for {
		current := atomic.LoadInt64(&b.value)
		next := current - amount
		if next < 0 {
			return ErrInsufficientFunds
		}

		if atomic.CompareAndSwapInt64(&b.value, current, next) {
			atomic.AddInt64(&b.trx, 1)
			atomic.StoreInt64(&b.updated, time.Now().UnixNano())
			return nil
		}
	}
}
