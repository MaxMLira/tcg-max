# Glossary

## Purpose

This document defines the terminology used throughout the project.

The goal is to ensure that developers, documentation, AI agents, and game rules use the same terminology consistently.

The Glossary does not define game rules or architecture decisions.

* `game-rules.md` defines **how the game works**.
* `product-requirements.md` defines **what the product must provide**.
* `architecture.md` defines **how the software is structured**.
* `glossary.md` defines **what project terms mean**.

---

# 1. Core Game Concepts

## Duel

A single match between two participants.

A Duel contains the complete runtime state required to execute a game.

A Duel may occur in:

* Story Mode;
* Multiplayer.

The Duel Engine is responsible for applying the game rules.

---

## Participant

An entity participating in a Duel.

A Duel always has two Participants.

A Participant may be controlled by:

* a human player;
* an AI controller.

The Duel Engine does not need to know how the Participant is controlled.

---

## Player

A User participating in the game.

For the current architecture:

`User = Player`

A Player may own:

* a Collection;
* multiple Decks;
* Story progress.

---

## Character

A named opponent used by Story Mode.

Example:

`Richard`

A Character is controlled by an AI Controller during a Story Duel.

---

## AI Controller

A component responsible for deciding actions for an AI-controlled Participant.

The AI Controller does not implement the game rules.

It selects actions that are then validated by the Duel Engine.

---

# 2. Card Concepts

## Card

The definition/template of a card.

A Card describes what the card is.

Examples of Card attributes include:

* name;
* type;
* attack;
* defense;
* weight;
* effect.

A Card does not represent a specific copy currently participating in a Duel.

---

## CardInstance

A concrete copy of a Card participating in a Duel.

A CardInstance has runtime state.

Examples:

* owner;
* location;
* position;
* attacks performed.

The same Card definition may have multiple CardInstances.

---

## Collection

The set of cards owned by a Player.

The Collection represents persistent ownership.

A Collection is different from a Deck and from CardInstances participating in a Duel.

---

## Deck

A collection of cards selected by a Player for use in a Duel.

A Player may own multiple Decks.

Deck construction rules are defined in `game-rules.md`.

---

## Starter Deck

A predefined Deck provided to a newly registered Player.

The Player receives one Starter Deck selected from the predefined Starter Deck options.

---

# 3. Zones

## Zone

A location where cards can exist during a Duel.

The primary zones are:

* Deck;
* Hand;
* Field;
* Graveyard.

---

## Deck Zone

The collection of CardInstances that have not yet been drawn during a Duel.

The Deck is also used as a persistent concept for deck composition, but the runtime Duel Deck consists of CardInstances.

---

## Hand

The cards currently held by a Participant during a Duel.

The Hand has a maximum size defined by the Game Rules.

---

## Field

The area where cards currently in play are placed.

The Field contains:

* Monster cards;
* Magic/Trap cards.

The Field has separate limits for Monsters and Magic/Trap cards.

---

## Graveyard

A public zone containing cards that have been moved there according to the Game Rules.

Cards in the Graveyard are not deleted.

A CardInstance remains the same instance when moving to the Graveyard.

---

# 4. Card Types

## Monster

A Card Type representing a Monster card.

A Monster may have:

* Attack;
* Defense;
* Weight;
* Position;
* Effects.

---

## Magic

A Card Type representing a Magic card.

Magic cards must be set before activation according to the Game Rules.

---

## Trap

A Card Type representing a Trap card.

Trap cards must be set before activation and may respond to appropriate game events.

---

# 5. Monster Concepts

## Weight

A Monster attribute used to determine the number of Tributes required for Summoning.

The exact Weight-to-Tribute relationship is defined in `game-rules.md`.

---

## Attack Position

A Monster position in which the Monster may perform attacks.

---

## Defense Position

A Monster position in which the Monster cannot attack.

---

## Face-Up

A card state in which the card's information is visible to both players.

---

## Face-Down

A card state in which the card's information is hidden from the opponent.

A face-down Monster is placed in Defense Position.

---

## Summon

The action of placing a Monster onto the Field according to the Summoning rules.

Summon is a game action.

It is not an Effect Type.

