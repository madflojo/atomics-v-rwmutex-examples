package balance

import (
	"sync"
	"sync/atomic"
	"testing"

	atomicbugsfull "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/bugs/full"
	atomicbugssimple "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/bugs/simple"
	atomiccasfull "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/cas/full"
	atomiccassimple "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/cas/simple"
	mutexfull "github.com/madflojo/atomics-v-rwmutex-examples/implementations/mutex/full"
	mutexsimple "github.com/madflojo/atomics-v-rwmutex-examples/implementations/mutex/simple"
	rwmutexfull "github.com/madflojo/atomics-v-rwmutex-examples/implementations/rwmutex/full"
	rwmutexsimple "github.com/madflojo/atomics-v-rwmutex-examples/implementations/rwmutex/simple"

	"github.com/madflojo/testlazy/helpers/counter"
)

func TestBalanceImplementations(t *testing.T) {
	testCases := []struct {
		name      string
		balance   Balance
		hasMeta   bool
		expectBug bool
	}{
		{
			name:      "Atomic Balance (bugs/simple)",
			balance:   atomicbugssimple.New(),
			hasMeta:   false,
			expectBug: true,
		},
		{
			name:      "Atomic Balance (bugs/full)",
			balance:   atomicbugsfull.New(),
			hasMeta:   true,
			expectBug: true,
		},
		{
			name:      "Atomic Balance (CAS/simple)",
			balance:   atomiccassimple.New(),
			hasMeta:   false,
			expectBug: false,
		},
		{
			name:      "Atomic Balance (CAS/full)",
			balance:   atomiccasfull.New(),
			hasMeta:   true,
			expectBug: false,
		},
		{
			name:      "RWMutex Balance (simple)",
			balance:   rwmutexsimple.New(),
			hasMeta:   false,
			expectBug: false,
		},
		{
			name:      "RWMutex Balance (full)",
			balance:   rwmutexfull.New(),
			hasMeta:   true,
			expectBug: false,
		},
		{
			name:      "Mutex Balance (simple)",
			balance:   mutexsimple.New(),
			hasMeta:   false,
			expectBug: false,
		},
		{name: "Mutex Balance (full)", balance: mutexfull.New(), hasMeta: true, expectBug: false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Helper()

			acct := tc.balance
			bal := counter.New()
			trx := counter.New()

			t.Run("initial state", func(t *testing.T) {
				if acct.Balance() != bal.Value() {
					t.Fatalf("expected zero balance, got %d", acct.Balance())
				}

				if tc.hasMeta {
					if acct.TransactionCount() != trx.Value() {
						t.Fatalf("expected zero transaction count, got %d", acct.TransactionCount())
					}
					if acct.LastUpdated() != 0 {
						t.Fatalf("expected zero last updated, got %d", acct.LastUpdated())
					}
				} else {
					if acct.TransactionCount() != 0 {
						t.Fatalf("simple balances must not track transactions, got %d", acct.TransactionCount())
					}
					if acct.LastUpdated() != 0 {
						t.Fatalf("simple balances must not track last updated, got %d", acct.LastUpdated())
					}
				}
			})

			t.Run("single threaded", func(t *testing.T) {
				t.Run("deposit", func(t *testing.T) {
					prev := acct.LastUpdated()
					acct.Add(1000)
					bal.Add(1000)
					trx.Add(1)

					if acct.Balance() != bal.Value() {
						t.Fatalf("deposit mismatch, got %d want %d", acct.Balance(), bal.Value())
					}

					if tc.hasMeta {
						if acct.LastUpdated() < prev {
							t.Fatalf("last updated failed to advance after deposit")
						}
					} else if acct.LastUpdated() != 0 {
						t.Fatalf("simple balances must not track last updated")
					}

					prev = acct.LastUpdated()
					for i := 0; i < 1000; i++ {
						acct.Add(100)
						bal.Add(100)
						trx.Add(1)
					}

					if acct.Balance() != bal.Value() {
						t.Fatalf(
							"balance mismatch after deposits, got %d want %d",
							acct.Balance(),
							bal.Value(),
						)
					}

					if tc.hasMeta {
						if acct.LastUpdated() < prev {
							t.Fatalf("last updated failed to advance after deposits")
						}
					} else if acct.LastUpdated() != 0 {
						t.Fatalf("simple balances must not track last updated")
					}
				})

				t.Run("withdraw", func(t *testing.T) {
					prev := acct.LastUpdated()
					for i := 0; i < 500; i++ {
						if err := acct.Subtract(50); err != nil {
							t.Fatalf("unexpected subtract error: %v", err)
						}
						bal.Subtract(50)
						trx.Add(1)
					}

					if acct.Balance() != bal.Value() {
						t.Fatalf(
							"balance mismatch after withdrawals, got %d want %d",
							acct.Balance(),
							bal.Value(),
						)
					}

					if tc.hasMeta {
						if acct.LastUpdated() < prev {
							t.Fatalf("last updated failed to advance after withdrawals")
						}
					} else if acct.LastUpdated() != 0 {
						t.Fatalf("simple balances must not track last updated")
					}
				})
			})

			t.Run("concurrent subtract", func(t *testing.T) {
				const (
					deposit  = 1_000
					withdraw = 25
					workers  = 32
					iters    = 80
				)

				if current := acct.Balance(); current > deposit {
					if err := acct.Subtract(current - deposit); err != nil {
						t.Fatalf("failed to normalize balance before concurrent test: %v", err)
					}
					bal.Subtract(current - deposit)
					trx.Add(1)
				}

				prev := acct.LastUpdated()
				acct.Add(deposit)
				bal.Add(deposit)
				trx.Add(1)

				if tc.hasMeta {
					if acct.LastUpdated() < prev {
						t.Fatalf("last updated failed to advance after prep deposit")
					}
				} else if acct.LastUpdated() != 0 {
					t.Fatalf("simple balances must not track last updated")
				}
				prev = acct.LastUpdated()

				var wg sync.WaitGroup
				var success int64
				var fail int64
				total := workers * iters

				for w := 0; w < workers; w++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						for i := 0; i < iters; i++ {
							if err := acct.Subtract(withdraw); err != nil {
								atomic.AddInt64(&fail, 1)
								continue
							}

							bal.Subtract(withdraw)
							trx.Add(1)
							atomic.AddInt64(&success, 1)
						}
					}()
				}

				wg.Wait()

				successCount := atomic.LoadInt64(&success)
				failCount := atomic.LoadInt64(&fail)
				if got := successCount + failCount; got != int64(total) {
					t.Fatalf("unexpected subtract attempts, got %d want %d", got, total)
				}

				if tc.hasMeta {
					if acct.LastUpdated() < prev {
						t.Fatalf("last updated failed to advance after concurrent subtracts")
					}
				} else if acct.LastUpdated() != 0 {
					t.Fatalf("simple balances must not track last updated")
				}

				if tc.expectBug {
					if acct.Balance() >= 0 {
						t.Fatalf(
							"expected buggy implementation to go negative, got %d",
							acct.Balance(),
						)
					}
					if failCount == 0 {
						t.Fatalf("expected at least one subtract failure")
					}
					return
				}

				if acct.Balance() < 0 {
					t.Fatalf("balance went negative, got %d", acct.Balance())
				}

				if acct.Balance() != bal.Value() {
					t.Fatalf(
						"balance mismatch after concurrent subtracts, got %d want %d",
						acct.Balance(),
						bal.Value(),
					)
				}

				if failCount == 0 {
					t.Fatalf("expected at least one subtract failure")
				}
			})

			t.Run("insufficient funds", func(t *testing.T) {
				prevBal := bal.Value()
				prevTrx := trx.Value()
				prevUpdate := acct.LastUpdated()
				amount := prevBal + 123

				if err := acct.Subtract(amount); err == nil {
					t.Fatalf("expected error when subtracting past balance")
				}

				if acct.Balance() != prevBal {
					t.Fatalf(
						"balance changed after failed subtract, got %d want %d",
						acct.Balance(),
						prevBal,
					)
				}

				if tc.hasMeta {
					if acct.TransactionCount() != prevTrx {
						t.Fatalf(
							"transaction count changed after failed subtract, got %d want %d",
							acct.TransactionCount(),
							prevTrx,
						)
					}
					if acct.LastUpdated() != prevUpdate {
						t.Fatalf("last updated changed after failed subtract")
					}
				} else {
					if acct.TransactionCount() != 0 {
						t.Fatalf("simple balances must not track transactions, got %d", acct.TransactionCount())
					}
					if acct.LastUpdated() != 0 {
						t.Fatalf("simple balances must not track last updated")
					}
				}
			})

			if tc.hasMeta {
				t.Run("transaction counter", func(t *testing.T) {
					if acct.TransactionCount() != trx.Value() {
						t.Fatalf(
							"transaction counter mismatch, got %d want %d",
							acct.TransactionCount(),
							trx.Value(),
						)
					}
				})
			}
		})
	}
}
