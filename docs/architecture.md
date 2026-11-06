# Architecture

## 1. Overview

This project is a competitive 1v1 Trading Card Game (TCG) with its own domain rules and game mechanics.

The system has two complementary goals:

1. Build a maintainable and extensible game system.
2. Use the project as an AI Engineering laboratory to apply context engineering, agents, skills, tools, MCP, RAG, guardrails, evaluation, and human-in-the-loop workflows.

The architecture prioritizes:

* clear domain boundaries;
* explicit business rules;
* testability;
* low coupling;
* extensibility;
* deterministic game behavior;
* controlled AI-assisted development.

The game engine is the source of truth for game rules.

AI agents must not invent or implicitly modify game rules.

---

# 2. Architectural Style

## 2.1 Modular Monolith

The backend initially uses a Modular Monolith architecture.

All domain modules run inside a single deployable Go application while maintaining clear internal boundaries.

The architecture intentionally avoids microservices because:

* the domain is still evolving;
* the project does not currently require distributed processing;
* a single process reduces operational complexity;
* domain boundaries should be discovered before introducing network boundaries;
* most game operations are naturally coordinated around a single Duel.

Future extraction of modules into independent services is possible if a real requirement appears.

Microservices must not be introduced solely because the project may eventually scale.

---

## 2.2 Domain-Driven Organization

The backend is organized around business and domain concepts rather than technical layers alone.

Important concepts include:

* User;
* Collection;
* Deck;
* Card;
* CardInstance;
* Duel;
* Participant;
* Character;
* Turn;
* Phase;
* Chain;
* Effect;
* Combat;
* Deck;
* Hand;
* Field;
* Graveyard.

The domain must remain independent from:

* HTTP;
* REST;
* WebSocket;
* PostgreSQL;
* React;
* JWT;
* ORM;
* external infrastructure.

---

# 3. System Architecture

The initial system consists of:

```text
                         ┌──────────────────────┐
                         │   React + TypeScript │
                         └──────────┬───────────┘
                                    │
                         ┌──────────┴───────────┐
                         │                      │
                        REST                WebSocket
                         │                      │
                         └──────────┬───────────┘
                                    │
                         ┌──────────▼───────────┐
                         │       Go Backend     │
                         │                      │
                         │ API / Application    │
                         │ Domain / Infra       │
                         └──────────┬───────────┘
                                    │
                         ┌──────────▼───────────┐
                         │      PostgreSQL      │
                         └──────────────────────┘
```

REST is the initial general-purpose communication mechanism.

WebSocket is used for real-time Multiplayer communication.

The Duel Engine must remain independent from both communication mechanisms.

---

# 4. Backend Architecture

The backend uses Go and follows the Modular Monolith architecture.

Initial conceptual structure:

```text
cmd/
└── server/

internal/
├── auth/
├── user/
├── collection/
├── deck/
├── card/
├── story/
├── matchmaking/
├── duel/
├── participant/
├── zone/
├── combat/
├── effect/
└── ...
```

Supporting infrastructure may include:

```text
internal/
├── http/
├── websocket/
├── repository/
├── database/
└── ...
```

The exact package structure may evolve as implementation progresses.

The architecture should favor domain boundaries over prematurely rigid package structures.

---

# 5. Core Domain Modules

## 5.1 User

The User represents the player account.

For the current architecture:

```text
User = Player
```

The same concept represents both the authenticated account and the player participating in the game.

A User may own:

* a Collection;
* multiple Decks;
* Story progress;
* future persistent game information.

Authentication concerns such as JWT validation belong outside the domain.

---

## 5.2 Collection

The Collection represents the cards owned by a User.

Conceptually:

```text
User
 │
 └── Collection
       ├── Card A
       ├── Card A
       ├── Card B
       └── Card C
```

The Collection represents persistent ownership.

It is distinct from the runtime state of cards participating in a Duel.

---

## 5.3 Deck

A User may own multiple Decks.

Conceptually:

```text
User
 │
 ├── Deck A
 │     └── Cards
 │
 ├── Deck B
 │     └── Cards
 │
 └── Deck C
       └── Cards
```

