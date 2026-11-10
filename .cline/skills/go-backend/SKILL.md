---

name: go-backend
description: Build idiomatic Go backend code for the TCG modular monolith, following domain-driven boundaries, explicit dependencies, and testable design.
----------------------------------------------------------------------------------------------------------------------------------------------------------

# Go Backend

## Objective

Provide implementation guidance for Go code in the TCG backend.

The project is a modular monolith with domain-oriented organization.

Prefer simple, explicit, idiomatic Go over unnecessary abstraction.

## Project Structure

Respect the existing separation between:

* `internal/domain`
* `internal/engine`
* `internal/repository`
* application/use-case code
* transport/API code
* infrastructure

Do not move responsibilities between layers merely for convenience.

## Domain Modeling

Domain objects should represent meaningful game concepts.

Examples include:

* Duel
* Participant
* Card
* CardInstance
* Deck
* Collection
* Zone
* Monster
* Effect
* Chain
* ResponseWindow

Use domain terminology consistently.

Do not introduce alternative names for concepts that already exist in the glossary.

## Interfaces

Define interfaces where they represent meaningful boundaries.

Good examples:

```go
type Repository interface {
    FindByID(ctx context.Context, id string) (...)
}
```

Avoid creating interfaces solely because Go code "should have interfaces".

Prefer defining interfaces close to the consumer when appropriate.

Avoid:

```go
type UserServiceInterface interface {
    ...
}
```

when the abstraction provides no meaningful boundary.

## Constructors and Dependencies

Prefer explicit dependency injection.

Example:

```go
func NewService(repo Repository) *Service {
    return &Service{
        repo: repo,
    }
}
```

Dependencies should be visible.

Avoid:

* global mutable state
* hidden service locators
* package-level dependency containers
* unnecessary singletons

## Context

Use `context.Context` for operations that may:

* access a database
* perform network calls
* interact with infrastructure
* require cancellation or deadlines

Do not store context inside long-lived domain objects.

Do not use context as a generic dependency container.

The domain should not depend on transport-specific context behavior.

## Errors

Errors should be explicit and meaningful.

Prefer domain-specific errors when callers need to distinguish business failures.

Examples:

```go
var ErrInvalidDeck = errors.New("invalid deck")
var ErrInsufficientCards = errors.New("insufficient cards")
```

Do not expose infrastructure errors directly through the public API when they reveal implementation details.

## Concurrency

Concurrency must have a concrete reason.

Before introducing:

* goroutines
* channels
* mutexes
* worker pools
* background workers

identify:

1. what is running concurrently
2. why concurrency is necessary
3. who owns the goroutine
4. how it stops
5. how errors are handled
6. how shared state is protected

Never introduce concurrency merely because Go makes it easy.

## Duel Engine

The Duel Engine is authoritative for game execution.

Transport code must not implement game rules.

For example, WebSocket code should communicate:

```text
player action
    ↓
application
    ↓
Duel Engine
    ↓
new duel state / events
    ↓
WebSocket
```

The engine must remain independent of whether the participant is:

* a human player
* an AI controller

Do not add:

```go
if opponent.IsAI() {
    ...
}
```

inside the Duel Engine.

## State Mutation

Prefer explicit state transitions.

Before mutating important state:

1. validate the action
2. verify the current phase
3. verify the participant
4. verify relevant game conditions
5. calculate the result
6. mutate state

## Persistence

Repositories provide persistence boundaries.

Domain code must not directly depend on:

* PostgreSQL
* SQL drivers
* HTTP
* WebSocket
* ORM-specific APIs

Repository implementations belong outside the domain.

## HTTP

HTTP handlers should be thin.

They should:

1. decode request
2. validate transport-level concerns
3. call application logic
4. map the result to an HTTP response

They should not implement game rules.

## WebSocket

WebSocket is a transport mechanism.

Do not put domain state or game rules inside WebSocket handlers.

WebSocket should translate between:

```text
client message
↔
application command/event
```

## Testing

Prefer testing business behavior at the lowest appropriate level.

For domain logic:

```text
domain → unit tests
```

For application orchestration:

```text
application → unit/integration tests
```

For repositories:

```text
repository → integration tests
```

For HTTP:

```text
handler → HTTP tests
```

Avoid testing implementation details when behavior can be tested directly.

## Quality Rules

Prefer:

* readable code
* small functions
* explicit dependencies
* meaningful names
* standard library
* deterministic behavior
* focused packages

Avoid:

* premature abstractions
* excessive generics
* reflection without a concrete reason
* global state
* unnecessary frameworks
* speculative distributed architecture
