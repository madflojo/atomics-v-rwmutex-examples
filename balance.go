package balance

// Balance describes a concurrency-safe account that tracks a value,
// transaction count, and the last update timestamp.
type Balance interface {
	// Balance returns the current account value.
	Balance() int64

	// TransactionCount reports how many mutating operations have been applied.
	TransactionCount() int64

	// LastUpdated returns a monotonic timestamp (nanoseconds) of the latest
	// successful mutation.
	LastUpdated() int64

	// Add increases the account balance by amount. Implementations must treat
	// negative values as undefined behavior.
	Add(amount int64)

	// Subtract decreases the balance by amount or returns an error if the
	// resulting balance would fall below zero.
	Subtract(amount int64) error
}
