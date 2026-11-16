// Package memory provides in-memory reference implementations of user
// repository contracts (architecture.md: infrastructure adapters).
//
// S9 ships a user repository so that the authentication use cases and their
// integration tests run without external services. The same contract is the
// one a production store (PostgreSQL) will implement.
package memory

import (
	"context"
	"sync"

	"github.com/MaxMLira/tcg-max/internal/user"
)

// UserRepository is an in-memory user.Repository.
//
// It is safe for concurrent use (S9 decision 5): all state is guarded by a
// single RWMutex, and Create performs the username-uniqueness check and the
// insert under the write lock, so duplicate detection is atomic.
//
// It stores users keyed by their canonical (lowercase) username; lookups by
// id scan the small map. This is a reference implementation, not an
// optimization target.
type UserRepository struct {
	mu         sync.RWMutex
	byID       map[user.UserID]user.User
	byUsername map[string]user.User
}

// NewUserRepository returns an empty in-memory UserRepository.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:       make(map[user.UserID]user.User),
		byUsername: make(map[string]user.User),
	}
}

// Create implements user.Repository. It returns user.ErrUsernameTaken when
// the canonical username is already taken; the check and the insert happen
// under the same write lock, so concurrent creates of the same username
// cannot both succeed.
func (r *UserRepository) Create(_ context.Context, u user.User) error {
	username := u.Username()
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, taken := r.byUsername[username]; taken {
		return user.ErrUsernameTaken
	}
	r.byUsername[username] = u
	r.byID[u.ID()] = u
	return nil
}

// ByID implements user.Repository.
func (r *UserRepository) ByID(_ context.Context, id user.UserID) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return user.User{}, user.ErrUserNotFound
	}
	return u, nil
}

// ByUsername implements user.Repository: it normalizes name and resolves
// the account, so lookups are case-insensitive (S9 decision 6).
func (r *UserRepository) ByUsername(_ context.Context, name string) (user.User, error) {
	username := user.NormalizeUsername(name)
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byUsername[username]
	if !ok {
		return user.User{}, user.ErrUserNotFound
	}
	return u, nil
}
