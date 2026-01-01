package simple

import (
	"errors"
	"sync"
)

// ErrInsufficientFunds indicates a withdrawal would push the balance negative.
var ErrInsufficientFunds = errors.New("insufficient funds")

// RWMutexSimpleBalance uses an RWMutex to guard just the balance value.
type RWMutexSimpleBalance struct {
	// mu protects value.
	mu sync.RWMutex
	// value stores the running balance.
	value int64
}

// New constructs a zeroed RWMutexSimpleBalance.
func New() (*RWMutexSimpleBalance, error) {
	return &RWMutexSimpleBalance{}, nil
}

// Balance returns the current value under a read lock.
func (b *RWMutexSimpleBalance) Balance() int64 {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.value
}

// TransactionCount always returns zero because the simple variant does not
// track metadata.
func (b *RWMutexSimpleBalance) TransactionCount() int64 {
	return 0
}

// LastUpdated always reports zero because timestamps are not recorded.
func (b *RWMutexSimpleBalance) LastUpdated() int64 {
	return 0
}

// Add increments the value with exclusive access.
func (b *RWMutexSimpleBalance) Add(amount int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.value += amount
}

// Subtract decrements the value or returns ErrInsufficientFunds.
func (b *RWMutexSimpleBalance) Subtract(amount int64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.value-amount < 0 {
		return ErrInsufficientFunds
	}

	b.value -= amount
	return nil
}