Deck rules are defined by the Game Rules document.

The Deck domain object is responsible for protecting deck-related invariants.

A Deck is a player's selected composition of cards, while the Collection represents all cards owned by the player.

---

## 5.4 Card

`Card` represents the definition/template of a card.

Conceptually:

```text
Card
├── ID
├── Name
├── Type
├── Attack
├── Defense
├── Weight
├── MonsterSubtype
├── AcceptedMonsterSubtype
└── Effect
```

A Card represents what the card is, not its runtime state during a Duel.

### MonsterSubtype

`MonsterSubtype` is the subtype of a Monster (for example Dragon, Warrior, Mage, Aquatic).

It is an extensible validated value type, not a closed enumeration:

* it has a private name and a validating constructor;
* the name is trimmed of surrounding whitespace and must not be blank;
* subtypes are compared case-insensitively through `Equals`;
* the zero value means "no subtype".

Every Monster has a subtype. Magic and Trap cards have none.

New subtypes can be added without changing the domain model.

### AcceptedMonsterSubtype

`AcceptedMonsterSubtype` is the Monster subtype an Equipment card is compatible with.

It is a separate fact from the card's own `MonsterSubtype`: an Equipment is a Magic card, so its own subtype is always empty, and it declares the subtype of the Monsters it can equip.

A Magic card without an `AcceptedMonsterSubtype` is valid.

`Card.AcceptsMonsterSubtype` only answers whether the definition is compatible with a subtype. Placing, activating and linking Equipment, and the ATK bonus, are Duel behavior and do not belong to the `Card` definition.

---

## 5.5 CardInstance

`CardInstance` represents a concrete copy of a Card participating in a Duel.

Conceptually:

```text
Card
  │
  ├── CardInstance #101 → FIELD
  ├── CardInstance #102 → HAND
  └── CardInstance #103 → GRAVEYARD
```

A CardInstance contains runtime state such as:

```text
CardInstance
├── ID
├── CardID
├── OwnerID
├── Location
├── Position
└── AttacksPerformed
```

CardInstance is not the representation of Collection ownership.

The separation is:

```text
Card
   ↓
card definition

Collection
   ↓
cards owned by User

Deck
   ↓
cards selected for a deck

CardInstance
   ↓
concrete card copy participating in a Duel
```

This separation prevents persistent ownership concepts from being mixed with runtime Duel state.

---

# 6. Duel Engine

The Duel Engine is the core of the game.

It contains the rules and state required to execute a Duel.

The same Duel Engine is used by:

```text
Story Mode
Player vs AI-controlled Character

Multiplayer
Player vs Player
```

The Duel Engine must not contain separate game rules for AI and human opponents.

The difference between Story and Multiplayer belongs outside the Duel rules.

---

# 7. Duel

The `Duel` represents an active game between two Participants.

Conceptually:

```text
Duel
├── ID
├── Participants
├── CurrentTurn
├── Phase
├── Chain
├── ResponseWindow
└── Status
```

The Duel coordinates the game but should not become a God Object.

Complex calculations and specialized behavior should be delegated to domain services and domain objects.

---

## 7.1 Participants

A Duel contains two Participants.

A Participant represents an entity capable of taking actions in a Duel.

The Duel does not need to know whether the Participant is controlled by a human or an AI.

Conceptually:

```text
Duel
├── Participant A
└── Participant B
```

Each Participant is associated with the appropriate controller outside the core game rules.

---

# 8. Player and AI Controllers

The system distinguishes the identity of a Participant from the mechanism that generates its actions.

Conceptually:

```text
Participant
     │
     ├── PlayerController
     │
     └── AIController
```

For Multiplayer:

```text
User
 ↓
PlayerController
 ↓
Duel
```

For Story Mode:

```text
Character
 ↓
AIController
 ↓
Duel
```

The Duel Engine only receives valid game actions.

It must not contain logic such as:

```text
if opponent is AI
```

inside its core game rules.

This allows Story Mode and Multiplayer to reuse the same Duel Engine.

