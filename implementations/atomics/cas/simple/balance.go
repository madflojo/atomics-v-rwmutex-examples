package simple

import (
	"errors"
	"sync/atomic"
)

// ErrInsufficientFunds indicates the balance would become negative.
var ErrInsufficientFunds = errors.New("insufficient funds")

// AtomicCASSimpleBalance keeps only the balance value while using CAS to
// ensure atomic read-modify-write semantics.
type AtomicCASSimpleBalance struct {
	// value stores the running balance.
	value atomic.Int64
}

// New returns a zeroed AtomicCASSimpleBalance.
func New() *AtomicCASSimpleBalance {
	return &AtomicCASSimpleBalance{}
}

// Balance returns the current value.
func (b *AtomicCASSimpleBalance) Balance() int64 {
	return b.value.Load()
}

// TransactionCount always returns zero because the simple variant does
// not track counts.
func (b *AtomicCASSimpleBalance) TransactionCount() int64 {
	return 0
}

// LastUpdated always returns zero because timestamps are not tracked.
func (b *AtomicCASSimpleBalance) LastUpdated() int64 {
	return 0
}

// Add increments the balance using atomic addition.
func (b *AtomicCASSimpleBalance) Add(amount int64) {
	b.value.Add(amount)
}

// Subtract decrements the balance while guaranteeing the update via CAS.
func (b *AtomicCASSimpleBalance) Subtract(amount int64) error {
	for {
		current := b.value.Load()
		next := current - amount
		if next < 0 {
			return ErrInsufficientFunds
		}

		if b.value.CompareAndSwap(current, next) {
			return nil
		}
	}
}
