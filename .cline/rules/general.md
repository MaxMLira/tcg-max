# General AI Rules

## 1. Source of Truth

Use the project documentation as the source of truth.

Do not invent requirements or game mechanics.

Relevant sources include:

* `docs/game-rules.md`
* `docs/product-requirements.md`
* `docs/architecture.md`
* `docs/glossary.md`
* `docs/groud-rules.md`
* `docs/adr/`

---

## 2. Undefined Behavior

If a requirement or Game Rule is undefined:

* do not guess;
* do not invent behavior;
* report the ambiguity;
* request clarification when necessary.

---

## 3. Scope

Only modify files that are relevant to the requested task.

Do not perform unrelated refactors.

---

## 4. Existing Code

Inspect existing code before creating new abstractions.

Prefer reusing existing abstractions when appropriate.

---

## 5. Dependencies

Do not introduce a new dependency without explaining:

* why it is necessary;
* what problem it solves;
* what alternatives were considered.

---

## 6. Destructive Operations

Do not perform destructive operations without explicit authorization.

Examples:

* deleting files;
* deleting data;
* rewriting large areas of the project;
* changing persistent data structures.

---

## 7. Human Approval

Human approval is required for:

* Game Rule changes;
* architectural changes;
* major dependencies;
* security boundary changes;
* domain boundary changes.

---

## 8. Communication

When encountering uncertainty, clearly distinguish:

* known facts;
* assumptions;
* recommendations;
* open questions.
