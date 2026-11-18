package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// errPHC marks a PHC string that is not in the exact format this package
// produces and accepts. It is internal: the Argon2id adapter converts it
// into its own verification error before it reaches a use case.
var errPHC = errors.New("phc: malformed hash")

// phcVersion is the Argon2 version field (19 = Argon2 1.3). Only this
// version is accepted (S9: strict PHC format).
const phcVersion = "19"

// argon2idPHCPrefix identifies the argon2id algorithm in a PHC string.
const argon2idPHCPrefix = "argon2id"

// phcBase64 is the PHC string alphabet: base64url without padding
// (A-Z a-z 0-9 - _), as defined for the PHC string format.
var phcBase64 = base64.URLEncoding.WithPadding(base64.NoPadding)

// Fixed S9 argon2id parameters (S9 decision 8): 64 MiB memory, 3
// iterations, parallelism 2, 16-byte salt, 32-byte key. They are part of
// the PHC format this package produces and accepts.
const (
	argon2idMemoryKiB   = 65536
	argon2idIterations  = 3
	argon2idParallelism = 2
	argon2idSaltBytes   = 16
	argon2idKeyBytes    = 32
)

// phcParts is the parsed form of a PHC string.
type phcParts struct {
	m, t, p int
	salt    []byte
	key     []byte
}

// encodePHC renders the PHC string for salt and key with the fixed S9
// argon2id parameters. It is the inverse of parsePHC: encoding a parsed
// value and re-parsing it yields the same value.
func encodePHC(salt, key []byte) string {
	return "$" + argon2idPHCPrefix + "$v=" + phcVersion +
		"$m=" + strconv.Itoa(argon2idMemoryKiB) +
		",t=" + strconv.Itoa(argon2idIterations) +
		",p=" + strconv.Itoa(argon2idParallelism) +
		"$" + phcBase64.EncodeToString(salt) +
		"$" + phcBase64.EncodeToString(key)
}

// parsePHC parses a PHC string of the exact shape this package produces:
//
// $argon2id$v=19$m=<int>,t=<int>,p=<int>$<salt>$<key>
//
// Every component is validated strictly:
//   - exactly five '$'-separated fields: empty head, algorithm, version,
//     parameters, salt, key;
//   - algorithm "argon2id" and version "19" only;
//   - parameter fields exactly "m", "t", "p" in that order, each a
//     positive decimal integer without sign;
//   - salt and key in PHC base64 (url alphabet, no padding), non-empty.
//
// Any deviation yields an error wrapping errPHC.
func parsePHC(encoded string) (phcParts, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return phcParts{}, fmt.Errorf("%w: expected 5 fields, got %d", errPHC, len(parts)-1)
	}
	if parts[0] != "" || parts[1] != argon2idPHCPrefix {
		return phcParts{}, fmt.Errorf("%w: unknown algorithm", errPHC)
	}
	if parts[2] != "v="+phcVersion {
		return phcParts{}, fmt.Errorf("%w: unsupported version %q", errPHC, parts[2])
	}
	m, t, p, err := parsePHCParams(parts[3])
	if err != nil {
		return phcParts{}, err
	}
	salt, err := decodePHCField(parts[4])
	if err != nil {
		return phcParts{}, err
	}
	key, err := decodePHCField(parts[5])
	if err != nil {
		return phcParts{}, err
	}
	return phcParts{m: m, t: t, p: p, salt: salt, key: key}, nil
}

// parsePHCParams parses the parameter section "m=<int>,t=<int>,p=<int>"
// with the fixed field order mandated by the PHC specification for argon2.
func parsePHCParams(section string) (int, int, int, error) {
	names := [3]string{"m", "t", "p"}
	fields := strings.Split(section, ",")
	if len(fields) != 3 {
		return 0, 0, 0, fmt.Errorf("%w: parameter section", errPHC)
	}
	values := [3]int{}
	for i, field := range fields {
		key, value, ok := strings.Cut(field, "=")
		if !ok || key != names[i] {
			return 0, 0, 0, fmt.Errorf("%w: parameter %q", errPHC, field)
		}
		n, err := parsePHCInt(value)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("%w: parameter %s", errPHC, names[i])
		}
		values[i] = n
	}
	return values[0], values[1], values[2], nil
}

// parsePHCInt parses a positive decimal integer without sign, padding
// whitespace or non-digit characters.
func parsePHCInt(value string) (int, error) {
	if value == "" {
		return 0, fmt.Errorf("%w: empty value", errPHC)
	}
	for i := 0; i < len(value); i++ {
		if c := value[i]; c < '0' || c > '9' {
			return 0, fmt.Errorf("%w: non-numeric value %q", errPHC, value)
		}
	}
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("%w: non-positive value %q", errPHC, value)
	}
	return n, nil
}

// decodePHCField decodes a non-empty PHC base64 field (url alphabet, no
// padding). Strict decoding rejects any character outside the alphabet,
// including the standard-base64 characters '+', '/' and the padding '='.
func decodePHCField(field string) ([]byte, error) {
	if field == "" {
		return nil, fmt.Errorf("%w: empty field", errPHC)
	}
	decoded, err := phcBase64.Strict().DecodeString(field)
	if err != nil {
		return nil, fmt.Errorf("%w: field encoding", errPHC)
	}
	return decoded, nil
}
