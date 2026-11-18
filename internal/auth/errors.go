// Package auth implements the S9 authentication use cases: user
// registration, login, and JWT token verification.
//
// It owns password hashing (Argon2id, PHC string format) and HS256 JWT
// issuance/verification built only on the standard library. It depends on
// internal/user for the account model and repository contract, and never on
// any transport (HTTP, WebSocket, REST): transport layers wrap this package.
package auth

import "errors"

// Errors describe authentication failures. Callers distinguish them with
// errors.Is; every error wraps one of these sentinels plus the detail that
// explains the failure.
var (
	// ErrInvalidInput reports a use-case argument that is empty or malformed
	// (empty username/password on Register, empty subject on Issue, nil
	// dependency on New).
	ErrInvalidInput = errors.New("invalid input")

	// ErrUsernameTaken reports that Register found the username already
	// registered (uniqueness is case-insensitive).
	ErrUsernameTaken = errors.New("username already taken")

	// ErrInvalidCredentials reports that Login could not verify the
	// supplied username/password pair.
	//
	// SECURITY: this is deliberately the only error Login can return for a
	// failed credential check, whether the user is unknown or the password
	// is wrong. Distinct errors would let an attacker enumerate registered
	// usernames (S9 security decision).
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidToken reports that a JWT failed verification: malformed
	// shape, malformed claims, wrong signature, or expired.
	//
	// SECURITY: no failure mode of verification is exposed separately; a
	// caller cannot tell an expired token from a tampered one.
	ErrInvalidToken = errors.New("invalid token")

	// ErrNoSecret reports that a token service was built without the
	// required HMAC secret.
	ErrNoSecret = errors.New("jwt secret must not be empty")

	// ErrShortSecret reports that the HMAC secret is below the minimum
	// length for HS256 (32 bytes).
	ErrShortSecret = errors.New("jwt secret must be at least 32 bytes")
)
