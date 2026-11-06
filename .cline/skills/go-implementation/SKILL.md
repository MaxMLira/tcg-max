---
name: go-implementation
description: Implement an approved TCG task in Go following the project layering, explicit dependencies and error handling, then run gofmt, go vet and tests. Use when writing or changing Go code.
---

# Go Implementation

## Before coding

* Read the approved plan and its acceptance criteria. Implement only what was approved.
* Inspect the existing code, tests and `.cline/rules/` before writing anything. Reuse existing abstractions.
* If a needed rule or decision is missing from the plan and the docs, stop and report it. Do not guess.

## Architecture

* Keep responsibilities separate: domain (rules), engine, application (use cases), repositories and handlers (infrastructure and transport). Dependencies point inward.
* Game rules never live in handlers, repositories or DTOs. DTOs are used at external boundaries; domain entities are not exposed directly.
* Inject dependencies through constructors. No global mutable state or service locators.
* Randomness and time are injected so tests are deterministic.

## Code

* Handle every error explicitly; wrap with `%w` and use the sentinel errors from the domain.
* Prefer the standard library. A new dependency needs explicit approval in the plan.
* Keep changes small and focused. No unrelated refactors.
* Do not commit or push. Commit rules in `.cline/rules/commit-strategy.md` apply only when the user asks for commits.

## Tests

* Add or update tests for every behavior change (happy path, validation failures, boundaries). Follow the `go-testing` skill when it was provided.
* Never delete or weaken a test to make it pass.

## Verify

Run from the repository root and report the exact commands and results:

```text
gofmt -l .
go vet ./...
go test ./...
```

Then inspect `git status` and `git diff` and confirm no unrelated file changed.

## Report

What was implemented, files changed, tests added or changed, commands executed with results, design decisions, remaining issues.
