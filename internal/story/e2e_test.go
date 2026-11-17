package story_test

import (
	"errors"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/domain"
	"github.com/MaxMLira/tcg-max/internal/story"
	"github.com/MaxMLira/tcg-max/internal/story/seed"
)

// TestFullCampaignProgression plays the whole seeded campaign end to end:
// every duel is played (with a loss and a retry before the victory), every
// victory grants exactly the three predefined cards of that duel, and the
// campaign cannot be advanced after the final victory.
func TestFullCampaignProgression(t *testing.T) {
	c, err := seed.Campaign()
	if err != nil {
		t.Fatalf("seed.Campaign(): %v", err)
	}
	p := story.NewProgress()
	total := c.DuelCount()
	for i := 0; i < total; i++ {
		level, duel, ok := p.Current(c)
		if !ok {
			t.Fatalf("progression %d: no current duel", i)
		}
		// A loss never changes the state: the same duel can be retried.
		if _, err := p.RecordLoss(c, level.Number(), duel.Number()); err != nil {
			t.Fatalf("progression %d: RecordLoss: %v", i, err)
		}
		next, err := p.RecordVictory(c, level.Number(), duel.Number())
		if err != nil {
			t.Fatalf("progression %d: RecordVictory(%d,%d): %v", i, level.Number(), duel.Number(), err)
		}
		// The victory grants exactly the three predefined cards of the
		// duel that was just won (game-rules §58).
		got, ok := next.LastRewards(c)
		if !ok || len(got) != 3 {
			t.Fatalf("progression %d: LastRewards() = (%v, %v), want 3 cards", i, got, ok)
		}
		want := make([]domain.CardID, 0, 3)
		for _, card := range duel.Rewards() {
			want = append(want, card.ID())
		}
		for i, card := range got {
			if card.ID() != want[i] {
				t.Fatalf("progression %d: reward %d = %q, want %q", i, i, card.ID(), want[i])
			}
		}
		p = next
	}
	if !p.IsComplete(c) {
		t.Fatal("campaign not complete after 20 victories")
	}
	if _, _, ok := p.Current(c); ok {
		t.Fatal("current duel on a complete campaign, want none")
	}
	if _, err := p.RecordVictory(c, 1, 1); !errors.Is(err, story.ErrCampaignComplete) {
		t.Fatalf("RecordVictory after completion err = %v, want wrapped %v", err, story.ErrCampaignComplete)
	}
	for l := 1; l <= c.LevelCount(); l++ {
		for d := 1; d <= story.DuelsPerLevel; d++ {
			if !p.HasWon(c, l, d) {
				t.Fatalf("HasWon(%d,%d) = false on a complete campaign, want true", l, d)
			}
		}
	}
}

// TestStarterDeckIsAValidDeck checks the player's starting deck against the
// deck rules (game-rules §4).
func TestStarterDeckIsAValidDeck(t *testing.T) {
	deck, err := seed.StarterDeck()
	if err != nil {
		t.Fatalf("seed.StarterDeck(): %v", err)
	}
	if err := deck.Validate(); err != nil {
		t.Fatalf("starter deck invalid: %v", err)
	}
}
