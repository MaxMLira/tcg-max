package user

import "context"

// Repository is the persistence contract for user accounts.
//
// Implementations must store users keyed by their canonical (lowercase)
// username and behave as a concurrency-safe map: every method may be called
// from any goroutine (S9 decision 5, in-memory reference implementation).
//
// The domain does not prescribe the backing store; S9 ships an in-memory
// reference implementation (internal/repository/memory) so that integration
// tests run without external services. A production store (PostgreSQL,
// following architecture.md) implements the same contract.
type Repository interface {
	// Create stores user. It must be atomic: concurrent attempts to create
	// the same username result in exactly one success and ErrUsernameTaken
	// for the rest. It returns ErrUsernameTaken when the canonical username
	// is already taken.
	Create(ctx context.Context, user User) error

	// ByID returns the user with the given id, or ErrUserNotFound.
	ByID(ctx context.Context, id UserID) (User, error)

	// ByUsername returns the user whose canonical username equals
	// NormalizeUsername(name). Uniqueness is case-insensitive: "Alice" and
	// "alice" resolve to the same account (S9 decision 6). It returns
	// ErrUserNotFound when no such user exists.
	ByUsername(ctx context.Context, name string) (User, error)
}
