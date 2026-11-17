package story

import (
	"errors"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// testLevel builds a valid level with the given number and fresh duels.
func testLevel(t *testing.T, number int) Level {
	t.Helper()
	duels := [DuelsPerLevel]Duel{}
	for d := 1; d <= DuelsPerLevel; d++ {
		duels[d-1] = testDuel(t, d, testCharacter(t,
			"char-lvl", "Character",
		))
	}
	level, err := NewLevel(number, duels)
	if err != nil {
		t.Fatalf("building test level %d: %v", number, err)
	}
	return level
}

func TestNewLevel(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		level := testLevel(t, 1)
		if level.Number() != 1 {
			t.Fatalf("Number() = %d, want 1", level.Number())
		}
		if len(level.Duels()) != DuelsPerLevel {
			t.Fatalf("len(Duels()) = %d, want %d", len(level.Duels()), DuelsPerLevel)
		}
	})

	t.Run("level 5 valid", func(t *testing.T) {
		if _, err := NewLevel(5, testLevel(t, 5).duels); err != nil {
			t.Fatalf("NewLevel(5) error: %v", err)
		}
	})

	t.Run("level 0 rejected by the domain rule", func(t *testing.T) {
		_, err := NewLevel(0, testLevel(t, 1).duels)
		if !errors.Is(err, ErrInvalidLevel) {
			t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidLevel)
		}
		if !errors.Is(err, domain.ErrInvalidLevel) {
			t.Fatalf("err = %v, want wrapped %v", err, domain.ErrInvalidLevel)
		}
	})

	t.Run("level above the supported range rejected", func(t *testing.T) {
		_, err := NewLevel(domain.MaxSupportedLevel+1, testLevel(t, 1).duels)
		if !errors.Is(err, ErrInvalidLevel) {
			t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidLevel)
		}
		if !errors.Is(err, domain.ErrInvalidLevel) {
			t.Fatalf("err = %v, want wrapped %v", err, domain.ErrInvalidLevel)
		}
	})

	t.Run("duel slot does not match its number", func(t *testing.T) {
		duels := testLevel(t, 1).duels
		duels[0] = testDuel(t, 2, testCharacter(t, "char-mis", "Mis"))
		_, err := NewLevel(1, duels)
		if !errors.Is(err, ErrInvalidLevel) {
			t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidLevel)
		}
	})

	t.Run("zero duel in a slot rejected", func(t *testing.T) {
		duels := testLevel(t, 1).duels
		duels[2] = Duel{}
		_, err := NewLevel(1, duels)
		if !errors.Is(err, ErrInvalidLevel) {
			t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidLevel)
		}
	})
}

func TestLevelInitialLifePoints(t *testing.T) {
	tests := []struct {
		level int
		want  int
	}{
		{level: 1, want: 4000},
		{level: 2, want: 4000},
		{level: 3, want: 4000},
		{level: 4, want: 4000},
		{level: 5, want: 8000},
	}
	for _, tt := range tests {
		t.Run("level "+itoa(tt.level), func(t *testing.T) {
			got, err := testLevel(t, tt.level).InitialLifePoints()
			if err != nil {
				t.Fatalf("InitialLifePoints() error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("InitialLifePoints() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestLevelDuelAccessor(t *testing.T) {
	level := testLevel(t, 3)
	for d := 1; d <= DuelsPerLevel; d++ {
		got, ok := level.Duel(d)
		if !ok || got.Number() != d {
			t.Fatalf("Duel(%d) = (%v, %v), want number %d", d, got.Number(), ok, d)
		}
	}
	if _, ok := level.Duel(0); ok {
		t.Fatal("Duel(0) ok = true, want false")
	}
	if _, ok := level.Duel(DuelsPerLevel + 1); ok {
		t.Fatal("Duel(5) ok = true, want false")
	}
}

// itoa avoids importing strconv for a single digit.
func itoa(n int) string {
	return string(rune('0' + n))
}
