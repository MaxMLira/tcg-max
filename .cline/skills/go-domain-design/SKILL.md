---
name: go-domain-design
description: Model TCG domain concepts in Go with explicit invariants, validated constructors, immutability of definitions and no infrastructure dependencies. Use when designing or changing domain types.
---

# Go Domain Design

## Principles

* Inspect `internal/domain` first and follow its conventions: private fields with accessors, constructors that validate and return an error, sentinel errors in `errors.go` wrapped with `%w`, doc comments that cite the rule they implement.
* Make invalid states hard to build. Validate in the constructor; do not rely on callers remembering to call `Validate()`.
* Prefer explicit state over boolean flags that duplicate facts.
* Keep the domain free of HTTP, database, JWT, WebSocket and other infrastructure imports.
* `Card` is a definition and is never mutated by a duel. Runtime state belongs to `CardInstance` and duel state. Collection ownership is separate from runtime state.
* Calculate before mutating: validate and compute the result, then apply it. A rejected action must leave state unchanged.
* Keep the engine independent of whether a participant is a human or an AI.
* Where the MVP requires extensibility (for example Monster subtypes), use a validated value type rather than a closed list hard-wired into the model.

## Restraint

* Do not add interfaces, generics or layers without a concrete need in the current task.
* Reuse existing types before creating new ones.
* Add no dependency without explaining why, and without approval.

## Compatibility

* Preserve exported names, signatures and behavior unless the approved task explicitly requires a breaking change.
* When a breaking change is required, list every affected caller and test in the design.

## Output of a design

Types and signatures, invariants each type protects, errors added, files touched, what is deliberately not included.
