---
name: go-testing
description: Design and run Go tests for TCG domain behavior using table-driven cases, boundary values and invariants, and report exact results. Use when writing, validating or auditing tests.
---

# Go Testing

## Design

* Follow the existing style: Arrange / Act / Assert comments, helpers such as `mustCard`, assertions with `errors.Is` on sentinel errors.
* Use table-driven tests when cases share a shape.
* Test observable behavior (public API, returned errors, resulting state), not private implementation.
* Name tests after the behavior, for example `TestNewDeckDefinition_RejectsDuplicateCard`.
* Keep tests deterministic: no real time, no unseeded randomness, no shared mutable state.

## Coverage to consider

* Success cases.
* Validation failures, one invalid input per case.
* Boundaries: zero, minimum, maximum, one below and one above a limit, empty and nil inputs.
* Domain invariants (hand limit 4, deck size 30, max 2 copies, field limits).
* A rejected action leaves state unchanged.
* Immutability and copy semantics where the type promises them.
* Regression of existing behavior near the change.

## Validation role

* Derive expected behavior from the acceptance criteria and the docs (`docs/mvp-story-mode.md`, `docs/game-rules.md`), not from the implementation.
* Do not assume undocumented behavior. Report it as an ambiguity.
* Do not weaken a test to make it pass. Do not modify production code.
* Add or change test files only when the Orchestrator allows it; say which ones.

## Run and report

Run from the repository root and quote the real output:

```text
gofmt -l .
go vet ./...
go test ./...
go test -cover ./...
```

Report: tests executed, passed, failed, missing scenarios, rule violations, regression risks, status (`PASS`, `PASS_WITH_WARNINGS`, `FAIL`). Coverage alone is not evidence of correctness.
