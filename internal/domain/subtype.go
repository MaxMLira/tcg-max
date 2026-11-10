package domain

import (
	"fmt"
	"strings"
)

// MonsterSubtype is the subtype of a Monster, for example Dragon or Warrior
// (game-rules §10.1.1, docs/mvp-story-mode.md §1.3 and assumption A3).
//
// It is a validated, non-blank value type rather than a closed enumeration, so
// new subtypes can be added without changing the domain model: any non-blank
// name built with NewMonsterSubtype is a valid subtype. The four subtypes of
// the Story Mode MVP (Dragon, Warrior, Mage and Aquatic) are plain values of
// this type, provided by SubtypeDragon, SubtypeWarrior, SubtypeMage and
// SubtypeAquatic.
//
// The zero value means "no subtype" (the case of Magic and Trap cards). Names
// are stored trimmed, keeping the case they were given, and are compared
// case-insensitively. Compare subtypes with Equals, never with ==: == is
// case-sensitive on the stored name.
type MonsterSubtype struct {
	name string
}

// NewMonsterSubtype validates and builds a Monster subtype. Leading and
// trailing whitespace is removed from the name. A name that is blank after
// trimming is rejected with ErrInvalidSubtype and the zero value is returned.
func NewMonsterSubtype(name string) (MonsterSubtype, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return MonsterSubtype{}, fmt.Errorf("%w: subtype name must not be blank", ErrInvalidSubtype)
	}
	return MonsterSubtype{name: trimmed}, nil
}

// String returns the subtype name, trimmed and in the case it was given.
func (s MonsterSubtype) String() string { return s.name }

// IsZero reports whether s is the zero value, meaning "no subtype".
func (s MonsterSubtype) IsZero() bool { return s.name == "" }

// Equals reports whether s and other are the same subtype.
//
// Normalization: NewMonsterSubtype already removes leading and trailing
// whitespace, so every non-zero value is stored trimmed. Equals does not trim
// again and does not touch whitespace inside the name ("Sea Dragon" and
// "Sea  Dragon" are different subtypes).
//
// Comparison: names are compared case-insensitively with strings.EqualFold, so
// "Dragon", "dragon" and "DRAGON" are the same subtype. This is simple case
// folding only; no special Unicode handling is applied (for example "Straße"
// and "STRASSE" are different subtypes). Two zero values are equal; a zero
// value is never equal to a non-zero one.
//
// Use Equals instead of ==, which compares the stored name case-sensitively.
func (s MonsterSubtype) Equals(other MonsterSubtype) bool {
	return strings.EqualFold(s.name, other.name)
}

// mustSubtype builds one of the provided MVP subtypes from a known non-blank
// name.
func mustSubtype(name string) MonsterSubtype {
	subtype, err := NewMonsterSubtype(name)
	if err != nil {
		panic(err)
	}
	return subtype
}

// SubtypeDragon returns the Dragon subtype (Story Mode MVP, game-rules §10.1.1).
func SubtypeDragon() MonsterSubtype { return mustSubtype("Dragon") }

// SubtypeWarrior returns the Warrior subtype (Story Mode MVP, game-rules §10.1.1).
func SubtypeWarrior() MonsterSubtype { return mustSubtype("Warrior") }

// SubtypeMage returns the Mage subtype (Story Mode MVP, game-rules §10.1.1).
func SubtypeMage() MonsterSubtype { return mustSubtype("Mage") }

// SubtypeAquatic returns the Aquatic subtype (Story Mode MVP, game-rules §10.1.1).
func SubtypeAquatic() MonsterSubtype { return mustSubtype("Aquatic") }
