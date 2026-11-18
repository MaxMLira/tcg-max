package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MaxMLira/tcg-max/internal/repository/memory"
)

// newIntegrationService wires the real S9 components: in-memory repository,
// Argon2id hasher, HS256 token service. It is slower than the unit tests
// (real KDF) and exercises the security paths end to end.
func newIntegrationService(t *testing.T, ttl time.Duration) *Service {
	t.Helper()
	repo := memory.NewUserRepository()
	tokens, err := NewTokenService("integration-test-secret-32-bytes!", ttl)
	if err != nil {
		t.Fatalf("NewTokenService: %v", err)
	}
	svc, err := New(repo, NewArgon2idHasher(), tokens)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return svc
}

func TestIntegrationRegisterLoginAuthenticate(t *testing.T) {
	svc := newIntegrationService(t, 0)
	ctx := context.Background()

	u, err := svc.Register(ctx, "Integration", "correct horse battery staple")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	token, err := svc.Login(ctx, "integration", "correct horse battery staple")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	got, err := svc.Authenticate(ctx, token)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if got.ID() != u.ID() {
		t.Fatalf("Authenticate.ID = %q, want %q", got.ID(), u.ID())
	}
}

func TestIntegrationLoginSecurity(t *testing.T) {
	svc := newIntegrationService(t, 0)
	ctx := context.Background()

	if _, err := svc.Register(ctx, "Secure", "correct horse"); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if _, err := svc.Login(ctx, "secure", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login(wrong) = %v, want ErrInvalidCredentials", err)
	}
	if _, err := svc.Login(ctx, "ghost", "wrong"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login(unknown) = %v, want ErrInvalidCredentials", err)
	}
}

func TestIntegrationTokenTTL(t *testing.T) {
	// A sub-second TTL expires within the same Unix second.
	svc := newIntegrationService(t, time.Nanosecond)
	ctx := context.Background()

	if _, err := svc.Register(ctx, "Short", "pw-pw-pw-pw"); err != nil {
		t.Fatalf("Register: %v", err)
	}
	token, err := svc.Login(ctx, "short", "pw-pw-pw-pw")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if _, err := svc.Authenticate(ctx, token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Authenticate(expired) = %v, want ErrInvalidToken", err)
	}
}
