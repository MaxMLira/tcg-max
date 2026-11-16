package auth

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestPHCRoundTrip(t *testing.T) {
	salt := []byte("0123456789abcdef")
	key := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123")
	encoded := encodePHC(salt, key)

	parts, err := parsePHC(encoded)
	if err != nil {
		t.Fatalf("parsePHC(encodePHC): %v", err)
	}
	if parts.m != argon2idMemoryKiB || parts.t != argon2idIterations || parts.p != argon2idParallelism {
		t.Errorf("params = (m=%d, t=%d, p=%d), want (%d, %d, %d)",
			parts.m, parts.t, parts.p, argon2idMemoryKiB, argon2idIterations, argon2idParallelism)
	}
	if !bytes.Equal(parts.salt, salt) {
		t.Errorf("salt = %x, want %x", parts.salt, salt)
	}
	if !bytes.Equal(parts.key, key) {
		t.Errorf("key = %x, want %x", parts.key, key)
	}
	if got := encodePHC(parts.salt, parts.key); got != encoded {
		t.Errorf("encodePHC(parsePHC(x)) = %q, want %q", got, encoded)
	}
}

func TestParsePHCRejectsMalformed(t *testing.T) {
	valid := encodePHC([]byte("0123456789abcdef"), []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123"))
	mutate := func(section int, replacement string) string {
		parts := strings.Split(valid, "$")
		parts[section] = replacement
		return strings.Join(parts, "$")
	}

	cases := map[string]string{
		"not a phc string":    "not-a-phc-string",
		"wrong field count":   "$argon2id$v=19$m=65536,t=3,p=2$" + "salt",
		"missing leading $":   "argon2id$v=19$m=65536,t=3,p=2$salt$key",
		"wrong algorithm":     mutate(1, "argon2i"),
		"wrong version":       mutate(2, "v=16"),
		"params reordered":    mutate(3, "t=3,m=65536,p=2"),
		"param extra field":   mutate(3, "m=65536,t=3,p=2,x=1"),
		"param missing field": mutate(3, "m=65536,t=3"),
		"param non numeric":   mutate(3, "m=64k,t=3,p=2"),
		"param negative":      mutate(3, "m=-1,t=3,p=2"),
		"param plus sign":     mutate(3, "m=+65536,t=3,p=2"),
		"param zero":          mutate(3, "m=65536,t=0,p=2"),
		"salt with padding":   mutate(4, "AAAA="),
		"key std base64":      mutate(5, "a+b/"),
		"empty salt":          mutate(4, ""),
		"empty key":           mutate(5, ""),
	}
	for name, encoded := range cases {
		if _, err := parsePHC(encoded); !errors.Is(err, errPHC) {
			t.Errorf("%s: parsePHC = %v, want errPHC", name, err)
		}
	}
}
