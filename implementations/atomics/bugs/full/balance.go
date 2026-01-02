package full

import (
	"errors"
	"sync/atomic"
	"time"
)

// ErrInsufficientFunds signals that a subtract would create a negative balance.
var ErrInsufficientFunds = errors.New("insufficient funds")

// AtomicBugsFullBalance mirrors the feature set of the other full implementations
// but intentionally omits CAS protection to highlight logic races.
type AtomicBugsFullBalance struct {
	// value stores the running balance.
	value atomic.Int64
	// trx counts successful mutations.
	trx atomic.Int64
	// updated records the timestamp of the last mutation in nanoseconds.
	updated atomic.Int64
}

// New constructs a zeroed AtomicBugsFullBalance.
func New() *AtomicBugsFullBalance {
	return &AtomicBugsFullBalance{}
}

// Balance returns the current value.
func (b *AtomicBugsFullBalance) Balance() int64 {
	return b.value.Load()
}

// TransactionCount reports how many writes have completed.
func (b *AtomicBugsFullBalance) TransactionCount() int64 {
	return b.trx.Load()
}

// LastUpdated returns the timestamp for the most recent mutation.
func (b *AtomicBugsFullBalance) LastUpdated() int64 {
	return b.updated.Load()
}

// Add increments the balance and metadata without locking.
func (b *AtomicBugsFullBalance) Add(amount int64) {
	b.value.Add(amount)
	b.trx.Add(1)
	b.updated.Store(time.Now().UnixNano())
}

// Subtract decrements the balance but intentionally lacks CAS protection,
// making it vulnerable to lost updates.
func (b *AtomicBugsFullBalance) Subtract(amount int64) error {
	current := b.value.Load()
	time.Sleep(100 * time.Microsecond)
	if current-amount < 0 {
		return ErrInsufficientFunds
	}

	b.value.Add(-amount)
	b.trx.Add(1)
	b.updated.Store(time.Now().UnixNano())
	return nil
}