---

## Set

The action of placing a Monster face-down in Defense Position or placing a Magic/Trap card face-down according to the Game Rules.

Set and Summon are distinct game actions even though both place cards onto the Field.

---

## Tribute

A Monster given up from the Field as part of a Summoning requirement.

Tributed Monsters are moved to the Graveyard.

---

# 6. Turn Concepts

## Turn

A period in which one Participant is allowed to perform the actions permitted by the current game state.

A Turn has:

* a number;
* a Participant;
* a current Phase.

---

## Phase

A stage of the Turn that determines which actions are currently permitted.

The current game contains:

* Phase 1;
* Phase 2.

---

## Phase 1

The phase responsible for placement and position-related actions defined by the Game Rules.

---

## Phase 2

The phase responsible for attacks, Magic activation, and other permitted actions defined by the Game Rules.

---

## First Turn

The first Turn of a Duel.

The first player follows special restrictions defined by the Game Rules.

---

# 7. Combat Concepts

## Attack

An action in which a Monster attempts to attack another Monster or directly attack the opposing Participant when permitted.

---

## Direct Attack

An Attack performed against the opposing Participant when no opposing Monsters are present, subject to the Game Rules.

---

## Combat

The process of resolving an Attack.

Combat determines outcomes such as:

* destruction;
* Life Point damage;
* survival.

---

## CombatResolver

A domain service responsible for calculating the result of Combat.

It calculates the result but does not directly mutate the Duel state.

---

## Piercing

An effect that allows excess Attack power against a Defense-position Monster to cause Life Point damage according to the Game Rules.

---

## Double Attack

An effect that allows a Monster to perform up to two attacks during the appropriate Turn.

---

# 8. Effects

## Effect

A behavior associated with a card that modifies or interacts with normal game behavior.

Effects may be:

* passive;
* automatically triggered;
* activated under specific conditions.

---

## Effect Type

A classification describing the behavior performed by an Effect.

Examples include:

* `ADD_TO_HAND`;
* `SEARCH`;
* `DESTROY`;
* `RECOVER_LP`;
* `BUFF_ATK`;
* `DEBUFF_ATK`;
* `FIELD_TYPE`.

An Effect Type is not necessarily a unique card.

Multiple cards may use the same Effect Type with different conditions or values.

---

## Passive Effect

An Effect that continuously modifies or influences game behavior without requiring a player activation.

Example:

`Double Attack`

---

## Activated Effect

An Effect that is activated when its conditions are satisfied.

An activated Effect may have a cost.

---

## Modifier

A runtime change applied to an attribute without changing the base Card definition.

Example:

```text
Base ATK: 1500
Modifier: +500
Effective ATK: 2000
```

Modifiers may have a duration.

---

## Duration

The period during which a Modifier or Effect remains active.

Examples include:

* entire Duel;
* until end of Turn;
* until a specific condition.

---

# 9. Chain Concepts

## Chain

A sequence of activated Effects or card activations that are resolved in reverse activation order.

---

## Chain Link

A single activation participating in a Chain.

---

## Chain Resolution

The process of resolving Chain Links.

Resolution occurs in Last-In-First-Out order.

---

## Response Window

A limited period during which a Participant may respond to an event with an eligible action.

The current rules define a 5-second response window for Trap decisions.

---

# 10. Game State

## Duel State

The complete runtime state of a Duel.

It includes information such as:

* Participants;
* Life Points;
* Turn;
* Phase;
* zones;
* CardInstances;
* active Chain;
* Response Window;
* Duel Status.

---

## Duel Status

The current lifecycle state of a Duel.

Examples may include:

* waiting to start;
* active;
* finished.

The exact valid states are defined by the implementation and Game Rules.

---

## Life Points (LP)

The amount of life remaining for a Participant.

A Participant loses when their Life Points reach zero or below according to the Game Rules.

---

# 11. Multiplayer Concepts

## Matchmaking

The process of finding an opponent for a Player.

The MVP uses a queue-based matchmaking system.

---

## Queue

The waiting state in which a Player is looking for an opponent.

The MVP does not include private rooms or invitations.

---

## Match