---

# 9. Story Mode

Story Mode introduces named Characters controlled by AI.

Example:

```text
Player
   vs
Richard
```

Richard is a Character participating in the Duel, while an AI Controller determines his actions.

The Character may contain Story-specific information such as:

* name;
* deck;
* difficulty;
* future progression data;
* AI strategy configuration.

Story-specific information must not become part of the core Duel rules unless it represents an actual game rule.

The Story system is responsible for:

* progression;
* levels;
* opponents;
* duel sequence;
* rewards;
* Story-specific configuration.

---

# 10. Matchmaking

Multiplayer initially uses a simple queue-based matchmaking model.

The MVP does not include:

* private rooms;
* invitations;
* room codes;
* party systems.

Initial flow:

```text
Player 1
   ↓
Queue

Player 2
   ↓
Queue

Matchmaking
   ↓
Match

Player 1 ←→ Player 2
   ↓
Duel
```

The matchmaking system is separate from the Duel Engine.

The Duel begins only after two compatible players have been matched.

---

# 11. Multiplayer Lifecycle

The Multiplayer lifecycle is intentionally simple.

Conceptually:

```text
QUEUE
  ↓
MATCHED
  ↓
DUEL
  ↓
FINISHED
```

## QUEUE

The player is waiting for an opponent.

## MATCHED

Two players have been paired.

The system creates or initializes the Duel.

## DUEL

The players participate in the Duel Engine.

## FINISHED

The Duel reaches a terminal state.

Possible causes include:

* Life Points reaching zero;
* surrender;
* disconnect;
* other terminal conditions defined by Game Rules.

The Multiplayer lifecycle must not be represented as states inside the Duel when those states belong to matchmaking.

---

# 12. Turn and Phase

The current game has:

```text
PHASE_1
PHASE_2
```

Phase is modeled as a state/value rather than an independent entity.

## Turn

A Turn represents the current turn state.

Conceptually:

```text
Turn
├── Number
└── ParticipantID
```

The Duel coordinates turn transitions.

Turn-specific state must be reset when a new turn begins.

For example:

```text
CardInstance.AttacksPerformed = 0
```

for the appropriate monsters.

---

# 13. Zones

The main Duel zones are:

```text
Deck
Hand
Field
Graveyard
```

Each zone protects its own invariants.

---

## 13.1 Deck

Responsibilities include:

* drawing cards;
* adding/removing cards when permitted;
* checking whether cards remain;
* protecting deck invariants.

Deck exhaustion follows the Game Rules.

---

## 13.2 Hand

The Hand represents cards currently held by a Participant.

Responsibilities include:

* adding cards;
* removing cards;
* checking membership;
* enforcing hand-related invariants.

The hand limit and draw behavior are defined by Game Rules.

---

## 13.3 Field

The Field represents cards currently in play.

Current constraints:

```text
Maximum Monsters: 4
Maximum Magic/Trap: 4
```

The Field is responsible for protecting these invariants.

Conceptually:

```text
Field
├── Monsters
└── MagicTraps
```

---

## 13.4 Graveyard

The Graveyard contains cards moved there according to the Game Rules.

Destroyed cards are not deleted.

A CardInstance remains the same instance while changing location.

Example:

```text
DECK
  ↓
HAND
  ↓
FIELD
  ↓
GRAVEYARD
  ↓
FIELD
```

Revived cards remain the same CardInstance unless future rules explicitly introduce different behavior.

---

# 14. Combat

Combat resolution is delegated to a domain service.

## CombatResolver

`CombatResolver` calculates the result of combat.

It should not directly mutate:

* Life Points;
* Field;
* Graveyard;
* Duel state.

Conceptually:

```text
Attack
  ↓
Duel validates action
  ↓
CombatResolver
  ↓
CombatResult
  ↓
Duel applies result
```

---

## CombatContext

`CombatContext` contains the information necessary to resolve combat.

Potential information includes:

* attacker;
* defender;
* current game state;
* applicable effects;
* relevant modifiers.

---

## CombatResult

