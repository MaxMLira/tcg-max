package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/user"
)

// makeUser builds a valid user account for repository tests.
func makeUser(t *testing.T, name string) (user.User, error) {
	t.Helper()
	id, err := user.NewUserID()
	if err != nil {
		return user.User{}, fmt.Errorf("NewUserID: %w", err)
	}
	return user.NewUser(id, name, "test-hash")
}

func TestUserRepositoryCreateAndByID(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	u, err := makeUser(t, "Alice")
	if err != nil {
		t.Fatalf("makeUser: %v", err)
	}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.ByID(ctx, u.ID())
	if err != nil {
		t.Fatalf("ByID: %v", err)
	}
	if got.ID() != u.ID() || got.Username() != "alice" {
		t.Fatalf("ByID = (%q, %q), want (%q, %q)", got.ID(), got.Username(), u.ID(), "alice")
	}

	if got.PasswordHash() != u.PasswordHash() {
		t.Fatalf("ByID.PasswordHash = %q, want %q", got.PasswordHash(), u.PasswordHash())
	}

	unknown := user.UserID("00000000000000000000000000000000")
	if _, err := repo.ByID(ctx, unknown); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("ByID(unknown) = %v, want ErrUserNotFound", err)
	}
}

func TestUserRepositoryByUsernameIsCaseInsensitive(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	u, err := makeUser(t, "Alice_Bob")
	if err != nil {
		t.Fatalf("makeUser: %v", err)
	}
	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	for _, lookup := range []string{"alice_bob", "ALICE_BOB", "Alice_Bob", "aLiCe_bOb"} {
		got, err := repo.ByUsername(ctx, lookup)
		if err != nil {
			t.Fatalf("ByUsername(%q): %v", lookup, err)
		}
		if got.ID() != u.ID() {
			t.Fatalf("ByUsername(%q).ID = %q, want %q", lookup, got.ID(), u.ID())
		}
	}

	if _, err := repo.ByUsername(ctx, "nobody"); !errors.Is(err, user.ErrUserNotFound) {
		t.Fatalf("ByUsername(nobody) = %v, want ErrUserNotFound", err)
	}
}

func TestUserRepositoryCreateRejectsDuplicateUsername(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	first, err := makeUser(t, "Alice")
	if err != nil {
		t.Fatalf("makeUser: %v", err)
	}
	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("Create(first): %v", err)
	}

	// Same canonical username with different casing: rejected (S9 decision 6).
	dup, err := makeUser(t, "ALICE")
	if err != nil {
		t.Fatalf("makeUser: %v", err)
	}
	if err := repo.Create(ctx, dup); !errors.Is(err, user.ErrUsernameTaken) {
		t.Fatalf("Create(dup) = %v, want ErrUsernameTaken", err)
	}

	// The original account is untouched.
	got, err := repo.ByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("ByUsername: %v", err)
	}
	if got.ID() != first.ID() {
		t.Fatalf("ByUsername.ID = %q, want %q", got.ID(), first.ID())
	}
}

func TestUserRepositoryConcurrentCreate(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	const workers = 16
	users := make([]user.User, workers)
	for i := range users {
		var err error
		users[i], err = makeUser(t, "Racer")
		if err != nil {
			t.Fatalf("makeUser: %v", err)
		}
	}

	var wg sync.WaitGroup
	results := make(chan error, workers)
	for i := range users {
		wg.Add(1)
		go func(u user.User) {
			defer wg.Done()
			results <- repo.Create(ctx, u)
		}(users[i])
	}
	wg.Wait()
	close(results)

	created, taken := 0, 0
	for err := range results {
		switch {
		case err == nil:
			created++
		case errors.Is(err, user.ErrUsernameTaken):
			taken++
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if created != 1 || taken != workers-1 {
		t.Fatalf("created = %d, taken = %d; want exactly one created", created, taken)
	}

	// Exactly one stored account, resolvable by username.
	got, err := repo.ByUsername(ctx, "racer")
	if err != nil {
		t.Fatalf("ByUsername: %v", err)
	}
	found := false
	for i := range users {
		if users[i].ID() == got.ID() {
			found = true
		}
	}
	if !found {
		t.Fatalf("stored id %q is not one of the contenders", got.ID())
	}
}

func TestUserRepositoryConcurrentReadWrite(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	const seeds = 8
	users := make([]user.User, seeds)
	for i := range users {
		var err error
		users[i], err = makeUser(t, fmt.Sprintf("seed_%d", i))
		if err != nil {
			t.Fatalf("makeUser: %v", err)
		}
		if err := repo.Create(ctx, users[i]); err != nil {
			t.Fatalf("Create(seed %d): %v", i, err)
		}
	}

	var wg sync.WaitGroup
	errs := make(chan error, 2*seeds*4)
	for i := range users {
		wg.Add(1)
		go func(u user.User) {
			defer wg.Done()
			for j := 0; j < 4; j++ {
				if _, err := repo.ByID(ctx, u.ID()); err != nil {
					errs <- fmt.Errorf("ByID: %w", err)
					return
				}
				if _, err := repo.ByUsername(ctx, u.Username()); err != nil {
					errs <- fmt.Errorf("ByUsername: %w", err)
					return
				}
			}
		}(users[i])
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent read: %v", err)
	}
}