The pairing of two Players who will participate in a Multiplayer Duel.

---

## Multiplayer

The game mode in which two human Players participate in the same Duel.

---

# 12. Story Mode Concepts

## Story Mode

The game mode in which the Player progresses through predefined levels and fights AI-controlled Characters.

---

## Story Level

A progression stage containing a predefined sequence of Duels.

The current product contains 10 Levels with 4 Duels per Level.

---

## Story Opponent

A Character faced by the Player during Story Mode.

Example:

`Richard`

---

## AI Strategy

The decision-making behavior used by an AI Controller to select actions during a Duel.

The AI Strategy must operate within the Game Rules.

---

# 13. Product Concepts

## Reward

A card or other progression-related benefit granted after a successful Duel according to the Product Requirements and Game Rules.

The current Story Mode reward is 3 cards per victory.

---

## Card Shop

A future product feature through which Players may obtain cards.

The Card Shop is outside the MVP.

---

## Stars

An in-game currency planned for a future economy system.

Stars and betting are outside the MVP.

---

## Wager

A future Multiplayer mechanic in which Players agree to risk Stars on the outcome of a Duel.

Wagers are outside the MVP.

---

# 14. Architecture Terms

## Duel Engine

The domain component responsible for executing the rules of a Duel.

The Duel Engine is shared by Story Mode and Multiplayer.

---

## Application Layer

The layer responsible for coordinating application use cases.

It connects external requests with domain behavior and infrastructure.

---

## Infrastructure

Technical implementations required by the application but not part of the game domain.

Examples include:

* PostgreSQL;
* HTTP;
* WebSocket;
* JWT;
* repositories.

---

## Repository

An abstraction used to access persistent data without coupling the domain to a specific database implementation.

---

## DTO

Data Transfer Object.

A structure used to transfer data across application boundaries.

API DTOs should not be the same objects as domain entities.

---

## Modular Monolith

An architecture in which the application is deployed as a single process/application while maintaining clear internal module boundaries.

The project uses a Modular Monolith as its initial architecture.

---

# 15. AI Engineering Terms

## AI Agent

An AI system configured to perform a specific development responsibility.

Examples:

* Planner;
* Architect;
* Developer;
* Tester;
* Reviewer.

---

## Context

The information provided to an AI Agent to perform its task correctly.

Project context may include:

* Game Rules;
* Product Requirements;
* Architecture;
* Glossary;
* ADRs;
* source code.

---

## Rule

An explicit instruction that constrains AI Agent behavior.

Rules are different from Game Rules.

Game Rules define the game.

AI Rules define how an AI Agent must behave while working on the project.

---

## Skill

A reusable procedure that defines how an AI Agent should perform a particular task.

Examples:

* implement-feature;
* create-domain-tests;
* architecture-review;
* review-pr.

---

## MCP

Model Context Protocol.

A mechanism through which controlled tools and capabilities can be exposed to AI Agents.

---

## RAG

Retrieval-Augmented Generation.

A mechanism for retrieving relevant project information before providing context to an AI Agent.

RAG is optional and should only be introduced when it provides measurable value.

---

## Human-in-the-Loop

A workflow in which humans retain authority over important decisions while AI Agents assist with analysis and implementation.

Architectural and Game Rule changes require human approval.

---

# 16. Naming Principles

The following terms should be used consistently throughout the project:

| Preferred      | Avoid        |
| -------------- | ------------ |
| `Card`         | CardTemplate |
| `CardInstance` | RuntimeCard  |
| `Collection`   | Inventory    |
| `Deck`         | CardList     |
| `Participant`  | Opponent     |
| `Character`    | AIPlayer     |
| `AIController` | AIPlayer     |
| `Duel`         | Match        |
| `Match`        | Duel         |
| `Graveyard`    | DiscardPile  |
| `Summon`       | InvokeEffect |
| `Effect`       | CardAction   |
| `Chain`        | EffectStack  |
| `Phase`        | GameStep     |
| `Matchmaking`  | LobbySystem  |
| `Queue`        | WaitingRoom  |

These conventions exist to maintain consistency between Game Rules, Product Requirements, Architecture, source code, tests, and AI-generated work.
