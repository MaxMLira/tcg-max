# Story Mode MVP — Scope, Decisions and Backlog

Deadline: **2026-11-12**.

**Status: all MVP product decisions are finalized.** This document is authoritative for the scope of the Story Mode MVP. It is a delivery slice, not a change to the product: `docs/game-rules.md` describes the broader game and marks MVP behavior in blocks labeled "Story Mode MVP"; where they differ, this document and those blocks prevail for the MVP. Full-game requirements stay in `docs/product-requirements.md`.

---

## 1. Finalized decisions

### 1.1 Platform

| Area | Decision |
| --- | --- |
| Persistence | In-memory repositories and in-memory duel state. No PostgreSQL. |
| Authentication | Username and password only. No email, social login, password recovery or external identity provider. |
| Password hashing | Argon2id (section 1.6). Never plaintext, never raw HMAC. |
| Sessions | JWT. |
| Backend | Go modular monolith. |
| Frontend | React + TypeScript + Vite. |
| Communication | REST. No WebSocket, no Multiplayer, no matchmaking. |
| Mode | Story Mode: campaign with levels and duels against a basic AI. |
| Rewards | Three predefined cards for a victory. |
| Retry | A player can retry a level after losing. |
| Timer | No 40-second turn timer in Story Mode. |
| AI | Chooses legal actions, prioritizing a good summon and then attacks. |

### 1.2 General duel rules

* Both players start with exactly **4 random cards from their own decks**.
* There is no additional opening card or draw advantage for the first player.
* The hand is replenished **automatically** to 4 cards whenever a resolved action leaves fewer than 4. There is no manual draw action.
* If the deck cannot provide enough cards to replenish the hand to 4, the player loses immediately.
* Hand limit: 4 cards.
* The first player cannot attack during their first turn.
* A duel ends immediately when a player's Life Points reach 0 or less.
* A full monster field prevents summoning unless tributes free a slot.
* Both Participants use the same Duel Engine; the engine has no human/AI-specific rules.

### 1.3 Card types and subtypes

* The domain supports three distinct card types: **Monster**, **Magic** and **Trap**.
* Trap remains a valid domain card type. Trap gameplay, activation and effects are **excluded from the MVP**.
* Magic and Trap cards have no ATK, DEF, Weight or Monster subtype. (Assumption A1 covers Trap stats.)
* **Every Monster has a subtype.** Magic and Trap have none.
* Monster subtypes are an **extensible concept**: the model must not depend on a permanently fixed list, and new subtypes must be addable without redesigning the domain.
* The MVP initially supports four subtypes: **Dragon, Warrior, Mage, Aquatic**.

### 1.4 Magic Equipment

The MVP includes exactly four Magic cards, all Magic Equipment, each granting **+800 ATK**:

| Card | Compatible subtype |
| --- | --- |
| Dragon Equipment | Dragon |
| Warrior Equipment | Warrior |
| Mage Equipment | Mage |
| Aquatic Equipment | Aquatic |

Rules:

* Equipment must first be placed in a free Magic/Trap zone (through the action in section 1.5).
* Activation is allowed only during **Phase 2**.
* The player selects an existing compatible Monster on their own field. The target must be **face-up**.
* Activation is rejected if no compatible face-up Monster exists.
* After activation the Equipment stays on the field, occupies its Magic/Trap zone and is linked to the target Monster.
* Several Equipment cards may be linked to the same Monster; bonuses add up.
* If the equipped Monster leaves the field, all Equipment linked to it moves to the Graveyard.
* Equipment cannot be activated from the hand.
* Combining Monster placement and Equipment activation in one action is excluded.
* An existing Magic/Trap card cannot be replaced when all four zones are occupied.
* Placement and activation are distinct steps. No Trap behavior is introduced.
* **Burn and Power Up are excluded.** No cards, effects, rules or implementation exist for them.

### 1.5 Card selection, discard and placement

One atomic player action, performed in Phase 1 and subject to one placement per Phase 1:

