package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestArgon2idHashRoundTrip(t *testing.T) {
	h := NewArgon2idHasher()
	ctx := context.Background()

	hash, err := h.Hash(ctx, "correct horse battery staple")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if !strings.HasPrefix(hash, "$argon2id$v=19$m=65536,t=3,p=2$") {
		t.Fatalf("Hash prefix = %q, want the S9 parameter prefix", hash)
	}
	if err := h.Verify(ctx, "correct horse battery staple", hash); err != nil {
		t.Fatalf("Verify(correct): %v", err)
	}
	if err := h.Verify(ctx, "wrong password", hash); err == nil {
		t.Fatal("Verify(wrong) = nil, want an error")
	}
}

func TestArgon2idHashUsesRandomSalt(t *testing.T) {
	h := NewArgon2idHasher()
	ctx := context.Background()

	a, err := h.Hash(ctx, "same-password")
	if err != nil {
		t.Fatalf("Hash(a): %v", err)
	}
	b, err := h.Hash(ctx, "same-password")
	if err != nil {
		t.Fatalf("Hash(b): %v", err)
	}
	if a == b {
		t.Fatal("two Hash calls produced identical PHC strings; the salt must be random")
	}
}

func TestArgon2idVerifyRejectsMalformedHash(t *testing.T) {
	h := NewArgon2idHasher()
	ctx := context.Background()

	cases := map[string]string{
		"empty":           "",
		"garbage":         "not-a-hash",
		"other algorithm": "$argon2i$v=19$m=65536,t=3,p=2$salt$key",
		"truncated":       "$argon2id$v=19$m=65536,t=3,p=2$salt",
	}
	for name, hash := range cases {
		if err := h.Verify(ctx, "password", hash); !errors.Is(err, errPHC) {
			t.Errorf("%s: Verify = %v, want errPHC", name, err)
		}
	}
}
