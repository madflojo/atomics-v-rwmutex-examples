package simple

import (
	"errors"
	"sync"
)

// ErrInsufficientFunds indicates a withdrawal would push the balance negative.
var ErrInsufficientFunds = errors.New("insufficient funds")

// MutexSimpleBalance uses a standard Mutex to guard just the balance value.
type MutexSimpleBalance struct {
	mu    sync.Mutex
	value int64
}

// New constructs a zeroed MutexSimpleBalance.
func New() *MutexSimpleBalance { return &MutexSimpleBalance{} }

// Balance returns the current value under a lock.
func (b *MutexSimpleBalance) Balance() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.value
}

// TransactionCount always returns zero because the simple variant does not track metadata.
func (b *MutexSimpleBalance) TransactionCount() int64 { return 0 }

// LastUpdated always reports zero because timestamps are not recorded.
func (b *MutexSimpleBalance) LastUpdated() int64 { return 0 }

// Add increments the value with exclusive access.
func (b *MutexSimpleBalance) Add(amount int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.value += amount
}

// Subtract decrements the value or returns ErrInsufficientFunds.
func (b *MutexSimpleBalance) Subtract(amount int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.value-amount < 0 {
		return ErrInsufficientFunds
	}
	b.value -= amount
	return nil
}
