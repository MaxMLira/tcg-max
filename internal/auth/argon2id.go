package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// errArgon2idMismatch is the uniform verification failure for a password
// that does not match the hash.
var errArgon2idMismatch = errors.New("argon2id: password mismatch")

// Argon2idHasher is the production Hasher: Argon2id in the PHC string
// format, parameters fixed by S9 decision 8.
//
// It is safe for concurrent use: it holds no mutable state.
type Argon2idHasher struct{}

// NewArgon2idHasher returns the Argon2id Hasher used by production and by
// the integration tests.
func NewArgon2idHasher() *Argon2idHasher {
	return &Argon2idHasher{}
}

// Hash computes the Argon2id PHC string for password using a fresh
// 16-byte random salt and the fixed S9 parameters.
func (h *Argon2idHasher) Hash(_ context.Context, password string) (string, error) {
	salt := make([]byte, argon2idSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("argon2id: generate salt: %w", err)
	}
	key := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(argon2idIterations),
		uint32(argon2idMemoryKiB),
		uint8(argon2idParallelism),
		argon2idKeyBytes,
	)
	return encodePHC(salt, key), nil
}

// Verify recomputes the key from the password and the parameters stored in
// hash, and compares the result in constant time.
//
// Any malformed hash or password mismatch yields a non-nil error, and the
// error is uniform on purpose: a caller must not be able to distinguish a
// wrong password from a malformed hash (S9 security decision). Corrupted
// stored parameters cannot cause a panic or an accept; they only produce a
// mismatch.
func (h *Argon2idHasher) Verify(_ context.Context, password, hash string) error {
	parts, err := parsePHC(hash)
	if err != nil {
		return err
	}
	key := argon2.IDKey(
		[]byte(password),
		parts.salt,
		uint32(parts.t),
		uint32(parts.m),
		uint8(parts.p),
		uint32(len(parts.key)),
	)
	if subtle.ConstantTimeCompare(key, parts.key) != 1 {
		return errArgon2idMismatch
	}
	return nil
}
