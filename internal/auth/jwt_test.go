package auth

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// testSecret is exactly 32 bytes, the minimum HS256 secret length.
const testSecret = "0123456789abcdef0123456789abcdef"

func newTestTokens(t *testing.T, ttl time.Duration) *TokenService {
	t.Helper()
	ts, err := NewTokenService(testSecret, ttl)
	if err != nil {
		t.Fatalf("NewTokenService: %v", err)
	}
	return ts
}

func TestTokenServiceRejectsWeakSecret(t *testing.T) {
	if _, err := NewTokenService("", 0); !errors.Is(err, ErrNoSecret) {
		t.Errorf("empty secret: %v, want ErrNoSecret", err)
	}
	if _, err := NewTokenService("too-short", 0); !errors.Is(err, ErrShortSecret) {
		t.Errorf("short secret: %v, want ErrShortSecret", err)
	}
}

func TestIssueAndAuthenticate(t *testing.T) {
	ts := newTestTokens(t, 0)
	token, err := ts.Issue("user-1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	sub, err := ts.Authenticate(token)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if sub != "user-1" {
		t.Fatalf("subject = %q, want user-1", sub)
	}
}

func TestIssueRejectsEmptySubject(t *testing.T) {
	ts := newTestTokens(t, 0)
	if _, err := ts.Issue(""); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("Issue(\"\") = %v, want ErrInvalidInput", err)
	}
}

func TestAuthenticateRejectsTamperedToken(t *testing.T) {
	ts := newTestTokens(t, 0)
	token, err := ts.Issue("user-1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	parts := strings.Split(token, ".")

	// 1. Signature altered: flip the last base64url character.
	sig := []byte(parts[2])
	if sig[len(sig)-1] == 'A' {
		sig[len(sig)-1] = 'B'
	} else {
		sig[len(sig)-1] = 'A'
	}
	if _, err := ts.Authenticate(parts[0] + "." + parts[1] + "." + string(sig)); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("tampered signature: %v, want ErrInvalidToken", err)
	}

	// 2. Forged payload with the original signature.
	forgedPayload := jwtB64.EncodeToString([]byte(`{"sub":"user-2","iat":1,"exp":9999999999}`))
	if _, err := ts.Authenticate(parts[0] + "." + forgedPayload + "." + parts[2]); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("forged payload: %v, want ErrInvalidToken", err)
	}

	// 3. Token signed with a different secret.
	other, err := NewTokenService("fedcba9876543210fedcba9876543210", 0)
	if err != nil {
		t.Fatalf("NewTokenService(other): %v", err)
	}
	if _, err := other.Authenticate(token); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("wrong secret: %v, want ErrInvalidToken", err)
	}

	// 4. "alg":"none" attack: re-encoded header, empty signature.
	noneHeader := jwtB64.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	if _, err := ts.Authenticate(noneHeader + "." + parts[1] + "."); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("alg none: %v, want ErrInvalidToken", err)
	}
}

func TestAuthenticateRejectsExpiredToken(t *testing.T) {
	// A sub-second TTL expires within the same Unix second.
	ts := newTestTokens(t, time.Nanosecond)
	token, err := ts.Issue("user-1")
	if err != nil {
		t.Fatalf("Issue: %v", err)
	}
	if _, err := ts.Authenticate(token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expired: %v, want ErrInvalidToken", err)
	}
}

func TestAuthenticateRejectsMalformedToken(t *testing.T) {
	ts := newTestTokens(t, 0)

	shapes := map[string]string{
		"empty":         "",
		"one part":      "abc",
		"two parts":     "abc.def",
		"four parts":    "abc.def.ghi.jkl",
		"bad header":    "!not-base64.def.ghi",
		"bad signature": "abc.def.!not-base64",
	}
	for name, token := range shapes {
		if _, err := ts.Authenticate(token); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("%s: %v, want ErrInvalidToken", name, err)
		}
	}

	// Valid signature but a mandatory claim missing.
	noSub, err := ts.sign(jwtClaims{IssuedAt: 1, ExpiresAt: 9999999999})
	if err != nil {
		t.Fatalf("sign(no sub): %v", err)
	}
	noIat, err := ts.sign(jwtClaims{Subject: "x", ExpiresAt: 9999999999})
	if err != nil {
		t.Fatalf("sign(no iat): %v", err)
	}
	noExp, err := ts.sign(jwtClaims{Subject: "x", IssuedAt: 1})
	if err != nil {
		t.Fatalf("sign(no exp): %v", err)
	}
	for name, token := range map[string]string{
		"missing subject": noSub,
		"missing iat":     noIat,
		"missing exp":     noExp,
	} {
		if _, err := ts.Authenticate(token); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("%s: %v, want ErrInvalidToken", name, err)
		}
	}
}