`CombatResult` represents the result of a combat calculation.

Conceptually:

```text
CombatResult
├── Attacker
├── Defender
├── AttackerDestroyed
├── DefenderDestroyed
├── Damage
└── ...
```

The result is calculated before the Duel mutates its state.

---

# 15. Effects

Effects represent behavior that modifies or interacts with the normal game rules.

Effects may be:

* passive;
* automatically triggered;
* activated under conditions.

The current architecture supports concepts such as:

```text
Piercing
DoubleAttack
IncreaseAttack
DecreaseAttack
Destroy
RecoverLifePoints
Draw
Search
Revive
FieldType
```

The exact effect taxonomy is defined by the Game Rules and card definitions.

The architecture should avoid spreading effect-specific conditionals throughout the game engine.

The system should favor composition and explicit effect behavior.

---

# 16. Attribute Modifiers

Effects that modify card attributes should not necessarily mutate the base Card definition.

For example:

```text
Card.Attack = 1500
```

should remain the card's base value.

A Duel may apply a modifier:

```text
Base Attack: 1500
Modifier: +500
Effective Attack: 2000
```

Modifiers may have durations.

Examples:

```text
Permanent for the Duel
Until end of turn
Until a specified condition
```

The Card definition remains immutable while the Duel state contains runtime modifications.

This prevents one Duel's effects from modifying the global definition of a Card.

---

# 17. Chain and Response System

The game supports response chains.

A Chain represents a sequence of activated effects.

Conceptually:

```text
Magic A
   ↓
Chain Link 1

Trap B
   ↓
Chain Link 2
```

Resolution occurs in Last-In-First-Out order:

```text
Chain Link 2
   ↓
Chain Link 1
```

The order of resolution depends on activation order, not card set order.

---

## 17.1 Chain

Conceptually:

```text
Chain
└── Links
    ├── ChainLink 1
    ├── ChainLink 2
    └── ...
```

A Chain Link represents one activation participating in the current chain.

The Chain is a Duel concept and exists only while a chain is active.

---

## 17.2 Response Window

Some actions allow a response window.

The current rules define a 5-second decision window for Trap activation.

The domain should represent the existence and deadline of the response window.

Infrastructure should be responsible for the actual passage of time.

The domain should not depend directly on goroutines, timers, WebSocket implementation, or HTTP.

---

## 17.3 Chain Resolution

Once chain resolution begins:

* no new activation occurs during resolution;
* links resolve in reverse activation order;
* an effect does not automatically activate another effect;
* each activation must satisfy its own conditions.

---

# 18. Domain Rules vs Application Logic

## Domain

The domain contains game rules and invariants.

Examples:

* field limits;
* deck rules;
* hand rules;
* summoning requirements;
* tribute requirements;
* combat calculations;
* attack limits;
* Life Point damage;
* card locations;
* valid phase actions;
* effect behavior;
* chain behavior;
* victory conditions.

The domain must not depend on:

* HTTP;
* REST;
* WebSocket;
* JSON;
* PostgreSQL;
* React;
* JWT;
* external APIs.

---

## Application Layer

The Application Layer coordinates use cases.

Examples:

```text
CreateUser
AuthenticateUser
CreateDeck
UpdateDeck
GetCollection

EnterMatchmaking
CreateMatch
StartDuel

StartStoryDuel
SubmitDuelAction
EndTurn
ActivateEffect
ResolveChain
FinishDuel
```

The Application Layer coordinates the domain and infrastructure without implementing game rules that belong inside the domain.

---

# 19. Dependency Direction

Dependencies should point toward the domain.

Conceptually:

```text
HTTP / WebSocket
       │
       ▼
 Application
       │
       ▼
    Domain
       ▲
       │
 Infrastructure
```

The domain must not depend on:

* HTTP;
* WebSocket;
* PostgreSQL;
* ORM;
* JWT;
* external services.

Repositories and infrastructure implementations depend on contracts defined by the appropriate inner layers.

The goal is to keep the Duel Engine executable and testable without a database or network transport.

---

# 20. Authentication and Security

