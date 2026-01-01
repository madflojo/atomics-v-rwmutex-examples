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
	value int64
	// trx counts successful mutations.
	trx int64
	// updated records the timestamp of the last mutation in nanoseconds.
	updated int64
}

// New constructs a zeroed AtomicBugsFullBalance.
func New() (*AtomicBugsFullBalance, error) {
	return &AtomicBugsFullBalance{}, nil
}

// Balance returns the current value.
func (b *AtomicBugsFullBalance) Balance() int64 {
	return atomic.LoadInt64(&b.value)
}

// TransactionCount reports how many writes have completed.
func (b *AtomicBugsFullBalance) TransactionCount() int64 {
	return atomic.LoadInt64(&b.trx)
}

// LastUpdated returns the timestamp for the most recent mutation.
func (b *AtomicBugsFullBalance) LastUpdated() int64 {
	return atomic.LoadInt64(&b.updated)
}

// Add increments the balance and metadata without locking.
func (b *AtomicBugsFullBalance) Add(amount int64) {
	atomic.AddInt64(&b.value, amount)
	atomic.AddInt64(&b.trx, 1)
	atomic.StoreInt64(&b.updated, time.Now().UnixNano())
}

// Subtract decrements the balance but intentionally lacks CAS protection,
// making it vulnerable to lost updates.
func (b *AtomicBugsFullBalance) Subtract(amount int64) error {
	current := atomic.LoadInt64(&b.value)
	time.Sleep(100 * time.Microsecond)
	if current-amount < 0 {
		return ErrInsufficientFunds
	}

	atomic.AddInt64(&b.value, -amount)
	atomic.AddInt64(&b.trx, 1)
	atomic.StoreInt64(&b.updated, time.Now().UnixNano())
	return nil
}
