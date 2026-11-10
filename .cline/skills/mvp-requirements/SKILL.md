---
name: mvp-requirements
description: Read the TCG game rules and the Story Mode MVP scope, treat the MVP document as authoritative for implementation, and extract acceptance criteria and constraints without inventing rules. Use before planning or designing any feature.
---

# MVP Requirements

## Sources

Read, in this order:

1. `docs/mvp-story-mode.md` - authoritative for what is implemented in the Story Mode MVP (decisions, assumptions, backlog with task IDs and acceptance criteria).
2. `docs/game-rules.md` - the general TCG rules.
3. When relevant: `docs/glossary.md`, `docs/architecture.md`, `docs/product-requirements.md`.

Use the backlog entry of the task in `docs/mvp-story-mode.md` as the starting point for scope and acceptance criteria.

## Precedence

* For the MVP, `docs/mvp-story-mode.md` and the blocks labeled "Story Mode MVP" in `docs/game-rules.md` prevail over the general rule they annotate.
* Text in `docs/game-rules.md` outside those blocks is the general game. It is not implemented unless the MVP document includes it.
* Anything listed under "Outside the MVP" is excluded (for example Trap gameplay, Burn, Power Up, Multiplayer, WebSocket, PostgreSQL).

## Rules of use

* Never invent a game rule. If a rule is needed and neither document defines it, say so explicitly.
* If the MVP document already defines the behavior, do not raise it as a question. If it states an assumption, use that assumption.
* If two documents conflict, report the conflict and apply the precedence above.
* Do not modify `docs/game-rules.md`, `docs/mvp-story-mode.md` or `docs/product-requirements.md`.

## What to produce

1. The rules that apply, each with its source (file and section).
2. What is in the MVP and what is explicitly out.
3. Acceptance criteria, testable and traceable to a rule.
4. Constraints (architecture, dependencies, domain boundaries).
5. Genuine ambiguities only: what is undefined, why it blocks, and the simplest assumption consistent with the documents.
