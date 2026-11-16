package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/user"
)

// dummyPassword is hashed once at construction and re-verified on every
// failed login so that an unknown user and a wrong password take the same
// amount of time (S9 security decision: timing must not reveal which
// usernames exist).
const dummyPassword = "auth-timing-equalization"

// Service implements the S9 authentication use cases: Register, Login,
// Authenticate.
//
// It is transport-agnostic: it receives plain values and returns values or
// errors. Future HTTP/WebSocket layers wrap it (S12+).
type Service struct {
	repo      user.Repository
	hasher    Hasher
	tokens    *TokenService
	dummyHash string
}

// New builds a Service.
//
// The hasher is used to precompute the one-off dummy hash used for login
// timing equalization. The token service must already have been configured
// with its secret (NewTokenService enforces the secret rules).
func New(repo user.Repository, hasher Hasher, tokens *TokenService) (*Service, error) {
	if repo == nil {
		return nil, fmt.Errorf("%w: repository", ErrInvalidInput)
	}
	if hasher == nil {
		return nil, fmt.Errorf("%w: hasher", ErrInvalidInput)
	}
	if tokens == nil {
		return nil, fmt.Errorf("%w: token service", ErrInvalidInput)
	}
	dummyHash, err := hasher.Hash(context.Background(), dummyPassword)
	if err != nil {
		return nil, fmt.Errorf("auth: seed dummy hash: %w", err)
	}
	return &Service{repo: repo, hasher: hasher, tokens: tokens, dummyHash: dummyHash}, nil
}

// Register creates a new account.
//
// The username must satisfy the user package rules (3-32 characters,
// [A-Za-z0-9_]); uniqueness is case-insensitive. The password is hashed by
// the service's Hasher before anything touches the repository.
//
// Errors: ErrInvalidInput (invalid username, empty password),
// ErrUsernameTaken (canonical username already registered), or unexpected
// hashing/repository errors.
func (s *Service) Register(ctx context.Context, username, password string) (user.User, error) {
	if err := user.ValidateUsername(username); err != nil {
		return user.User{}, fmt.Errorf("%w: username: %v", ErrInvalidInput, err)
	}
	if password == "" {
		return user.User{}, fmt.Errorf("%w: password must not be empty", ErrInvalidInput)
	}
	hash, err := s.hasher.Hash(ctx, password)
	if err != nil {
		return user.User{}, fmt.Errorf("auth: hash password: %w", err)
	}
	id, err := user.NewUserID()
	if err != nil {
		return user.User{}, fmt.Errorf("auth: new user id: %w", err)
	}
	u, err := user.NewUser(id, username, hash)
	if err != nil {
		return user.User{}, fmt.Errorf("auth: new user: %w", err)
	}
	if err := s.repo.Create(ctx, u); err != nil {
		if errors.Is(err, user.ErrUsernameTaken) {
			return user.User{}, ErrUsernameTaken
		}
		return user.User{}, fmt.Errorf("auth: create user: %w", err)
	}
	return u, nil
}

// Login verifies username/password and returns a signed token for the
// account on success.
//
// SECURITY (S9): a failed login always returns ErrInvalidCredentials,
// whether the username is unknown or the password is wrong, and both paths
// perform one full KDF verification (unknown users verify against the
// precomputed dummy hash) so timing does not reveal which usernames exist.
func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	u, err := s.repo.ByUsername(ctx, username)
	if err != nil {
		if !errors.Is(err, user.ErrUserNotFound) {
			return "", fmt.Errorf("auth: look up user: %w", err)
		}
		// Unknown user: perform one KDF verification against the dummy
		// hash so timing matches the wrong-password path, then fail with
		// the same error. The dummy password is fixed and private, so a
		// match is a programming error at worst and still rejected.
		_ = s.hasher.Verify(ctx, password, s.dummyHash)
		return "", fmt.Errorf("%w: username or password", ErrInvalidCredentials)
	}
	if verr := s.hasher.Verify(ctx, password, u.PasswordHash()); verr != nil {
		return "", fmt.Errorf("%w: username or password", ErrInvalidCredentials)
	}
	token, err := s.tokens.Issue(string(u.ID()))
	if err != nil {
		return "", fmt.Errorf("auth: issue token: %w", err)
	}
	return token, nil
}

// Authenticate verifies a token and returns the account it was issued for.
//
// It returns ErrInvalidToken for any invalid token (bad shape, bad
// signature, malformed claims, expired), and ErrUserNotFound if the token
// is valid but the account no longer exists.
func (s *Service) Authenticate(ctx context.Context, token string) (user.User, error) {
	subject, err := s.tokens.Authenticate(token)
	if err != nil {
		return user.User{}, err
	}
	u, err := s.repo.ByID(ctx, user.UserID(subject))
	if err != nil {
		return user.User{}, fmt.Errorf("auth: load user: %w", err)
	}
	return u, nil
}
