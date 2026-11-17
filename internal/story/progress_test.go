package story

import (
	"errors"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// winFrom plays victories until the progress reaches Count() == want,
// always playing the current duel.
func winFrom(t *testing.T, c Campaign, p Progress, want int) Progress {
	t.Helper()
	for p.Count() < want {
		level, duel, ok := p.Current(c)
		if !ok {
			t.Fatalf("progress %d: no current duel", p.Count())
		}
		var err error
		p, err = p.RecordVictory(c, level.Number(), duel.Number())
		if err != nil {
			t.Fatalf("progress %d: RecordVictory(%d,%d): %v", p.Count(), level.Number(), duel.Number(), err)
		}
	}
	return p
}

func TestProgressFresh(t *testing.T) {
	c := testCampaign(t, 5)
	p := NewProgress()
	if p.Count() != 0 {
		t.Fatalf("Count() = %d, want 0", p.Count())
	}
	level, duel, ok := p.Current(c)
	if !ok || level.Number() != 1 || duel.Number() != 1 {
		t.Fatalf("Current() = (%d, %d, %v), want (1, 1, true)", level.Number(), duel.Number(), ok)
	}
	if p.IsComplete(c) {
		t.Fatal("IsComplete() = true on a fresh progress, want false")
	}
	if p.HasWon(c, 1, 1) {
		t.Fatal("HasWon(1,1) = true on a fresh progress, want false")
	}
	if _, ok := p.LastRewards(c); ok {
		t.Fatal("LastRewards() ok = true on a fresh progress, want false")
	}
}

func TestProgressRecordVictoryAdvancesLinearly(t *testing.T) {
	c := testCampaign(t, 5)
	p := winFrom(t, c, NewProgress(), 4)
	level, duel, ok := p.Current(c)
	if !ok || level.Number() != 2 || duel.Number() != 1 {
		t.Fatalf("Current() = (%d, %d, %v), want (2, 1, true) after 4 victories", level.Number(), duel.Number(), ok)
	}
	if p.Count() != 4 {
		t.Fatalf("Count() = %d, want 4", p.Count())
	}
}

func TestProgressRecordLossKeepsState(t *testing.T) {
	c := testCampaign(t, 5)
	p := winFrom(t, c, NewProgress(), 2)
	if _, err := p.RecordLoss(c, 1, 3); err != nil {
		t.Fatalf("RecordLoss(1,3): %v", err)
	}
	got, err := p.RecordLoss(c, 1, 3)
	if err != nil {
		t.Fatalf("RecordLoss(1,3) error: %v", err)
	}
	if got.Count() != p.Count() {
		t.Fatalf("loss changed Count() to %d, want %d", got.Count(), p.Count())
	}
	level, duel, ok := got.Current(c)
	if !ok || level.Number() != 1 || duel.Number() != 3 {
		t.Fatalf("Current() = (%d, %d, %v), want (1, 3, true) after a loss", level.Number(), duel.Number(), ok)
	}
	// The retry wins and the player keeps the previous victories.
	retried, err := got.RecordVictory(c, 1, 3)
	if err != nil {
		t.Fatalf("RecordVictory after retry: %v", err)
	}
	if retried.Count() != 3 {
		t.Fatalf("Count() after retry = %d, want 3", retried.Count())
	}
}

func TestProgressRecordVictoryRejectsNonCurrentDuel(t *testing.T) {
	c := testCampaign(t, 5)
	p := winFrom(t, c, NewProgress(), 2)
	tests := []struct {
		name  string
		level int
		duel  int
	}{
		{name: "already won duel", level: 1, duel: 1},
		{name: "future duel in the level", level: 1, duel: 4},
		{name: "future level", level: 2, duel: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.RecordVictory(c, tt.level, tt.duel)
			if !errors.Is(err, ErrNotCurrentDuel) {
				t.Fatalf("err = %v, want wrapped %v", err, ErrNotCurrentDuel)
			}
			if got.Count() != p.Count() {
				t.Fatalf("rejected victory changed Count() to %d, want %d", got.Count(), p.Count())
			}
		})
	}
}

func TestProgressRecordVictoryRejectsUnknownDuel(t *testing.T) {
	c := testCampaign(t, 2)
	p := NewProgress()
	tests := []struct {
		name  string
		level int
		duel  int
	}{
		{name: "level above the campaign", level: 9, duel: 1},
		{name: "duel above the level", level: 1, duel: 9},
		{name: "zero level", level: 0, duel: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.RecordVictory(c, tt.level, tt.duel)
			if !errors.Is(err, ErrDuelNotFound) {
				t.Fatalf("err = %v, want wrapped %v", err, ErrDuelNotFound)
			}
			if got.Count() != 0 {
				t.Fatalf("rejected victory changed Count() to %d, want 0", got.Count())
			}
		})
	}
}

