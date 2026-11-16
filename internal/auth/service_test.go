package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/repository/memory"
)

// fakeHasher is a fast in-memory Hasher for Service unit tests. It is not
// secure and must never be used outside tests.
type fakeHasher struct {
	hashes map[string]string
	taken  int
}

func newFakeHasher() *fakeHasher {
	return &fakeHasher{hashes: map[string]string{}}
}

func (f *fakeHasher) Hash(_ context.Context, password string) (string, error) {
	hash := "hash-of:" + password
	f.hashes[password] = hash
	return hash, nil
}

func (f *fakeHasher) Verify(_ context.Context, password, hash string) error {
	f.taken++
	if f.hashes[password] != hash {
		return errors.New("mismatch")
	}
	return nil
}

func newTestService(t *testing.T) (*Service, *fakeHasher) {
	t.Helper()
	hasher := newFakeHasher()
	tokens := newTestTokens(t, 0)
	svc, err := New(memory.NewUserRepository(), hasher, tokens)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return svc, hasher
}

func TestRegisterAndLogin(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	u, err := svc.Register(ctx, "Alice", "S3cret-pass")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if u.Username() != "alice" {
		t.Fatalf("username = %q, want canonical lowercase", u.Username())
	}
	if u.PasswordHash() == "S3cret-pass" || u.PasswordHash() == "" {
		t.Fatalf("password hash = %q, want opaque non-empty hash", u.PasswordHash())
	}

	token, err := svc.Login(ctx, "Alice", "S3cret-pass")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if token == "" {
		t.Fatal("Login returned an empty token")
	}

	again, err := svc.Authenticate(ctx, token)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if again.ID() != u.ID() {
		t.Fatalf("Authenticate.ID = %q, want %q", again.ID(), u.ID())
	}
}

func TestRegisterValidation(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	cases := []struct {
		name, username, password string
	}{
		{"empty username", "", "pw"},
		{"short username", "ab", "pw"},
		{"long username", "abcdefghijklmnopqrstuvwxyzabcdefg", "pw"},
		{"bad chars", "al ice", "pw"},
		{"empty password", "alice", ""},
	}
	for _, tc := range cases {
		if _, err := svc.Register(ctx, tc.username, tc.password); !errors.Is(err, ErrInvalidInput) {
			t.Errorf("%s: Register = %v, want ErrInvalidInput", tc.name, err)
		}
	}
}

func TestRegisterUsernameTaken(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	if _, err := svc.Register(ctx, "Alice", "pw-1"); err != nil {
		t.Fatalf("Register(first): %v", err)
	}
	// Case-insensitive uniqueness: ALICE is the same account name.
	if _, err := svc.Register(ctx, "ALICE", "pw-2"); !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("Register(dup) = %v, want ErrUsernameTaken", err)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	if _, err := svc.Register(ctx, "Alice", "correct"); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if _, err := svc.Login(ctx, "alice", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login(wrong password) = %v, want ErrInvalidCredentials", err)
	}
}

func TestLoginUnknownUserSameErrorAndWork(t *testing.T) {
	svc, hasher := newTestService(t)
	ctx := context.Background()

	if _, err := svc.Register(ctx, "Alice", "correct"); err != nil {
		t.Fatalf("Register: %v", err)
	}

	// Baseline: a wrong password for a known user performs one verify.
	if _, err := svc.Login(ctx, "alice", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login(wrong) = %v, want ErrInvalidCredentials", err)
	}
	knownWork := hasher.taken

	// Unknown user: same error...
	hasher.taken = 0
	if _, err := svc.Login(ctx, "ghost", "anything"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login(unknown) = %v, want ErrInvalidCredentials", err)
	}
	// ...and the same KDF work (timing equalization).
	if hasher.taken != knownWork {
		t.Fatalf("verify work: unknown user = %d, known user = %d; must be equal",
			hasher.taken, knownWork)
	}
}

func TestAuthenticateInvalidToken(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	for name, token := range map[string]string{
		"empty":   "",
		"garbage": "not.a.jwt",
	} {
		if _, err := svc.Authenticate(ctx, token); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("%s: Authenticate = %v, want ErrInvalidToken", name, err)
		}
	}
}

func TestNewRejectsNilDependencies(t *testing.T) {
	tokens := newTestTokens(t, 0)
	if _, err := New(nil, newFakeHasher(), tokens); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("nil repo: %v, want ErrInvalidInput", err)
	}
	if _, err := New(memory.NewUserRepository(), nil, tokens); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("nil hasher: %v, want ErrInvalidInput", err)
	}
	if _, err := New(memory.NewUserRepository(), newFakeHasher(), nil); !errors.Is(err, ErrInvalidInput) {
		t.Errorf("nil tokens: %v, want ErrInvalidInput", err)
	}
}
