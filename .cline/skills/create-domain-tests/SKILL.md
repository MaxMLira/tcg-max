---

name: create-domain-tests
description: Create focused Go tests for TCG domain rules, state transitions, invariants, combat, effects, and duel behavior.
-----------------------------------------------------------------------------------------------------------------------------

# Create Domain Tests

## Objective

Create tests that validate the behavior of the TCG domain against the official game rules.

The primary source of truth is:

```text
docs/game-rules.md
```

Product behavior comes from:

```text
docs/product-requirements.md
```

Architecture constraints come from:

```text
docs/architecture.md
```

## Testing Philosophy

Test behavior and invariants, not implementation details.

A good domain test should answer:

> "Does the game behave correctly according to the rules?"

rather than:

> "Did this private method execute?"

## Before Writing Tests

1. Read the relevant game rules.
2. Identify the rule being tested.
3. Identify valid states.
4. Identify invalid states.
5. Identify important boundary conditions.
6. Inspect existing tests to avoid duplication.
7. Determine whether the behavior belongs to:

   * domain
   * Duel Engine
   * application
   * repository
   * transport

## Domain Test Structure

Prefer Arrange → Act → Assert.

Example:

```go
func TestCombatResolver_AttackVsAttack(t *testing.T) {
    // Arrange

    // Act

    // Assert
}
```

Use table-driven tests when several cases share the same structure.

## Game Rules to Prioritize

When relevant, explicitly test:

### Deck

* exactly 30 cards
* maximum 2 copies of the same card
* valid/invalid deck composition

### Hand

* maximum 4 cards
* draw fills hand up to 4
* full hand prevents drawing
* deck exhaustion causes immediate defeat when a draw is required

### Summoning

* weight 1–5 requires zero tributes
* weight 6–8 requires one tribute
* weight 9–12 requires two tributes
* weight 13+ requires three tributes
* tribute monsters must exist on the field
* newly placed monster cannot be used as tribute during the same Phase 1

### Phase 1

* exactly one placement action
* summon OR set
* cannot perform two placement actions
* position changes do not count as placement

### Combat

Test all relevant combinations:

* ATK vs ATK
* ATK vs DEF
* direct attack
* equal ATK
* attacker stronger
* defender stronger
* piercing
* Double Attack

Include boundary values.

### Face-Down Monsters

* information is hidden from opponent
* owner can identify the card
* attacked face-down monster is revealed before resolution
* surviving monster remains unable to attack during opponent's turn
* face-down monster can be tributed

### Effects

Test:

* valid activation
* invalid activation
* automatic monster effects
* optional Trap activation
* activation conditions
* costs
* duration
* modifiers
* effect resolution

### Chain

Test:

* activation order
* LIFO resolution
* response windows
* Trap responding to Magic
* no new activation during resolution

### Graveyard

Test:

* destroyed cards go to Graveyard
* tributed cards go to Graveyard
* discarded cards go to Graveyard
* used Magic/Trap goes to Graveyard
* Field Magic replacement
* Graveyard visibility/order where relevant

### Revival

Test:

* only allowed cards can be revived
* cannot revive opponent's cards
* revived monster chooses position
* revived monster cannot attack that turn
* revived monster cannot be tributed that turn
* destruction after revival returns it to Graveyard

## Invariants

Whenever possible, test invariants such as:

```text
Deck + Hand + Field + Graveyard
```

represent the cards currently belonging to the participant during the duel.

A card should not unexpectedly exist in two zones simultaneously.

Other useful invariants include:

* LP cannot become invalid without a rule allowing it
* field slot limits are respected
* hand limit is respected
* Phase 1 allows only one placement
* a finished duel cannot continue processing actions
* invalid actions do not partially mutate state

## State Transition Tests

For important actions, test both:

```text
valid state → expected new state
```

and:

```text
invalid state → error + unchanged state
```

The second case is especially important.

If an action fails validation, verify that the domain was not partially mutated.

## Edge Cases

Always consider:

* zero values
* maximum values
* boundary values
* empty collections
* full collections
* last card in deck
* equal ATK/DEF
* LP reaching exactly zero
* LP becoming negative
* empty field
* full field
* multiple eligible targets
* repeated attacks
* simultaneous destruction
* invalid phase
* invalid participant
* already finished duel

## Test Naming

Names should describe business behavior.

Prefer:

```go
TestDuel_CannotPlaceTwoCardsInPhaseOne
```

over:

```go
TestDuel_PhaseOne2
```

Prefer:

```go
TestCombatResolver_PiercingDealsDifferenceWhenAttackExceedsDefense
```

over:

```go
TestPiercing
```

## Test Independence

Tests must be deterministic.

Avoid:

* real time dependencies
* random behavior without controlled seeds
* external services
* shared mutable state
* test ordering dependencies

## Completion Checklist

Before finishing:

* [ ] Relevant game rule was read
* [ ] Happy path tested
* [ ] Invalid path tested
* [ ] Boundary cases considered
* [ ] State mutation verified
* [ ] Important invariants verified
* [ ] Tests are deterministic
* [ ] `go test ./...` passes

## Output

Report:

1. Rules covered
2. Tests created/updated
3. Edge cases covered
4. Invariants validated
5. Test command executed
6. Any uncovered behavior
