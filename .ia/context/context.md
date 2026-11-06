# Project Context

## Project

This project is a competitive 1v1 Trading Card Game (TCG).

The project is also an engineering laboratory for developing software engineering and AI engineering capabilities.

The system should prioritize maintainability, clear domain boundaries, testability, deterministic behavior, and controlled AI-assisted development.

---

## Primary Goals

The project has two main goals:

1. Build a maintainable and extensible TCG.
2. Use the project as an AI Engineering laboratory.

The AI Engineering side includes:

* context engineering;
* AI agents;
* skills;
* MCP;
* RAG;
* guardrails;
* evaluation;
* human-in-the-loop workflows.

---

## Technology

### Backend

* Go

### Frontend

* React
* TypeScript

### Database

* PostgreSQL

### Communication

* REST for the general API
* WebSocket for real-time Multiplayer

### Deployment

* Docker
* Docker Compose

---

## Architecture

The backend uses a Modular Monolith architecture.

The core game is implemented as a domain-oriented Duel Engine.

The domain must remain independent from:

* HTTP;
* REST;
* WebSocket;
* PostgreSQL;
* JWT;
* React;
* external services.

The architecture is documented in:

`docs/architecture.md`

---

## Product Modes

The MVP contains two game modes.

### Story Mode

Player versus AI-controlled Character.

Example:

```text
Player
  vs
Richard
```

Story Characters are controlled through AI Controllers.

### Multiplayer

Player versus Player.

The MVP uses queue-based matchmaking:

```text
Player 1
    ↓
Queue
    ↓
Player 2
    ↓
Match
    ↓
Duel
```

Private rooms and invitations are outside the MVP.

---

## Duel Engine

Story Mode and Multiplayer use the same Duel Engine.

The Duel Engine is responsible for executing the Game Rules.

The Duel Engine must not contain separate rules for:

* AI opponents;
* human opponents.

The difference is in how actions are generated.

```text
Story
Character
   ↓
AI Controller
   ↓
Duel Engine

Multiplayer
Player
   ↓
Player Controller
   ↓
Duel Engine
```

---

## Core Domain Concepts

Important concepts include:

```text
User
Collection
Deck
Card
CardInstance

Duel
Participant
Turn
Phase

Deck
Hand
Field
Graveyard

Combat
Effect
Chain
ResponseWindow
```

---

## Card Model

The project distinguishes between card definition, ownership, deck composition, and runtime state.

```text
Card
  ↓
Collection
  ↓
Deck
  ↓
CardInstance
  ↓
Duel
```

### Card

Defines what the card is.

### Collection

Represents cards owned by a User.

### Deck

Represents cards selected by the User for a deck.

### CardInstance

Represents a concrete card copy participating in a Duel.

These concepts must not be unnecessarily merged.

---

## Source of Truth

Different documents have different responsibilities.

### Game Rules

`docs/game-rules.md`

Defines how the game works.

### Product Requirements

`docs/product-requirements.md`

Defines what the product must provide.

### Architecture

`docs/architecture.md`

Defines how the software is structured.

### Glossary

`docs/glossary.md`

Defines project terminology.

### ADRs

`docs/adr/`

Record important architectural decisions and their rationale.

---

## Current MVP

The MVP includes:

* User registration;
* authentication;
* starter deck;
* Collection;
* multiple Decks;
* Story Mode;
* AI-controlled Characters;
* Duel Engine;
* Multiplayer;
* matchmaking queue;
* WebSocket Multiplayer communication.

The MVP does not include:

* Card Shop;
* Stars;
* betting;
* real-money purchases;
* ranking;
* seasons;
* divisions;
* tournaments;
* private rooms;
* invitations;
* social chat;
* spectator mode;
* replay system;
* Duel recovery after server failure.

---

## Development Philosophy

The project should evolve incrementally.

Do not introduce complexity merely because it may be useful in the future.

Prefer:

* simple solutions;
* explicit domain behavior;
* small modules;
* strong tests;
* clear boundaries;
* deterministic tooling.

Architectural complexity must be justified by a real requirement.

---

## AI Development Workflow

AI-assisted development follows:

```text
Requirement
    ↓
Planner
    ↓
Architect
    ↓
Human Approval
    ↓
Developer
    ↓
Tester
    ↓
Reviewer
    ↓
Human Approval
```

Agents assist development but do not own Game Rules or architectural decisions.

---

## Important Constraint

If a Game Rule or requirement is undefined, the AI must not invent behavior.

It must identify the ambiguity and request clarification.
