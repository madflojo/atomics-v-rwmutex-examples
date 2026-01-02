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
	value atomic.Int64
	// trx counts successful mutations.
	trx atomic.Int64
	// updated records the timestamp of the latest mutation.
	updated atomic.Int64
}

// New creates a zeroed AtomicCASFullBalance.
func New() *AtomicCASFullBalance {
	return &AtomicCASFullBalance{}
}

// Balance returns the current value.
func (b *AtomicCASFullBalance) Balance() int64 {
	return b.value.Load()
}

// TransactionCount reports completed mutations.
func (b *AtomicCASFullBalance) TransactionCount() int64 {
	return b.trx.Load()
}

// LastUpdated returns the timestamp for the latest mutation.
func (b *AtomicCASFullBalance) LastUpdated() int64 {
	return b.updated.Load()
}

// Add increments the value and metadata.
func (b *AtomicCASFullBalance) Add(amount int64) {
	b.value.Add(amount)
	b.trx.Add(1)
	b.updated.Store(time.Now().UnixNano())
}

// Subtract decrements the value via CAS and records metadata updates.
func (b *AtomicCASFullBalance) Subtract(amount int64) error {
	for {
		current := b.value.Load()
		next := current - amount
		if next < 0 {
			return ErrInsufficientFunds
		}

		if b.value.CompareAndSwap(current, next) {
			b.trx.Add(1)
			b.updated.Store(time.Now().UnixNano())
			return nil
		}
	}
}
