# Workflow: feature-development

Invoke with `/feature-development.md` followed by a feature description or a backlog ID from `docs/mvp-story-mode.md` (for example `S1b`).

The main agent is the **Orchestrator**. It never plays the other roles itself. Every role runs through a real `spawn_agent` call.

## Mandatory rules

1. Delegate with `spawn_agent` only. Never write role-based answers (Planner, Architect, Developer, Reviewer, Tester) in the main context.
2. Start independent agents in the same response so they run in parallel: Planner + Architect (Step 2), Reviewer + Tester (Step 5).
3. Pass earlier reports explicitly to downstream agents. Subagents do not see this conversation.
4. Keep each subagent task bounded: one role, one deliverable, a stated report format.
5. Planner and Architect never write or edit files. Reviewer never edits production code or tests. Tester never edits production code.
6. The Developer implements only the plan the user approved.
7. **No implementation before the user explicitly approves the plan** (Step 3).
8. Never commit or push. Commits happen only if the user asks for them.
9. If `spawn_agent` fails (for example `Unauthorized`), report the failure and the affected step, then stop. Do not do the agent's work yourself and do not present the step as done.
10. Do not ask product questions that `docs/mvp-story-mode.md` or `docs/game-rules.md` already answer. Ask only about genuine gaps.
11. Claim an agent or Skill as used only if the call was actually made, or the Skill file was actually read and injected.

## How agents and Skills are wired

Cline does not load `.cline/agents/*.md` as native profiles and has no setting that attaches Skills to a profile. `spawn_agent` accepts only `systemPrompt` and `task`. So the Orchestrator builds each `systemPrompt` itself:

1. Read the agent file with `read_files` and use its text as the base of `systemPrompt`.
2. Read each Skill file listed below with `read_files` and append its text under the heading `## Skill: <name>`.
3. Append this path note, because the agent files predate the repository layout:
   `Path mapping: .ai/rules/ = .cline/rules/ ; .ai/context/ = .ia/context/ ; docs/commit-strategy.md = .cline/rules/commit-strategy.md . docs/mvp-story-mode.md is authoritative for MVP scope.`
4. Append the constraints of the role (the Mandatory rules above that apply to it).

| Role | Agent file | Skills to read and inject |
| --- | --- | --- |
| Planner | `.cline/agents/planner.md` | `.cline/skills/mvp-requirements/SKILL.md` |
| Architect | `.cline/agents/architect.md` | `.cline/skills/mvp-requirements/SKILL.md`, `.cline/skills/go-domain-design/SKILL.md` |
| Developer | `.cline/agents/developer.md` | `.cline/skills/go-implementation/SKILL.md`, plus `.cline/skills/go-domain-design/SKILL.md` when domain code changes |
| Reviewer | `.cline/agents/reviewer.md` | `.cline/skills/code-review/SKILL.md` |
| Tester | `.cline/agents/tester.md` | `.cline/skills/go-testing/SKILL.md` |

Project rules in `.cline/rules/` apply to every role. Tell each agent to read them.

## Step 1 - Receive the feature (Orchestrator, no subagent)

1. Take the feature description or backlog ID from the user.
2. Read `docs/mvp-story-mode.md`, the relevant parts of `docs/game-rules.md`, `docs/architecture.md` and `docs/glossary.md`, the backlog entry, and the existing code and tests it touches.
3. Write down scope, constraints and which docs sections apply.
4. Do not modify anything.

## Step 2 - Parallel planning (two `spawn_agent` calls in one response)

**Planner** receives: the feature request, the backlog entry, the relevant MVP and game-rules sections, and the scope notes from Step 1.
Deliverable: objective, relevant rules with sources, affected modules, small ordered implementation tasks with dependencies, acceptance criteria, tests required, risks, genuine open questions. No code.

**Architect** receives: the feature request, the backlog entry, the relevant MVP sections, and pointers to `docs/architecture.md` and the source files to inspect.
Deliverable: assessment, affected boundaries, simplest compatible design (types, signatures, invariants), files to change, risks and trade-offs, alternatives, ADR recommendation, status `APPROVED`, `APPROVED_WITH_CHANGES` or `REJECTED`. No code.

The two agents work independently. Wait for both. If one fails, report which and stop.

## Step 3 - Consolidate and request approval (Orchestrator)

Present one plan to the user:

* feature scope and what is excluded;
* acceptance criteria;
* technical design;
* files expected to change;
* implementation tasks in order;
* risks, dependencies and any conflict between the Planner and Architect reports.

Then ask: "Do you approve this plan for implementation?"

**STOP. Take no further action until the user answers.**

* Approval means an explicit yes. Silence, a question or a partial comment is not approval.
* If the user asks for changes, revise the plan (re-run Planner and/or Architect if needed) and ask again.

## Step 4 - Implementation (only after explicit approval)

Call the **Developer** with `spawn_agent`. Provide:

* the original request and the **approved** plan, exactly as approved;
* the Architect report and the acceptance criteria;
* the instructions: implement only the approved plan; add or update tests; run `gofmt -l .`, `go vet ./...`, `go test ./...`; inspect `git diff`; do not commit; stop and report if an ambiguity or architectural problem appears.

Deliverable: what was implemented, files changed, tests added or changed, commands executed with real output, design decisions, remaining issues.

Wait for the report. If the Developer fails or reports a blocking ambiguity, tell the user and stop. Do not commit.

## Step 5 - Independent validation (two `spawn_agent` calls in one response)

Before calling them, the Orchestrator collects `git status --short` and the diff of the changed files.

Both **Reviewer** and **Tester** receive: the original request, the approved plan, the acceptance criteria, the Developer report, the list of changed files and the diff. Tell them not to rely on the Developer's conclusions: they must read the code and run commands themselves.

* **Reviewer** - correctness, game rules, design, regressions, scope. Read-only: no edits to production code or tests. Deliverable: findings with severity and `path:line` evidence, and a recommendation.
* **Tester** - acceptance criteria, scenarios, execution. May add or update test files only when a scenario is missing and must list them. Never touches production code. Deliverable: tests executed with exact output, missing scenarios, rule violations, regression risks, status.

Wait for both. If either fails, report it and stop. The step is not complete without both reports.

## Step 6 - Quality gate (Orchestrator)

Consolidate both reports.

**Pass** only when all of these hold: no `BLOCKER` or `HIGH` finding, Reviewer recommendation is `APPROVED` or `APPROVED_WITH_CHANGES` with no blocking change, Tester status is `PASS` or `PASS_WITH_WARNINGS`, and every acceptance criterion is shown to pass. On pass, present the validation report and ask the user to approve the final result. Do not commit.

**Otherwise**, present the findings, send the corrections to the Developer with `spawn_agent` (pass the findings verbatim), then repeat Step 5 on the new diff.

* At most **two** correction rounds.
* If blocking issues remain after the second round, stop and ask the user to intervene. Do not start a third round.
* Corrections stay inside the approved plan. A correction that needs a new decision goes back to the user.

## Step 7 - Final report (Orchestrator)

Report:

1. Feature implemented.
2. Agents actually invoked (role, step, result). Say so for any that failed or were skipped.
3. Skills actually used (read and injected).
4. Files changed.
5. Tests executed and exact results.
6. Reviewer findings.
7. Tester findings.
8. Remaining risks.
9. Whether the user must approve a commit. Offer the commit split required by `.cline/rules/commit-strategy.md`, but do not commit.
