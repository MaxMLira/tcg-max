// Package user holds the player account (architecture section 5.1: "User = Player").
//
// A User is the persistent identity that owns a Collection, Decks and story
// progress. This package owns the identity rules: opaque user ids, username
// validation and normalization, and the persistence contract (Repository).
// Password hashing, JWT issuance and the Register/Login use cases live in
// internal/auth; this package never sees a plaintext password.
package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// MinUsernameLength is the minimum number of characters of a valid username
// (S9 decision 6).
const MinUsernameLength = 3

// MaxUsernameLength is the maximum number of characters of a valid username
// (S9 decision 6).
const MaxUsernameLength = 32

// userIDRandomBytes is the number of random bytes an opaque UserID carries:
// 16 bytes = 32 hex characters (S9 decision 9).
const userIDRandomBytes = 16

// UserID is the opaque identity of a user account (S9 decision 9).
//
// It is 16 random bytes from crypto/rand rendered as 32 lowercase hex
// characters. Callers must treat it as an opaque string: it keys the user
// repository, future collections and story progress (S8/S10), and is the JWT
// "sub" claim (S9). It must not encode any business fact.
type UserID string

// NewUserID generates a new opaque user identity.
func NewUserID() (UserID, error) {
	buf := make([]byte, userIDRandomBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate user id: %w", err)
	}
	return UserID(hex.EncodeToString(buf)), nil
}

// NormalizeUsername returns the canonical form of a username: the same text in
// lowercase. Username uniqueness is case-insensitive (S9 decision 6) and this
// is the single place that canonicalizes: NewUser stores the canonical form
// and repositories index by it.
func NormalizeUsername(name string) string {
	return strings.ToLower(name)
}

// ValidateUsername reports whether username satisfies the account rules
// (S9 decision 6): 3 to 32 characters, each a letter, a digit or an
// underscore.
func ValidateUsername(username string) error {
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return fmt.Errorf("%w: must have between %d and %d characters",
			ErrInvalidUsername, MinUsernameLength, MaxUsernameLength)
	}
	for i := 0; i < len(username); i++ {
		c := username[i]
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
		case c == '_':
		default:
			return fmt.Errorf("%w: may only contain letters, digits or underscore",
				ErrInvalidUsername)
		}
	}
	return nil
}

// User is the player account (architecture section 5.1: User = Player).
//
// A User is the persistent identity that owns a Collection, Decks and story
// progress. It stores only the opaque password hash; the plaintext password
// is never kept (docs/mvp-story-mode.md section 1.6). The hash format is
// unknown to this package: user does not know which algorithm produced it
// (internal/auth does).
//
// Fields are private so a User can only be built through NewUser, which
// validates its invariants.
type User struct {
	id           UserID
	username     string
	passwordHash string
}

// NewUser validates and builds a user account.
//
// Enforced invariants:
//   - id must not be the zero value;
//   - username must satisfy ValidateUsername (S9 decision 6);
//   - passwordHash must not be empty.
//
// The username is stored canonicalized (lowercase), so uniqueness is
// case-insensitive (S9 decision 6).
func NewUser(id UserID, username, passwordHash string) (User, error) {
	if id == "" {
		return User{}, fmt.Errorf("%w: id must not be empty", ErrInvalidUser)
	}
	if err := ValidateUsername(username); err != nil {
		return User{}, err
	}
	if passwordHash == "" {
		return User{}, fmt.Errorf("%w: password hash must not be empty", ErrInvalidUser)
	}
	return User{id: id, username: NormalizeUsername(username), passwordHash: passwordHash}, nil
}

// ID returns the user identity.
func (u User) ID() UserID { return u.id }

// Username returns the canonical (lowercase) username.
func (u User) Username() string { return u.username }

// PasswordHash returns the opaque stored password hash. It must only be
// consumed by a password verifier (internal/auth); it is never the plaintext.
func (u User) PasswordHash() string { return u.passwordHash }
