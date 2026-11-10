# Development Rules

## 1. Tests

New domain behavior must include appropriate tests.

Game Rules should be tested explicitly.

---

## 2. Implementation Order

Prefer:

```text
Understand
  ↓
Plan
  ↓
Review
  ↓
Implement
  ↓
Test
  ↓
Review
```

Do not start coding before understanding the affected domain.

---

## 3. Small Changes

Prefer small, focused changes.

Avoid combining:

* feature implementation;
* unrelated refactoring;
* architecture changes;
* dependency changes

in the same task.

---

## 4. Formatting

Run the project's formatter after modifying Go code.

---

## 5. Validation

Before considering a task complete:

* run relevant tests;
* run the full test suite when appropriate;
* inspect the Git diff;
* verify that no unrelated files changed.

---

## 6. Error Handling

Do not silently ignore errors.

Errors should be handled at the appropriate architectural layer.

---

## 7. API Boundaries

Do not expose domain entities directly as public API contracts.

Use DTOs at external boundaries.

---

## 8. Infrastructure

Infrastructure implementations must not leak technical concerns into the domain.

Examples:

* database models;
* HTTP requests;
* WebSocket connections;
* JWT claims.

These concerns belong outside the core domain.

---

## 9. Refactoring

Refactor when it improves the current design.

Do not perform speculative refactoring solely to prepare for hypothetical future requirements.

---

## 10. Completion

A task is not complete merely because the code compiles.

Completion requires:

* implementation;
* relevant tests;
* validation;
* review of the resulting changes.