The initial authentication mechanism is JWT.

Authentication belongs to the application/infrastructure boundary.

Conceptually:

```text
HTTP Request
     ↓
JWT Middleware
     ↓
Authenticated User
     ↓
Application Use Case
     ↓
Domain
```

The domain does not parse or validate JWT tokens.

The backend remains authoritative for:

* authentication;
* authorization;
* player ownership;
* deck ownership;
* Duel actions;
* game state;
* game rules.

The client must never be trusted to determine whether an action is valid.

Future security concerns include:

* token expiration;
* authorization;
* ownership validation;
* rate limiting;
* WebSocket session security;
* protection against manipulated requests.

---

# 21. Persistence

PostgreSQL is the primary persistent database.

Persistence is separated from the domain model.

Conceptually:

```text
Domain Model
     ↕
Persistence Model
```

The persistence model must not dictate the domain model.

Persistent concepts may include:

```text
Users
Collections
Decks
Card definitions
Story progress
Matches
Duel history
```

The runtime Duel state is not required to support recovery from server failure in the MVP.

If the server loses the in-memory Duel state, the active Duel may be considered terminated.

Persistent Duel recovery may be introduced in a future iteration if it becomes a real product requirement.

---

# 22. API

The initial general API uses REST.

REST is responsible for operations such as:

```text
POST /users
POST /auth/login

GET /collection
GET /decks
POST /decks
PUT /decks/{id}

POST /story/duels

POST /matchmaking/join
DELETE /matchmaking/queue
```

Duel actions may initially be exposed through application-oriented endpoints.

The exact endpoint structure may evolve as implementation progresses.

---

## DTOs

HTTP request and response DTOs must remain separate from domain entities.

The API must not expose internal domain objects directly as its public contract.

---

# 23. WebSocket

WebSocket is used for real-time Multiplayer communication.

The architecture separates:

```text
Transport
   ↓
Application
   ↓
Duel
```

The Duel Engine does not depend on WebSocket.

The WebSocket layer is responsible for transporting events and actions between connected players and the application layer.

Potential events include:

```text
MatchFound
DuelStarted
TurnStarted
ActionPerformed
ResponseWindowOpened
ChainActivated
ChainResolved
DuelFinished
OpponentDisconnected
```

The exact protocol will be defined when Multiplayer implementation begins.

WebSocket is not required for Story Mode.

---

# 24. Frontend Architecture

The frontend uses:

* React;
* TypeScript.

The frontend is responsible for:

* rendering the game;
* displaying player state;
* displaying cards;
* displaying available actions;
* interacting with REST APIs;
* maintaining the WebSocket connection during Multiplayer;
* displaying game events;
* displaying results.

The frontend is not the source of truth for game rules.

For example, the frontend may disable an invalid attack for UX purposes, but the backend must independently validate the action.

---

# 25. Story and Multiplayer Reuse

Story Mode and Multiplayer must use the same Duel Engine.

Conceptually:

```text
                    Duel Engine
                         │
              ┌──────────┴──────────┐
              │                     │
          Story Mode           Multiplayer
              │                     │
         AI Controller        Player Controller
              │                     │
          Character                User
```

The difference exists in the action provider, not in the game rules.

This ensures that:

* combat behaves identically;
* effects behave identically;
* turns behave identically;
* chains behave identically;
* victory conditions behave identically.

---

# 26. Testing Strategy

Testing prioritizes the domain because game rules are the most critical part of the system.

## Unit Tests

Used heavily for:

* combat;
* summoning;
* Set;
* tribute calculation;
* field limits;
* hand limits;
* draw rules;
* deck exhaustion;
* turn transitions;
* phase validation;
* attack limits;
* Double Attack;
* Piercing;
* effects;
* modifiers;
* chain resolution;
* Life Point calculations;
* card movement;
* Graveyard;
* revival;
* victory conditions.

Examples:

```text
Weight 5 → 0 tributes
Weight 6 → 1 tribute
Weight 9 → 2 tributes
Weight 13 → 3 tributes

Normal monster → maximum 1 attack
Double Attack → maximum 2 attacks

Destroyed monster → Graveyard

Chain A → Chain B
Resolution → B then A
```

