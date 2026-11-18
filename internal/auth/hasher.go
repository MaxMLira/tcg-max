package auth

import "context"

// Hasher hashes and verifies passwords.
//
// The interface isolates the expensive key derivation from the use cases so
// that Service is testable with a fast fake, while production and the
// integration tests use the Argon2id adapter (NewArgon2idHasher).
//
// Implementations must:
//   - return an opaque, self-describing hash (the PHC string, see phc.go);
//   - be safe for concurrent use.
type Hasher interface {
	// Hash computes the password hash for password.
	Hash(ctx context.Context, password string) (string, error)

	// Verify reports whether password verifies against hash. hash is the
	// exact string previously returned by Hash. A non-nil error means the
	// password is not accepted; implementations must not let the error
	// distinguish a wrong password from a malformed hash.
	Verify(ctx context.Context, password, hash string) error
}
