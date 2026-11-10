package domain

import (
	"errors"
	"testing"
)

func TestNewMonsterSubtype_TrimsSurroundingWhitespace(t *testing.T) {
	// Arrange

	// Act
	subtype, err := NewMonsterSubtype("  Dragon \t")

	// Assert
	if err != nil {
		t.Fatalf("NewMonsterSubtype returned unexpected error: %v", err)
	}
	if subtype.String() != "Dragon" {
		t.Errorf("String() = %q, want %q", subtype.String(), "Dragon")
	}
}

func TestNewMonsterSubtype_KeepsTheGivenCase(t *testing.T) {
	// Arrange

	// Act
	subtype, err := NewMonsterSubtype("dRaGoN")

	// Assert
	if err != nil {
		t.Fatalf("NewMonsterSubtype returned unexpected error: %v", err)
	}
	if subtype.String() != "dRaGoN" {
		t.Errorf("String() = %q, want %q", subtype.String(), "dRaGoN")
	}
}

func TestNewMonsterSubtype_RejectsBlankName(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty", input: ""},
		{name: "spaces", input: "   "},
		{name: "tab", input: "\t"},
		{name: "newline and space", input: "\n "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			subtype, err := NewMonsterSubtype(tc.input)

			// Assert
			if !errors.Is(err, ErrInvalidSubtype) {
				t.Fatalf("error = %v, want ErrInvalidSubtype", err)
			}
			if !subtype.IsZero() {
				t.Errorf("subtype = %q, want the zero value", subtype)
			}
		})
	}
}

func TestMonsterSubtype_EqualsIgnoresCaseAndSurroundingSpaces(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "lower case", input: "dragon"},
		{name: "upper case with spaces", input: " DRAGON "},
		{name: "same value", input: "Dragon"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			other, err := NewMonsterSubtype(tc.input)
			if err != nil {
				t.Fatalf("NewMonsterSubtype(%q) failed: %v", tc.input, err)
			}

			// Act
			got := SubtypeDragon().Equals(other)
			gotReverse := other.Equals(SubtypeDragon())

			// Assert
			if !got || !gotReverse {
				t.Errorf("Dragon.Equals(%q) = %v, reverse = %v, want true both ways", tc.input, got, gotReverse)
			}
		})
	}
}

func TestMonsterSubtype_EqualsIsFalseForDifferentSubtypes(t *testing.T) {
	// Arrange
	dragon := SubtypeDragon()
	warrior := SubtypeWarrior()

	// Act
	got := dragon.Equals(warrior)

	// Assert
	if got {
		t.Error("Dragon.Equals(Warrior) = true, want false")
	}
}

func TestMonsterSubtype_EqualsHandlesTheZeroValue(t *testing.T) {
	tests := []struct {
		name string
		a, b MonsterSubtype
		want bool
	}{
		{name: "two zero values", a: MonsterSubtype{}, b: MonsterSubtype{}, want: true},
		{name: "zero and non-zero", a: MonsterSubtype{}, b: SubtypeDragon(), want: false},
		{name: "non-zero and zero", a: SubtypeDragon(), b: MonsterSubtype{}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			got := tc.a.Equals(tc.b)

			// Assert
			if got != tc.want {
				t.Errorf("Equals = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewMonsterSubtype_AcceptsNewSubtypeWithoutDomainChanges(t *testing.T) {
	// Arrange

	// Act
	beast, err := NewMonsterSubtype("Beast")

	// Assert
	if err != nil {
		t.Fatalf("NewMonsterSubtype returned unexpected error: %v", err)
	}
	if beast.IsZero() {
		t.Error("IsZero() = true, want false")
	}
	if beast.String() != "Beast" {
		t.Errorf("String() = %q, want %q", beast.String(), "Beast")
	}
	if !beast.Equals(beast) {
		t.Error("Beast.Equals(Beast) = false, want true")
	}
	for _, mvp := range []MonsterSubtype{SubtypeDragon(), SubtypeWarrior(), SubtypeMage(), SubtypeAquatic()} {
		if beast.Equals(mvp) {
			t.Errorf("Beast.Equals(%v) = true, want false", mvp)
		}
	}
}

func TestMvpSubtypes_AreNamedNonZeroAndDistinct(t *testing.T) {
	// Arrange
	subtypes := []struct {
		subtype MonsterSubtype
		want    string
	}{
		{subtype: SubtypeDragon(), want: "Dragon"},
		{subtype: SubtypeWarrior(), want: "Warrior"},
		{subtype: SubtypeMage(), want: "Mage"},
		{subtype: SubtypeAquatic(), want: "Aquatic"},
	}

	for i, a := range subtypes {
		// Act & Assert
		if a.subtype.IsZero() {
			t.Errorf("%s: IsZero() = true, want false", a.want)
		}
		if a.subtype.String() != a.want {
			t.Errorf("String() = %q, want %q", a.subtype.String(), a.want)
		}
		for j, b := range subtypes {
			if got := a.subtype.Equals(b.subtype); got != (i == j) {
				t.Errorf("%s.Equals(%s) = %v, want %v", a.want, b.want, got, i == j)
			}
		}
	}
}

func TestMonsterSubtype_IsZero(t *testing.T) {
	tests := []struct {
		name    string
		subtype MonsterSubtype
		want    bool
	}{
		{name: "zero value", subtype: MonsterSubtype{}, want: true},
		{name: "built subtype", subtype: SubtypeMage(), want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			got := tc.subtype.IsZero()

			// Assert
			if got != tc.want {
				t.Errorf("IsZero() = %v, want %v", got, tc.want)
			}
		})
	}
}
