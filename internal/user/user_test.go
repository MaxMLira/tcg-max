package user

import (
	"errors"
	"regexp"
	"testing"
)

// validHash returns a syntactically valid placeholder hash for user tests.
// The hash format is opaque to the user package (internal/auth owns it).
func validHash(t *testing.T) string {
	t.Helper()
	return "$argon2id$16$3$2$abcdefghijklmnop$abcdefghijklmnopqrstuvwx"
}

// hexUserID is the rendered form of an opaque user id: 32 lowercase hex
// characters (S9 decision 9).
var hexUserID = regexp.MustCompile(`^[0-9a-f]{32}$`)

func TestNewUserID(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 256; i++ {
		id, err := NewUserID()
		if err != nil {
			t.Fatalf("NewUserID: %v", err)
		}
		if !hexUserID.MatchString(string(id)) {
			t.Fatalf("NewUserID = %q, want 32 lowercase hex chars", id)
		}
		if seen[string(id)] {
			t.Fatalf("duplicate id after %d draws: %q", i+1, id)
		}
		seen[string(id)] = true
	}
}

func TestNewUserIDFormat(t *testing.T) {
	// 16 random bytes rendered as hex (S9 decision 9).
	id, err := NewUserID()
	if err != nil {
		t.Fatalf("NewUserID: %v", err)
	}

	if len(string(id)) != 32 {
		t.Fatalf("len(id) = %d, want 32", len(string(id)))
	}
}

func TestValidateUsername(t *testing.T) {
	valid := []string{
		"abc", "ABC", "a_1", "user_name", "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
		"Alice_Smith_123",
	}
	for _, name := range valid {
		if err := ValidateUsername(name); err != nil {
			t.Errorf("ValidateUsername(%q) = %v, want nil", name, err)
		}
	}

	invalid := []string{
		"", "ab", // too short
		"abcdefghijklmnopqrstuvwxyzabcdefg", // 33 chars: too long
		"a b", "a-b", "a.b", "a/b", "usér", "user!", " user",
	}
	for _, name := range invalid {
		if err := ValidateUsername(name); !errors.Is(err, ErrInvalidUsername) {
			t.Errorf("ValidateUsername(%q) = %v, want ErrInvalidUsername", name, err)
		}
	}
}

func TestNormalizeUsername(t *testing.T) {
	if got := NormalizeUsername("Alice_Bob-1"); got != "alice_bob-1" {
		t.Fatalf("NormalizeUsername = %q, want %q", got, "alice_bob-1")
	}
}

func TestNewUser(t *testing.T) {
	id, err := NewUserID()
	if err != nil {
		t.Fatalf("NewUserID: %v", err)
	}

	t.Run("accepts valid user and canonicalizes username", func(t *testing.T) {
		u, err := NewUser(id, "Alice_Bob", validHash(t))
		if err != nil {
			t.Fatalf("NewUser: %v", err)
		}
		if u.ID() != id {
			t.Errorf("ID = %q, want %q", u.ID(), id)
		}

		if u.Username() != "alice_bob" {
			t.Errorf("Username = %q, want canonical lowercase", u.Username())
		}
		if u.PasswordHash() != validHash(t) {
			t.Errorf("PasswordHash = %q, want stored hash", u.PasswordHash())
		}
	})

	t.Run("rejects empty id", func(t *testing.T) {
		_, err := NewUser("", "alice", validHash(t))
		if !errors.Is(err, ErrInvalidUser) {
			t.Fatalf("NewUser = %v, want ErrInvalidUser", err)
		}
	})

	t.Run("rejects invalid username", func(t *testing.T) {
		if _, err := NewUser(id, "a!", validHash(t)); !errors.Is(err, ErrInvalidUsername) {
			t.Fatalf("NewUser = %v, want ErrInvalidUsername", err)
		}
		if _, err := NewUser(id, "ab", validHash(t)); !errors.Is(err, ErrInvalidUsername) {
			t.Fatalf("NewUser = %v, want ErrInvalidUsername", err)
		}
	})

	t.Run("rejects empty password hash", func(t *testing.T) {
		_, err := NewUser(id, "alice", "")
		if !errors.Is(err, ErrInvalidUser) {
			t.Fatalf("NewUser = %v, want ErrInvalidUser", err)
		}
	})
}
