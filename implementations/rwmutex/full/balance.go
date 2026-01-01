package full

import (
	"errors"
	"sync"
	"time"
)

// ErrInsufficientFunds indicates the balance would go negative.
var ErrInsufficientFunds = errors.New("insufficient funds")

// RWMutexFullBalance protects balance metadata with an RWMutex while
// allowing concurrent reads.
type RWMutexFullBalance struct {
	// mu guards all fields.
	mu sync.RWMutex
	// value stores the running balance.
	value int64
	// trx counts successful mutations.
	trx int64
	// updated records the timestamp of the most recent mutation.
	updated int64
}

// New returns a zeroed RWMutexFullBalance.
func New() *RWMutexFullBalance {
	return &RWMutexFullBalance{}
}

// Balance returns the current value under a read lock.
func (b *RWMutexFullBalance) Balance() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.value
}

// TransactionCount returns how many mutations have executed.
func (b *RWMutexFullBalance) TransactionCount() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.trx
}

// LastUpdated returns the timestamp of the latest mutation.
func (b *RWMutexFullBalance) LastUpdated() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.updated
}

// Add increments the balance and records metadata.
func (b *RWMutexFullBalance) Add(amount int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.value += amount
	b.trx++
	b.updated = time.Now().UnixNano()
}

// Subtract decrements the balance or returns ErrInsufficientFunds.
func (b *RWMutexFullBalance) Subtract(amount int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.value-amount < 0 {
		return ErrInsufficientFunds
	}

	b.value -= amount
	b.trx++
	b.updated = time.Now().UnixNano()
	return nil
}