* The player selects **1 to 4** cards from the hand. Selection order matters.
* Every selected card except the last goes to the Graveyard.
* The last selected card is placed on the field and must satisfy all placement rules for its type.
* If the last card cannot legally be placed, **the entire action is rejected** and no selected card is discarded.
* On success, the hand is automatically replenished to 4 cards; if the deck cannot complete it, the player loses immediately.
* One selected card means placing that card with no discard. Four selected cards means the first three go to the Graveyard and the fourth is placed.
* This action **replaces the separate discard step**. There is no discard-only action.

### 1.6 Password hashing

Argon2id through `golang.org/x/crypto/argon2`:

| Parameter | Value |
| --- | --- |
| Memory | 64 * 1024 KiB (64 MiB) |
| Iterations | 3 |
| Parallelism | 2 |
| Salt | 16 bytes, from a cryptographically secure random source |
| Derived key length | 32 bytes |

The stored value must carry the algorithm, parameters, salt and derived hash in a verifiable format. Plaintext passwords are never stored or logged. Raw HMAC is not used as a password hash.

This adds one dependency, `golang.org/x/crypto`. It is needed because Argon2id is not in the Go standard library; it is the maintained official extended-library implementation.

### 1.7 Terminology

* **Duel** is the gameplay engine. **Match** is the multiplayer pairing concept.

---

## 2. Assumptions (simplest choice consistent with the decisions)

* **A1 — Trap card stats:** `NewCard` keeps its current Trap validation unchanged. Trap has no gameplay in the MVP.
* **A2 — Trap in the MVP:** Trap cards are not part of any MVP deck, and the placement action rejects a Trap.
* **A3 — Extensible subtype:** a subtype is a validated, non-blank value type rather than a closed enumeration. The four MVP subtypes are values of it. Equipment declares the one subtype it accepts.
* **A4 — Selection action:** it is only valid in Phase 1, selected cards must be distinct, and a Monster that needs tributes carries its tributes, position and face in the same action; invalid tributes reject the whole action.
* **A5 — Opening hand:** the 4 opening cards are drawn from the shuffled deck. The deck must hold at least 4 cards; decks are 30 cards.

## 3. Outside the MVP (future)

Multiplayer and matchmaking, WebSocket, PostgreSQL persistence, Collection and Deck builder, registration extras (email, social login, password recovery), Trap gameplay, Burn, Power Up, chains, revival, complex effects, the 40-second timer, the broader general draw/discard rules in `docs/game-rules.md` sections 5, 11–13.

## 4. Documentation notes

* `docs/glossary.md` section 16 lists `Duel` ↔ `Match` and `Match` ↔ `Duel` as preferred/avoid pairs. This contradicts the glossary definitions and section 1.7. It is outside this task and should be corrected in a separate docs change.
* `docs/product-requirements.md` section 14 puts Multiplayer and PostgreSQL in the product MVP. This milestone defers them, and in-memory state does not satisfy saved-progress requirements (PR-STORY-007/008).

---

## 5. Backlog

Agents: Planner, Developer, Reviewer, QA. Priority: Critical, High, Optional. No task depends on an open question; all decisions are final.