func TestProgressRecordLossRejectsNonCurrentAndUnknown(t *testing.T) {
	c := testCampaign(t, 2)
	p := winFrom(t, c, NewProgress(), 1)
	// The current duel is 1-2; 1-3 is in the future and 5-1 is unknown.
	if _, err := p.RecordLoss(c, 1, 3); !errors.Is(err, ErrNotCurrentDuel) {
		t.Fatalf("RecordLoss(1,3) err = %v, want wrapped %v", err, ErrNotCurrentDuel)
	}
	if _, err := p.RecordLoss(c, 5, 1); !errors.Is(err, ErrDuelNotFound) {
		t.Fatalf("RecordLoss(5,1) err = %v, want wrapped %v", err, ErrDuelNotFound)
	}
}

func TestProgressComplete(t *testing.T) {
	c := testCampaign(t, 2)
	p := winFrom(t, c, NewProgress(), c.DuelCount())
	if !p.IsComplete(c) {
		t.Fatal("IsComplete() = false after winning every duel, want true")
	}
	if _, _, ok := p.Current(c); ok {
		t.Fatal("Current() ok = true on a complete campaign, want false")
	}
	if _, err := p.RecordVictory(c, 1, 1); !errors.Is(err, ErrCampaignComplete) {
		t.Fatalf("RecordVictory err = %v, want wrapped %v", err, ErrCampaignComplete)
	}
	if _, err := p.RecordLoss(c, 1, 1); !errors.Is(err, ErrCampaignComplete) {
		t.Fatalf("RecordLoss err = %v, want wrapped %v", err, ErrCampaignComplete)
	}
	// A won duel can never be played or rewarded again.
	if !p.HasWon(c, 1, 1) || !p.HasWon(c, 2, 4) {
		t.Fatal("HasWon = false for won duels on a complete campaign, want true")
	}
}

func TestProgressHasWon(t *testing.T) {
	c := testCampaign(t, 5)
	p := winFrom(t, c, NewProgress(), 5)
	if !p.HasWon(c, 1, 4) {
		t.Fatal("HasWon(1,4) = false after 5 victories, want true")
	}
	if !p.HasWon(c, 2, 1) {
		t.Fatal("HasWon(2,1) = false after 5 victories, want true")
	}
	if p.HasWon(c, 2, 2) {
		t.Fatal("HasWon(2,2) = true after 5 victories, want false")
	}
	if p.HasWon(c, 9, 1) {
		t.Fatal("HasWon(9,1) = true for an unknown duel, want false")
	}
}

func TestProgressLastRewards(t *testing.T) {
	c := testCampaign(t, 5)

	t.Run("first victory grants the first duel rewards", func(t *testing.T) {
		p := winFrom(t, c, NewProgress(), 1)
		got, ok := p.LastRewards(c)
		if !ok || len(got) != 3 {
			t.Fatalf("LastRewards() = (%v, %v), want 3 cards", got, ok)
		}
		want := rewardIDs(t, c, 1, 1)
		for i, card := range got {
			if card.ID() != want[i] {
				t.Fatalf("LastRewards()[%d] = %q, want %q", i, card.ID(), want[i])
			}
		}
	})

	t.Run("fourth victory grants the last duel of the level", func(t *testing.T) {
		p := winFrom(t, c, NewProgress(), 4)
		got, ok := p.LastRewards(c)
		if !ok {
			t.Fatal("LastRewards() ok = false, want true")
		}
		want := rewardIDs(t, c, 1, 4)
		if got[0].ID() != want[0] {
			t.Fatalf("LastRewards()[0] = %q, want %q", got[0].ID(), want[0])
		}
	})

	t.Run("final victory grants the final duel rewards", func(t *testing.T) {
		p := winFrom(t, c, NewProgress(), c.DuelCount())
		got, ok := p.LastRewards(c)
		if !ok {
			t.Fatal("LastRewards() ok = false, want true")
		}
		want := rewardIDs(t, c, 5, 4)
		if got[0].ID() != want[0] {
			t.Fatalf("LastRewards()[0] = %q, want %q", got[0].ID(), want[0])
		}
	})
}

// rewardIDs returns the ids of the three rewards of the given duel.
func rewardIDs(t *testing.T, c Campaign, level, duel int) []domain.CardID {
	t.Helper()
	d, ok := c.Duel(level, duel)
	if !ok {
		t.Fatalf("Duel(%d,%d) not found", level, duel)
	}
	ids := make([]domain.CardID, 0, 3)
	for _, card := range d.Rewards() {
		ids = append(ids, card.ID())
	}
	return ids
}
