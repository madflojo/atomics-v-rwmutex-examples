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
	value int64
}

// New returns a zeroed AtomicCASSimpleBalance.
func New() (*AtomicCASSimpleBalance, error) {
	return &AtomicCASSimpleBalance{}, nil
}

// Balance returns the current value.
func (b *AtomicCASSimpleBalance) Balance() int64 {
	return atomic.LoadInt64(&b.value)
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
	atomic.AddInt64(&b.value, amount)
}

// Subtract decrements the balance while guaranteeing the update via CAS.
func (b *AtomicCASSimpleBalance) Subtract(amount int64) error {
	for {
		current := atomic.LoadInt64(&b.value)
		next := current - amount
		if next < 0 {
			return ErrInsufficientFunds
		}

		if atomic.CompareAndSwapInt64(&b.value, current, next) {
			return nil
		}
	}
}