---

## Integration Tests

Used for:

* PostgreSQL repositories;
* authentication;
* application use cases;
* matchmaking;
* REST API;
* WebSocket communication;
* persistence behavior.

---

## End-to-End Tests

Used for complete flows such as:

```text
Create User
   ↓
Login
   ↓
Receive Starter Deck
   ↓
Create/Edit Deck
   ↓
Enter Story
   ↓
Start Duel
   ↓
Draw
   ↓
Summon
   ↓
Attack
   ↓
Finish Duel
   ↓
Receive Reward
```

And Multiplayer:

```text
Login Player A
      ↓
Join Queue

Login Player B
      ↓
Join Queue

Match
  ↓
Duel
  ↓
WebSocket events
  ↓
Duel Finished
```

---

# 27. Error Handling

Errors should be separated by responsibility.

## Domain Errors

Represent invalid game operations.

Examples:

```text
InvalidPhase
InvalidAction
InvalidAttack
FieldFull
InsufficientTributes
CannotAttack
InvalidTarget
InvalidChainActivation
InvalidPosition
DeckEmpty
HandFull
```

## Application Errors

Represent failures in use-case orchestration.

Examples may include:

```text
UserNotFound
DeckNotFound
UnauthorizedDeckAccess
MatchNotFound
DuelNotFound
```

## HTTP Errors

Translate application and domain errors into appropriate HTTP responses.

The domain must not know about HTTP status codes.

---

# 28. Observability

The backend should eventually provide structured observability.

## Logging

Important events may include:

```text
UserRegistered
UserAuthenticated
MatchmakingJoined
MatchFound
DuelStarted
TurnStarted
CardSummoned
CardSet
AttackPerformed
MonsterDestroyed
CardMovedToGraveyard
EffectActivated
ChainStarted
ChainResolved
DuelFinished
PlayerDisconnected
```

## Metrics

Potential metrics include:

* active duels;
* players in matchmaking;
* matchmaking duration;
* duel duration;
* actions per duel;
* combat results;
* WebSocket connections;
* API latency;
* error rates.

## Tracing

Tracing may be introduced when the system becomes sufficiently complex or distributed, or when debugging complex flows provides measurable value.

Observability should be introduced incrementally.

---

# 29. AI Engineering Architecture

AI is treated as a development and engineering capability rather than as authority over the game domain.

The project may use AI agents to assist with:

* planning;
* architecture;
* implementation;
* testing;
* code review;
* documentation;
* debugging.

The AI workflow follows a human-in-the-loop model.

Conceptually:

```text
Human
  │
  ▼
Planner
  │
  ▼
Architect
  │
  ▼
Human Approval
  │
  ▼
Developer
  │
  ▼
Tester
  │
  ▼
Reviewer
  │
  ▼
Human Approval
```

Agents must not implicitly change game rules or architectural decisions.

---

# 30. Agent Responsibilities

## Planner

Responsible for transforming a requirement into an implementation plan.

The Planner should:

* inspect relevant context;
* identify affected modules;
* identify domain rules;
* identify required tests;
* identify risks.

The Planner does not modify production code.

---

## Architect

Responsible for evaluating architectural impact.

The Architect should:

* validate proposed designs;
* identify coupling;
* identify abstraction problems;
* verify architectural consistency;
* propose ADRs when necessary.

---

## Developer

Responsible for implementing an approved plan.

The Developer may modify source code and tests according to its permissions.

---

## Tester

Responsible for validating behavior.

The Tester should:

* execute tests;
* identify missing scenarios;
* verify domain rules;
* detect regressions.

---

## Reviewer

Responsible for independent review.

The Reviewer evaluates:

* correctness;
* architecture;
* tests;
* domain rule compliance;
* maintainability;
* security concerns.

---

# 31. AI Context

AI agents receive context from controlled project sources.

Primary sources include:

