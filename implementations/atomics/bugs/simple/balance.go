package simple

import (
	"errors"
	"sync/atomic"
	"time"
)

// ErrInsufficientFunds indicates a subtraction would push the balance below zero.
var ErrInsufficientFunds = errors.New("insufficient funds")

// AtomicBugsSimpleBalance is an intentionally incorrect atomic balance that only
// tracks its value, making it easy to demonstrate race conditions.
type AtomicBugsSimpleBalance struct {
	// value stores the raw account balance.
	value int64
}

// New creates a zeroed AtomicBugsSimpleBalance.
func New() (*AtomicBugsSimpleBalance, error) {
	return &AtomicBugsSimpleBalance{}, nil
}

// Balance returns the current value.
func (b *AtomicBugsSimpleBalance) Balance() int64 {
	return atomic.LoadInt64(&b.value)
}

// TransactionCount always reports zero because the simple implementation
// does not track metadata.
func (b *AtomicBugsSimpleBalance) TransactionCount() int64 {
	return 0
}

// LastUpdated always returns zero because timestamps are not recorded.
func (b *AtomicBugsSimpleBalance) LastUpdated() int64 {
	return 0
}

// Add increments the value atomically.
func (b *AtomicBugsSimpleBalance) Add(amount int64) {
	atomic.AddInt64(&b.value, amount)
}

// Subtract decrements the value without CAS protection, intentionally
// leaving room for lost updates under contention.
func (b *AtomicBugsSimpleBalance) Subtract(amount int64) error {
	current := atomic.LoadInt64(&b.value)
	time.Sleep(100 * time.Microsecond)
	if current-amount < 0 {
		return ErrInsufficientFunds
	}

	atomic.AddInt64(&b.value, -amount)
	return nil
}
