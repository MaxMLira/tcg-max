# Domain Rules

## 1. Game Rules

`docs/game-rules.md` is the authority for game behavior.

Never implement behavior that contradicts it.

---

## 2. Domain Independence

The domain must not depend on:

* HTTP;
* REST;
* WebSocket;
* PostgreSQL;
* JWT;
* React;
* ORM;
* external APIs.

---

## 3. Backend Authority

The backend is the authority for game validity.

Never rely on the frontend to enforce Game Rules.

---

## 4. State

Prefer explicit domain state over boolean flags that duplicate existing facts.

Represent domain facts directly whenever possible.

---

## 5. Domain Invariants

Domain objects are responsible for protecting their own invariants.

Examples:

* Field limits;
* Hand limits;
* Deck constraints;
* valid card locations;
* valid actions;
* valid phase transitions.

---

## 6. Calculate Before Mutating

Specialized domain services should calculate results before the Duel applies them.

Example:

```text
Attack
  ↓
CombatResolver
  ↓
CombatResult
  ↓
Duel applies result
```

---

## 7. Card Definitions

`Card` represents the definition of a card.

Do not mutate the base Card definition because of a runtime Duel effect.

Runtime state belongs to `CardInstance` and other Duel state.

---

## 8. Collection vs Duel State

Do not mix persistent ownership with runtime Duel state.

```text
Collection
    ↓
persistent ownership

CardInstance
    ↓
runtime Duel state
```

---

## 9. Effects

Do not spread effect-specific conditionals throughout the domain.

Prefer explicit effect behavior and composition.

---

## 10. Duel Engine

Story Mode and Multiplayer must use the same Duel Engine.

The Duel Engine must not contain rules based on whether a Participant is:

* human;
* AI.

---

## 11. Chain

Chain resolution follows the Game Rules.

The system must not invent additional priority or response mechanics.

---

## 12. Transport Independence

The Duel Engine must not depend on:

* REST;
* WebSocket;
* timers;
* network connections.

Transport and timing infrastructure belong outside the domain.
