---

name: implement-feature
description: Implement an approved feature while preserving the project's architecture, domain rules, dependency boundaries, and testing standards.
---------------------------------------------------------------------------------------------------------------------------------------------------

# Implement Feature

## Objective

Implement a previously approved feature or change in the TCG project.

The implementation must follow:

* `docs/game-rules.md`
* `docs/product-requirements.md`
* `docs/architecture.md`
* `.cline/rules/`
* the approved implementation plan

Do not redefine business rules during implementation.

## Before Coding

1. Read the relevant project documentation.
2. Inspect the existing implementation before creating new code.
3. Identify:

   * affected domain modules
   * application/use-case changes
   * repository/infrastructure changes
   * API changes
   * tests that must be added or changed
4. Verify existing abstractions before creating new ones.
5. Check whether the requested behavior already exists partially.

If an important requirement is undefined, stop and report the ambiguity instead of inventing behavior.

## Implementation Principles

### Domain First

Business rules belong to the domain.

Do not place game rules inside:

* HTTP handlers
* WebSocket handlers
* repositories
* database models
* controllers
* transport-specific code

### Dependency Direction

Respect the project's dependency direction:

```text
transport
    ↓
application
    ↓
domain
    ↑
infrastructure
```

The domain must not depend on infrastructure or transport.

### Explicit State

Prefer explicit domain state over implicit behavior.

Important state transitions should be represented clearly and validated.

### Calculate Before Mutate

For important state changes:

1. validate the action
2. calculate the result
3. apply the mutation
4. return the resulting state/event

Avoid partially mutating state before validation completes.

### Composition

Prefer small domain objects with focused responsibilities.

Do not create large classes or services that coordinate unrelated responsibilities.

### Existing Architecture

Do not introduce:

* unnecessary frameworks
* premature abstractions
* generic utility layers without a concrete need
* speculative scalability mechanisms
* duplicated domain logic

Prefer consistency with the existing project.

## Go Standards

Use idiomatic Go.

Prefer:

* small interfaces
* dependency injection through constructors
* explicit error handling
* context propagation where appropriate
* table-driven tests when useful
* pointer receivers only when mutation or identity requires them
* standard library solutions before introducing dependencies

Run formatting and tests after implementation.

## Testing

Every behavioral change must have appropriate tests.

At minimum consider:

* happy path
* invalid input
* domain invariant violations
* important edge cases
* state transitions

Do not change production code merely to make a test easier.

## Scope Control

Implement only what was approved.

If implementation reveals a necessary architectural decision that was not part of the plan:

1. stop
2. describe the issue
3. explain the alternatives
4. request the architectural dec