| ID | Task | Depends on | Acceptance criteria | Priority | Agent | Status |
| --- | --- | --- | --- | --- | --- | --- |
| S0 | Record MVP scope and fix module path | none | This document exists; `go.mod` module equals the Git remote path; build passes. | Critical | Planner | Done (uncommitted) |
| S1 | `NewDeckDefinition`, blank ID/name rejection, Magic without ATK/DEF/Weight | S0 | Constructor copies and validates; blank values rejected; Magic with non-zero ATK/DEF/Weight rejected; tests, gofmt, vet pass. | Critical | Developer, Reviewer, QA | Done (uncommitted, review pending) |
| S1b | Extensible Monster subtype on `Card` (sections 1.3, 1.4): Monster requires a valid subtype; Magic and Trap reject one; Equipment declares its one accepted subtype; the four MVP subtypes are provided | S1 | Monster without subtype rejected; Magic/Trap with a Monster subtype rejected; blank subtype rejected; a new subtype value works without domain changes; Equipment accepts only its subtype; tests pass. | Critical | Developer | Todo |
| S2 | CardInstance, zones (Deck, Hand, Field, Graveyard), Participant, LP by level | S1b | Hand limit 4; 4 monster slots; 4 Magic/Trap slots; shuffle injectable; LP 4000 for levels 1-4 and 8000 for levels 5-10. | Critical | Developer | Todo |
| S3 | Duel setup and automatic hand replenishment (1.2) | S2 | Both players start with exactly 4 random cards from their own decks; no extra card for the first player; replenishment to 4 after a resolved action; no manual draw; insufficient deck loses immediately; first player cannot attack on turn 1; fixed seed is deterministic. | Critical | Developer | Todo |
| S4 | Summon with tributes, face-down rules, change position, full-field rule, and the selection/discard/placement action (1.5) | S3 | Tribute tiers 0/1/2/3 by weight; full field blocks summon unless tributes free a slot; 1-4 selected cards, last placed, others to Graveyard; illegal last card rejects the whole action with no discard; hand replenished after success; no discard-only action; one placement per Phase 1. | Critical | Developer | Todo |
| S5 | Magic Equipment (1.4): placement in a free Magic/Trap zone, Phase 2 activation, linking, +800 ATK, Graveyard on monster leaving | S4 | Activation rejected outside Phase 2, from the hand, when not placed, or with no compatible face-up Monster; Equipment stays in its zone linked to the target; bonuses stack; matching subtype only; a monster leaving the field sends all linked Equipment to the Graveyard; a fifth Magic/Trap placement is rejected; base Card never mutated. | Critical | Developer | Todo |
| S6 | CombatResolver, attack action, direct attack, LP and win/lose | S4 | All combat cases of game-rules 26-27; ATK includes Equipment bonuses; direct attack only when the opponent has no monsters; duel ends at 0 LP or less. | Critical | Developer | Todo |
| S7 | Legal-action enumerator, controller interface, basic AI | S5, S6 | AI plays only enumerated actions; no human/AI branching in the engine; AI-vs-AI simulations finish without illegal actions. | Critical | Developer | Todo |
| S8 | Story domain: Character, Level, progress, seed data (Monster cards with subtypes, the four Equipment cards, decks), 3 predefined reward cards per victory | S1b | Seed decks pass `NewDeckDefinition`; no Trap, Burn or Power Up in seed data; retry after a loss allowed; a victory grants the three predefined cards. | Critical | Developer | Todo |
| S9 | Auth: username/password, Argon2id (1.6), JWT, in-memory user repository | none | Parameters exactly as in 1.6; random 16-byte salt; verifiable stored format; plaintext never stored or logged; duplicate username and wrong credentials rejected. | Critical | Developer | Todo |
| S10 | Application layer and in-memory duel registry: start story, list levels, start duel, perform action with AI turn, record result | S7, S8 | A scripted duel runs to the end and updates progress; a loss does not progress. | Critical | Developer | Todo |
| S11 | REST API, DTOs, error mapping, `cmd/server` | S9, S10 | DTOs only at the boundary; hidden information not leaked; legal actions returned. | Critical | Developer | Todo |
| S12 | Frontend (React + TypeScript + Vite): login, story map, duel board, result | S11 contract | The MVP flow is playable in a browser; no rules logic in the client. | Critical | Developer | Todo |
| S13 | End-to-end validation and bug-fix pass | S11, S12 | Full flow works from a clean start; `go test ./...` passes. | Critical | QA | Todo |
| S14 | Dockerfile, compose, README | S11 | One command starts the app. | Optional | Developer | Todo |

Work is delivered in separate backend, frontend and docs commits (commit-strategy rules) and reviewed before dependent tasks start.

## 6. Schedule

| Date | Plan |
| --- | --- |
| Nov 6 | S1b |
| Nov 7 | S2, S3, S8 |
| Nov 8 | S4, S5 |
| Nov 9 | S6, S7 |
| Nov 10 | S9, S10, S11 |
| Nov 11 | S12, S13 |
| Nov 12 | Final QA, review, split commits |

The frontend can start against a mocked contract once the S11 DTOs are agreed. If time runs short, cut S14 first.
