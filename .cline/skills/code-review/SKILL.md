---
name: code-review
description: Independently review a Go change for correctness, game rule compliance, regressions, maintainability and missing tests, with severity-ranked findings. Use when reviewing a Developer's work.
---

# Code Review

## Independence

* Form your own view. Run `git status` and `git diff`, and read new untracked files yourself; do not rely on the Developer's summary.
* Passing tests do not prove correctness. Read the code against the plan, the acceptance criteria and the docs.
* You may run read-only commands (build, vet, tests). Do not modify production code or tests, and do not silently fix anything.

## What to check, in order

1. Correctness and edge cases.
2. Game rule compliance (`docs/mvp-story-mode.md` first, then `docs/game-rules.md`).
3. Architecture and domain boundaries (no infrastructure in the domain, invariants protected by the type).
4. Regressions and compatibility of exported APIs.
5. Tests: do they assert behavior, cover failures and boundaries, and would they fail if the code were wrong?
6. Maintainability, duplication, unnecessary abstractions or dependencies.
7. Security, when relevant.
8. Scope: nothing beyond the approved plan, no unrelated files.

## Findings

Each finding has: severity, file and line (`path:line`), the relevant code, the problem, why it matters, and a recommendation.

Severity: `BLOCKER`, `HIGH`, `MEDIUM`, `LOW`, `INFO`. Order findings from most to least severe. Do not invent requirements or reject a solution only because another exists.

## Output

Summary, findings, game rule compliance, architecture compliance, test assessment, security assessment, and a recommendation: `APPROVED`, `APPROVED_WITH_CHANGES` or `CHANGES_REQUIRED`.