```text
docs/game-rules.md
docs/product-requirements.md
docs/architecture.md
docs/glossary.md
docs/adr/
.ai/context/
```

The Game Rules document is the source of truth for game behavior.

The Product Requirements document is the source of truth for product capabilities and MVP scope.

The Architecture document is the source of truth for approved technical architecture.

If an agent encounters an undefined rule, it must not invent one.

It must identify the ambiguity and request human clarification.

---

# 32. AI Rules and Guardrails

Agents must follow explicit project rules.

Examples:

* Never modify game rules implicitly.
* Never modify architecture decisions without approval.
* Never introduce dependencies without justification.
* Never bypass domain validation.
* Always create tests for new domain behavior.
* Inspect affected modules before implementation.
* Prefer existing abstractions when appropriate.
* Ask for clarification when requirements conflict or are undefined.
* Do not perform destructive operations without explicit authorization.
* Do not modify unrelated modules without justification.

---

# 33. AI Skills

Skills represent repeatable development procedures.

Initial skills may include:

```text
implement-feature
create-domain-tests
architecture-review
review-pr
debug-test-failure
create-adr
```

A Skill defines the procedure an agent follows.

A Skill does not define game rules.

---

# 34. AI Tool Access

Agents follow least-privilege principles.

Example:

```text
Planner
  READ: project/context/source
  WRITE: none

Architect
  READ: source/context/docs
  WRITE: proposals/ADR drafts

Developer
  READ: source/context/docs
  WRITE: source/tests

Tester
  READ: source/tests
  EXECUTE: tests

Reviewer
  READ: source/tests/diff/docs
  WRITE: review
```

Agents should only receive tools required for their responsibilities.

---

# 35. MCP and Tooling

MCP may expose controlled capabilities to AI agents.

Potential tools include:

```text
read_file
search_code
write_file
run_tests
run_linter
inspect_git_diff
search_documentation
```

Tool permissions must be scoped by agent.

The objective is to prevent agents from receiving unrestricted access to the development environment.

---

# 36. RAG

RAG may be introduced when the project's knowledge base becomes sufficiently large that direct context injection becomes inefficient.

Potential sources include:

* Game Rules;
* Product Requirements;
* Architecture;
* ADRs;
* card documentation;
* effect documentation;
* technical decisions;
* tests.

RAG is not a mandatory architectural component.

It should be introduced only when retrieval provides measurable value.

---

# 37. Deterministic Tooling

Deterministic operations should preferably be implemented as scripts or tools instead of relying on LLM reasoning.

Examples:

```text
run-tests
run-linter
validate-rules
generate-fixtures
check-formatting
```

AI agents should use deterministic tooling whenever a task can be objectively verified.

---

# 38. Human-in-the-Loop

Human approval is required for high-impact decisions.

Examples:

* changing game rules;
* changing architectural style;
* introducing major dependencies;
* changing domain boundaries;
* modifying security boundaries;
* approving architectural decisions.

AI should recommend and execute approved changes, not silently become the project's decision maker.

---

# 39. AI Evaluation

The AI workflow should eventually be evaluated using representative development tasks.

Examples:

```text
Add Double Attack
Add Piercing
Add Revive
Add a new Card Effect
Modify combat behavior
Add a new domain rule
Refactor combat
Add Multiplayer behavior
```

Evaluation dimensions may include:

* correctness;
* architecture compliance;
* test quality;
* rule compliance;
* unnecessary changes;
* tool usage;
* human interventions;
* failure rate.

The objective is to evaluate the AI engineering system itself rather than treating successful code generation as proof of successful engineering.

---

# 40. Architectural Decisions

Important architectural decisions are documented as ADRs.

## ADR-001 — Modular Monolith

The project starts as a Modular Monolith because the domain is evolving and there is no demonstrated requirement for distributed services.

## ADR-002 — Card vs CardInstance

`Card` represents a card definition while `CardInstance` represents a concrete copy participating in a Duel.

## ADR-003 — Graveyard

Cards moved from the Field or other zones are not deleted when the Game Rules require them to enter the Graveyard.

## ADR-004 — CombatResolver

