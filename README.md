# atomics-v-rwmutex-examples ⚖️

Atomic, RWMutex, and Mutex balances—side-by-side Go examples with tests and benchmarks

![Go Version](https://img.shields.io/github/go-mod/go-version/madflojo/atomics-v-rwmutex-examples)
[![Go Report Card](https://goreportcard.com/badge/github.com/madflojo/atomics-v-rwmutex-examples)](https://goreportcard.com/report/github.com/madflojo/atomics-v-rwmutex-examples)

---

## 🧠 What is atomics-v-rwmutex-examples?
A tiny playground showing several Go balance implementations—naïve atomics, CAS‑protected atomics, an RWMutex wrapper, and a standard Mutex wrapper—along with the tests/benchmarks that surfaced their differences for the accompanying blog post over at [madflojo.dev](https://madflojo.dev/posts/are-atomic-operations-faster-than-a-mutex-it-depends/).

- Compare implementation complexity and correctness risks without reading walls of theory.
- Run the test suite (including race-y scenarios) to watch the “atomic bugs” version misbehave.
- Capture real performance numbers with ready-made benchmarks for your own write-ups or talks.

---

Want to see contention fallout? Run `go test ./...` to exercise the same scenarios used in the article, or `go test -bench=. ./...` to capture your own latency numbers.

---

## 🧱 Structure

The project is organized into focused modules so you can depend only on what you need.

| Module/Path | Description | Docs |
| ----------- | ----------- | ---- |
| `balance.go` | Shared `Balance` interface every implementation satisfies. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples#Balance) |
| `implementations/atomics/bugs/simple` | Minimal atomic example with known race bugs; only tracks balance. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/bugs/simple) |
| `implementations/atomics/bugs/full` | Buggy atomic implementation plus transaction/timestamp tracking for apples-to-apples comparisons. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/bugs/full) |
| `implementations/atomics/cas/simple` | CAS-protected counter that only manages the balance value. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/cas/simple) |
| `implementations/atomics/cas/full` | CAS-protected counter with transaction counts and timestamps. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/atomics/cas/full) |
| `implementations/rwmutex/simple` | RWMutex-backed balance guarding just the value. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/rwmutex/simple) |
| `implementations/rwmutex/full` | Feature-complete RWMutex-backed balance mirroring the atomic versions. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/rwmutex/full) |
| `implementations/mutex/simple` | Mutex-backed balance guarding just the value. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/mutex/simple) |
| `implementations/mutex/full` | Feature-complete Mutex-backed balance mirroring the atomic versions. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples/implementations/mutex/full) |
| `balance_test.go` | End-to-end tests covering deposits, withdrawals, insufficient funds, and concurrent subtract races. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples#section-documentation) |
| `balance_benchmark_test.go` | Benchmarks for pure adds, read-before-write adds, and read-only paths to quantify each approach. | [Reference](https://pkg.go.dev/github.com/madflojo/atomics-v-rwmutex-examples#section-directories) |

---

## 🧪 Benchmarks

Benchmarks are designed to exercise different contention patterns:

- Add only: repeatedly calls `Add(1)` under parallel workers.
- Add with read: reads `Balance()` first, then computes an increment using the result to force a read-before-write dependency:
  - If current > 100: increment = current/100 + 1
  - Else: increment = current + 1
- Read only: repeatedly calls `Balance()` under parallel workers.

Implementation detail: the read-only benchmark writes the read value to a shared atomic sink to prevent compiler elision and avoid data races when using `RunParallel`. This can introduce minor contention and slightly skew results; the tradeoff is documented inline in the benchmark source.

Run locally:

- `make benchmarks` (or `go test -run=^$ -bench=. -benchmem ./...`)
- If your environment restricts write access to the global Go cache, set local caches:
  - `GOCACHE=$(pwd)/.gocache GOMODCACHE=$(pwd)/.gomodcache make benchmarks`

## 📦 Tech & Integrations

* Language: Go 1.25.5 (module path `github.com/madflojo/atomics-v-rwmutex-examples`)
* Key deps: `sync/atomic`, `sync.RWMutex`, `sync.Mutex`, and `github.com/madflojo/testlazy/helpers/counter` for test bookkeeping
* Tooling: `go test`, `go test -bench`, and the provided `Makefile` targets if you prefer `make tests`

---

## 📄 License

MIT — see [`LICENSE`](LICENSE).
