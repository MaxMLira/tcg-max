package story

import (
	"errors"
	"testing"
)

func TestNewCampaign(t *testing.T) {
	t.Run("valid five levels", func(t *testing.T) {
		c := testCampaign(t, 5)
		if c.LevelCount() != 5 {
			t.Fatalf("LevelCount() = %d, want 5", c.LevelCount())
		}
		if c.DuelCount() != 20 {
			t.Fatalf("DuelCount() = %d, want 20", c.DuelCount())
		}
	})

	t.Run("valid single level", func(t *testing.T) {
		c := testCampaign(t, 1)
		if c.LevelCount() != 1 || c.DuelCount() != 4 {
			t.Fatalf("counts = (%d, %d), want (1, 4)", c.LevelCount(), c.DuelCount())
		}
	})

	t.Run("empty campaign rejected", func(t *testing.T) {
		_, err := NewCampaign()
		if !errors.Is(err, ErrInvalidCampaign) {
			t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidCampaign)
		}
	})

	t.Run("non sequential level numbers rejected", func(t *testing.T) {
		levels := testCampaign(t, 2).Levels()
		// Replace the second level (number 2) with level 3: position 2
		// expects number 2.
		third, err := NewLevel(3, testLevel(t, 3).duels)
		if err != nil {
			t.Fatalf("building level 3: %v", err)
		}
		levels[1] = third
		_, err = NewCampaign(levels...)
		if !errors.Is(err, ErrInvalidCampaign) {
			t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidCampaign)
		}
	})

	t.Run("more levels than the domain supports rejected", func(t *testing.T) {
		// Whitebox: NewLevel cannot build level 6, so the campaign invariant
		// is exercised directly with the zero-based internal representation.
		levels := make([]Level, 0, 6)
		for l := 1; l <= 5; l++ {
			levels = append(levels, testLevel(t, l))
		}
		levels = append(levels, Level{number: 6})
		_, err := NewCampaign(levels...)
		if !errors.Is(err, ErrInvalidCampaign) {
			t.Fatalf("err = %v, want wrapped %v", err, ErrInvalidCampaign)
		}
	})
}

func TestCampaignAccessors(t *testing.T) {
	c := testCampaign(t, 5)

	t.Run("Level", func(t *testing.T) {
		for n := 1; n <= 5; n++ {
			got, ok := c.Level(n)
			if !ok || got.Number() != n {
				t.Fatalf("Level(%d) = (%d, %v), want (%d, true)", n, got.Number(), ok, n)
			}
		}
		if _, ok := c.Level(6); ok {
			t.Fatal("Level(6) ok = true, want false")
		}
		if _, ok := c.Level(0); ok {
			t.Fatal("Level(0) ok = true, want false")
		}
	})

	t.Run("Duel", func(t *testing.T) {
		got, ok := c.Duel(3, 2)
		if !ok || got.Number() != 2 || got.Character().ID() != "char-3-2" {
			t.Fatalf("Duel(3,2) = (%v, %q, %v), want number 2 and character char-3-2", got.Number(), got.Character().ID(), ok)
		}
		if _, ok := c.Duel(9, 1); ok {
			t.Fatal("Duel(9,1) ok = true, want false")
		}
		if _, ok := c.Duel(1, 5); ok {
			t.Fatal("Duel(1,5) ok = true, want false")
		}
	})

	t.Run("Levels copy is isolated", func(t *testing.T) {
		levels := c.Levels()
		levels[0] = Level{}
		if c.Levels()[0].Number() != 1 {
			t.Fatal("mutating the returned Levels() slice changed the campaign")
		}
	})
}
