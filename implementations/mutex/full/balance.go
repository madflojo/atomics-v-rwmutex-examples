package full

import (
	"errors"
	"sync"
	"time"
)

// ErrInsufficientFunds indicates the balance would go negative.
var ErrInsufficientFunds = errors.New("insufficient funds")

// MutexFullBalance protects balance metadata with a standard Mutex.
type MutexFullBalance struct {
	mu      sync.Mutex
	value   int64
	trx     int64
	updated int64
}

// New returns a zeroed MutexFullBalance.
func New() *MutexFullBalance { return &MutexFullBalance{} }

// Balance returns the current value under a lock.
func (b *MutexFullBalance) Balance() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.value
}

// TransactionCount returns how many mutations have executed.
func (b *MutexFullBalance) TransactionCount() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.trx
}

// LastUpdated returns the timestamp of the latest mutation.
func (b *MutexFullBalance) LastUpdated() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.updated
}

// Add increments the balance and records metadata.
func (b *MutexFullBalance) Add(amount int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.value += amount
	b.trx++
	b.updated = time.Now().UnixNano()
}

// Subtract decrements the balance or returns ErrInsufficientFunds.
func (b *MutexFullBalance) Subtract(amount int64) error {
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
