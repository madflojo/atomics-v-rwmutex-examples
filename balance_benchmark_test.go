package balance

import (
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
)

var benchmarkImplementations = []struct {
	name string
}{
	{
		name: "Atomic_Balance_bugs_simple",
	},
	{
		name: "Atomic_Balance_bugs_full",
	},
	{
		name: "Atomic_Balance_CAS_simple",
	},
	{
		name: "Atomic_Balance_CAS_full",
	},
	{
		name: "RWMutex_Balance_simple",
	},
	{
		name: "RWMutex_Balance_full",
	},
	{
		name: "Mutex_Balance_simple",
	},
	{
		name: "Mutex_Balance_full",
	},
}

// balanceSink ensures Balance() results are observed in read-only benchmarks.
// Note: writing to a shared sink can add contention and slightly skew results,
// but it prevents compiler elision and avoids data races in RunParallel.
var balanceSink int64

func BenchmarkBalanceAdd(b *testing.B) {
	for _, impl := range benchmarkImplementations {
		impl := impl
		b.Run(impl.name, func(b *testing.B) {
			account := newBalanceForName(impl.name)

			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					account.Add(1)
				}
			})
		})
	}
}

func BenchmarkBalanceAddWithRead(b *testing.B) {
	for _, impl := range benchmarkImplementations {
		impl := impl
		b.Run(impl.name, func(b *testing.B) {
			account := newBalanceForName(impl.name)

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
		})
	}
}

func BenchmarkBalanceReadOnly(b *testing.B) {
	for _, impl := range benchmarkImplementations {
		impl := impl
		b.Run(impl.name, func(b *testing.B) {
			account := newBalanceForName(impl.name)

			// Prime the value to avoid zero-edge quirks.
			account.Add(1)

			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					// Write to a shared sink to keep the read observable and avoid races.
					atomic.StoreInt64(&balanceSink, account.Balance())
				}
			})
		})
	}
}

// newBalanceForName constructs a fresh Balance for the given benchmark name.
// Keeping construction here avoids repeating switch logic and ensures each
// sub-benchmark gets a clean instance.
func newBalanceForName(name string) Balance {
	switch name {
	case "Atomic_Balance_bugs_simple":
		return atomicbugssimple.New()
	case "Atomic_Balance_bugs_full":
		return atomicbugsfull.New()
	case "Atomic_Balance_CAS_simple":
		return atomiccassimple.New()
	case "Atomic_Balance_CAS_full":
		return atomiccasfull.New()
	case "RWMutex_Balance_simple":
		return rwmutexsimple.New()
	case "RWMutex_Balance_full":
		return rwmutexfull.New()
	case "Mutex_Balance_simple":
		return mutexsimple.New()
	case "Mutex_Balance_full":
		return mutexfull.New()
	default:
		return rwmutexsimple.New()
	}
}