Combat resolution is delegated to a domain service to prevent the Duel from becoming responsible for every calculation.

## ADR-005 — Human-in-the-Loop AI

AI agents may assist development, but architectural and game-rule decisions remain under human control.

## ADR-006 — Story and Multiplayer Share the Duel Engine

Story Mode and Multiplayer use the same Duel Engine. The difference between AI and human control exists outside the core game rules.

## ADR-007 — REST and WebSocket

REST is used as the initial general-purpose API. WebSocket is introduced for real-time Multiplayer communication without coupling the Duel Engine to the transport mechanism.

## ADR-008 — Queue-Based Multiplayer

The MVP Multiplayer flow uses a simple matchmaking queue. Private rooms, invitations, and similar functionality are deferred.

## ADR-009 — JWT Authentication

JWT is used for initial authentication. Authentication concerns remain outside the domain.

## ADR-010 — No Duel Recovery in MVP

An active Duel does not require recovery after server failure in the MVP. If the server loses the active Duel state, the match may be terminated.

Additional ADRs should only be created when a meaningful architectural decision is made.

---

# 41. MVP Architecture Scope

The MVP includes:

```text
Authentication
    ↓
User / Player
    ↓
Collection
    ↓
Decks
    ↓
Starter Deck
    ↓
Duel Engine
    ├── Story Mode
    │     └── AI Characters
    │
    └── Multiplayer
          ├── Matchmaking Queue
          └── WebSocket
```

The MVP does not require:

```text
Private Rooms
Invitations
Card Shop
Stars
Betting
Real Money
Ranking
Seasons
Divisions
Tournaments
Social Chat
Achievements
Spectator Mode
Replay System
Duel Recovery
```

These features may be introduced later without being considered requirements of the initial architecture.

---

# 42. Evolution

The architecture is intentionally designed to evolve.

Current:

```text
                    Current
                       │
        ┌──────────────┼──────────────┐
        │              │              │
 Modular Monolith     REST        PostgreSQL
        │
        ├── Duel Engine
        ├── Story Mode
        ├── AI Controller
        ├── Matchmaking
        └── WebSocket Multiplayer
```

Future evolution may include:

```text
Authentication improvements
Private rooms
Invitations
Card Shop
Stars
Ranking
Tournaments
Duel persistence/recovery
Advanced observability
RAG knowledge system
More sophisticated AI workflows
```

Microservices should only be introduced when a concrete scalability, organizational, reliability, or deployment requirement justifies the added complexity.

The architecture should evolve based on observed requirements rather than speculative future needs.

---

# 43. Architectural Principles

The following principles guide the project:

1. **Domain first** — game rules belong to the domain.

2. **Backend is authoritative** — the client cannot determine game validity.

3. **Explicit state over implicit flags** — represent domain facts directly.

4. **Small responsibilities** — avoid God Objects.

5. **Calculate before mutating** — specialized resolvers calculate results before the game state applies them.

6. **Composition over conditionals** — effects should not become scattered type checks.

7. **Dependency inversion** — infrastructure depends on appropriate domain/application contracts.

8. **Test business rules heavily** — domain behavior is the highest-value test target.

9. **Same Duel Engine** — Story and Multiplayer must use the same game rules.

10. **Transport independence** — the Duel Engine must not depend on REST or WebSocket.

11. **AI assists, humans decide** — agents do not own architectural or game decisions.

12. **Least privilege for agents** — each agent receives only the tools and permissions it needs.

13. **Deterministic tools for deterministic work** — use scripts and tests instead of LLM reasoning whenever possible.

14. **Document meaningful decisions** — use ADRs for architectural decisions.

15. **Do not solve hypothetical scale** — introduce complexity only when justified by real requirements.

16. **Separate ownership from runtime state** — Collection and Deck represent persistent player data, while CardInstance represents a concrete card participating in a Duel.

17. **Undefined rules remain undefined** — neither humans nor AI should silently invent game mechanics.

18. **Infrastructure stays outside the domain** — authentication, persistence, transport, timers, and external services must not contaminate core game rules.
