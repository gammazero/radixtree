# Agent Guide for radixtree

## Project Overview

`radixtree` is a Go library implementing an Adaptive Radix Tree (compressed trie / compact prefix tree). Keys are byte sequences; values are generic (`T any`). The module path is `github.com/gammazero/radixtree`.

## Repository Layout

| File | Contents |
|---|---|
| `tree.go` | Core `Tree[T]` type: `New`, `Get`, `Put`, `Delete`, `DeletePrefix`, `Iter`, `IterAt`, `IterPath`, `Inspect`, and internal node helpers |
| `stepper.go` | `Stepper[T]` for byte-at-a-time incremental lookup; `Copy` allows concurrent branching |
| `doc.go` | Package-level doc comment |
| `tree_test.go` | Unit and integration tests for `Tree[T]` |
| `stepper_test.go` | Tests for `Stepper[T]` |
| `doc_test.go` | Example tests (used by `go doc`) |
| `bench_test.go` | Benchmarks |
| `web2`, `web2a` | Word-list fixtures used by benchmarks |
| `go.mod` | Module definition; minimum Go version 1.25 |

## Key Design Invariants

- **Read operations are allocation-free.** `Get`, `Iter`, `IterAt`, `IterPath`, and all `Stepper` methods must not allocate heap memory. Tests and benchmarks verify this.
- **Reads are goroutine-safe; writes are not.** Callers coordinate their own locking. Do not add internal locks.
- **`Get` distinguishes missing key from nil value.** The second return (`bool`) is the only sentinel; never use a zero value as a sentinel.
- **Iteration order is lexical.** Any change to `walk` or `iter` must preserve this.
- **`radices` and `nodes` are parallel slices** kept in sorted byte order. `indexEdge` uses binary search over the dense `radices` byte slice for cache efficiency. Maintain this invariant on every insert, split, and delete.

## Development Workflow

Run tests:

```bash
go test ./...
```

Run tests with the race detector (required before committing any concurrent-access change):

```bash
go test -race ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem ./...
```

Check for allocations in a specific benchmark:

```bash
go test -bench=BenchmarkGet -benchmem -run=^$ .
```

No build step, code generator, or linter configuration is committed. Standard `go vet` is expected to pass.

## Coding Conventions

- Match the existing comment density: exported types, functions, and non-obvious private helpers have doc comments; trivial helpers do not.
- Use plain ASCII punctuation in all source and documentation (no curly quotes, no em dashes).
- Keep changes small and focused. A PR that touches `tree.go` should have a corresponding test in `tree_test.go`.
- New public API needs an example in `doc_test.go`.
- Benchmarks live in `bench_test.go`; add one for any new hot path.

## What Agents Should Watch Out For

- **Split/merge correctness.** `radixNode.split`, `compress`, and `prune` are subtle. When modifying them, re-run the full test suite *and* the race detector.
- **Allocation budget.** After any change to a read path, verify with `-benchmem` that allocs/op for `Get` remains 0.
- **`Stepper` copy semantics.** `Copy` must produce a fully independent snapshot; the original and the copy must be usable concurrently without data races.
- **Parallel-slice consistency.** Every path that adds or removes a child must update both `radices` and `nodes` together and keep them sorted.
