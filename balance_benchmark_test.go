package balance

import (
	"sync/atomic"
	"testing"

	atomicbugsfull "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/bugs/full"
	atomicbugssimple "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/bugs/simple"
	atomiccasfull "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/cas/full"
	atomiccassimple "github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/cas/simple"
	rwmutexfull "github.com/madflojo/atomics-v-rwmutex-examples/implementations/rwmutex/full"
	rwmutexsimple "github.com/madflojo/atomics-v-rwmutex-examples/implementations/rwmutex/simple"
)

var benchmarkImplementations = []struct {
	name    string
	account func(tb testing.TB) Balance
}{
	{
		name: "Atomic_Balance_bugs_simple",
		account: func(tb testing.TB) Balance {
			tb.Helper()

			b, err := atomicbugssimple.New()
			if err != nil {
				tb.Fatalf("failed to create atomic balance: %v", err)
			}

			return b
		},
	},
	{
		name: "Atomic_Balance_bugs_full",
		account: func(tb testing.TB) Balance {
			tb.Helper()

			b, err := atomicbugsfull.New()
			if err != nil {
				tb.Fatalf("failed to create atomic balance: %v", err)
			}

			return b
		},
	},
	{
		name: "Atomic_Balance_CAS_simple",
		account: func(tb testing.TB) Balance {
			tb.Helper()

			b, err := atomiccassimple.New()
			if err != nil {
				tb.Fatalf("failed to create atomic balance: %v", err)
			}

			return b
		},
	},
	{
		name: "Atomic_Balance_CAS_full",
		account: func(tb testing.TB) Balance {
			tb.Helper()

			b, err := atomiccasfull.New()
			if err != nil {
				tb.Fatalf("failed to create atomic balance: %v", err)
			}

			return b
		},
	},
	{
		name: "RWMutex_Balance_simple",
		account: func(tb testing.TB) Balance {
			tb.Helper()

			b, err := rwmutexsimple.New()
			if err != nil {
				tb.Fatalf("failed to create mutex balance: %v", err)
			}

			return b
		},
	},
	{
		name: "RWMutex_Balance_full",
		account: func(tb testing.TB) Balance {
			tb.Helper()

			b, err := rwmutexfull.New()
			if err != nil {
				tb.Fatalf("failed to create mutex balance: %v", err)
			}

			return b
		},
	},
}

var balanceSink int64

func BenchmarkBalanceAdd(b *testing.B) {
	for _, impl := range benchmarkImplementations {
		impl := impl
		b.Run(impl.name, func(b *testing.B) {
			benchmarkAdd(b, impl.account)
		})
	}
}

func BenchmarkBalanceAddWithRead(b *testing.B) {
	for _, impl := range benchmarkImplementations {
		impl := impl
		b.Run(impl.name, func(b *testing.B) {
			benchmarkAddWithRead(b, impl.account)
		})
	}
}

func BenchmarkBalanceReadOnly(b *testing.B) {
	for _, impl := range benchmarkImplementations {
		impl := impl
		b.Run(impl.name, func(b *testing.B) {
			benchmarkReadOnly(b, impl.account)
		})
	}
}

func benchmarkAdd(b *testing.B, accountFactory func(tb testing.TB) Balance) {
	b.Helper()

	account := accountFactory(b)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			account.Add(1)
		}
	})
}

func benchmarkAddWithRead(b *testing.B, accountFactory func(tb testing.TB) Balance) {
    b.Helper()

    account := accountFactory(b)
    b.ReportAllocs()
    b.ResetTimer()

    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Read-before-write: compute increment from current balance
            cur := account.Balance()
            var inc int64
            if cur > 100 {
                inc = (cur / 100) + 1
            } else {
                inc = cur + 1
            }
            account.Add(inc)
        }
    })
}

func benchmarkReadOnly(b *testing.B, accountFactory func(tb testing.TB) Balance) {
	b.Helper()

	account := accountFactory(b)
	account.Add(1)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.StoreInt64(&balanceSink, account.Balance())
		}
	})
}
