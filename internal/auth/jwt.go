package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// jwtAlgorithm is the only JWT algorithm this package signs and accepts.
// The header is pinned to HS256 on both ends, so a token cannot be
// re-keyed as "none" or as any other algorithm (S9 security decision).
const jwtAlgorithm = "HS256"

// jwtB64 is base64url without padding, per RFC 7515.
var jwtB64 = base64.RawURLEncoding

// minSecretBytes is the minimum length of the HMAC secret for HS256.
// 32 bytes matches the output size of SHA-256 (RFC 7518 section 3.2).
const minSecretBytes = 32

// defaultTokenTTL is the default lifetime of an issued token (S9: 8 hours).
const defaultTokenTTL = 8 * time.Hour

// jwtClaims are the JWT claims this package issues and accepts (S9):
// subject, issued-at and expiry. None of them is optional: a token missing
// any of them is invalid.
type jwtClaims struct {
	Subject   string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// TokenService issues and verifies HS256 JWTs using only the standard
// library (no third-party JWT dependency; S9 decision 7).
type TokenService struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenService builds a TokenService.
//
// The secret is mandatory and must be at least 32 bytes: HS256 signs with
// HMAC-SHA-256, so a short secret would be guessable. There is no built-in
// default on purpose: a forgotten secret must fail loudly at startup, not
// silently fall back to a known value.
//
// ttl sets the lifetime of issued tokens; a non-positive ttl selects the
// default of 8 hours.
func NewTokenService(secret string, ttl time.Duration) (*TokenService, error) {
	if secret == "" {
		return nil, ErrNoSecret
	}
	if len(secret) < minSecretBytes {
		return nil, fmt.Errorf("%w: got %d bytes, want at least %d", ErrShortSecret, len(secret), minSecretBytes)
	}
	if ttl <= 0 {
		ttl = defaultTokenTTL
	}
	return &TokenService{secret: []byte(secret), ttl: ttl}, nil
}

// Issue creates a signed token for userID that expires after the service's
// TTL. The token uses the RFC 7519 compact serialization with header
// {"alg":"HS256","typ":"JWT"} and claims sub/iat/exp.
func (s *TokenService) Issue(userID string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("%w: subject", ErrInvalidInput)
	}
	now := time.Now()
	claims := jwtClaims{
		Subject:   userID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(s.ttl).Unix(),
	}
	return s.sign(claims)
}

// Authenticate verifies the token signature and claims, and returns the
// subject (the user id) of a valid token. Every failure mode - malformed
// shape, bad signature, malformed claims, expired - returns an error
// wrapping ErrInvalidToken without distinguishing the cause (S9 security
// decision).
func (s *TokenService) Authenticate(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("%w: shape", ErrInvalidToken)
	}
	headerJSON, err := jwtB64.Strict().DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("%w: header", ErrInvalidToken)
	}
	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil || header.Alg != jwtAlgorithm {
		return "", fmt.Errorf("%w: header", ErrInvalidToken)
	}

	// Verify the signature over "header.payload" before reading any claim,
	// so an unauthenticated payload is never interpreted.
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(parts[0] + "." + parts[1]))
	signature, err := jwtB64.Strict().DecodeString(parts[2])
	if err != nil {
		return "", fmt.Errorf("%w: signature", ErrInvalidToken)
	}
	if !hmac.Equal(mac.Sum(nil), signature) {
		return "", fmt.Errorf("%w: signature", ErrInvalidToken)
	}

	payloadJSON, err := jwtB64.Strict().DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("%w: payload", ErrInvalidToken)
	}
	var claims jwtClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return "", fmt.Errorf("%w: claims", ErrInvalidToken)
	}
	if claims.Subject == "" || claims.IssuedAt == 0 || claims.ExpiresAt == 0 {
		return "", fmt.Errorf("%w: claims", ErrInvalidToken)
	}
	if time.Now().Unix() >= claims.ExpiresAt {
		return "", fmt.Errorf("%w: expired", ErrInvalidToken)
	}
	return claims.Subject, nil
}

// sign serializes the claims and signs them with HS256.
func (s *TokenService) sign(claims jwtClaims) (string, error) {
	header, err := json.Marshal(struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}{Alg: jwtAlgorithm, Typ: "JWT"})
	if err != nil {
		return "", fmt.Errorf("jwt: marshal header: %w", err)
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("jwt: marshal claims: %w", err)
	}
	signingInput := jwtB64.EncodeToString(header) + "." + jwtB64.EncodeToString(payload)
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(signingInput))
	return signingInput + "." + jwtB64.EncodeToString(mac.Sum(nil)), nil
}
