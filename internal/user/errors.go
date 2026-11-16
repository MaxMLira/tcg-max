package user

import "errors"

// Errors describe user account failures. Callers distinguish them with
// errors.Is; every error wraps one of these sentinels plus the detail that
// explains the failure (domain convention, see internal/domain/errors.go).
var (
	// ErrInvalidUser reports a user that violates an account invariant
	// (missing identity, missing password hash).
	ErrInvalidUser = errors.New("invalid user")

	// ErrInvalidUsername reports a username that violates its format rules:
	// 3 to 32 characters made of [A-Za-z0-9_] (S9 decision 6).
	ErrInvalidUsername = errors.New("invalid username")

	// ErrUserNotFound reports a lookup that matched no stored user.
	ErrUserNotFound = errors.New("user not found")

	// ErrUsernameTaken reports the attempt to create an account whose username
	// already exists (uniqueness is case-insensitive, S9 decision 6).
	ErrUsernameTaken = errors.New("username already taken")
)
