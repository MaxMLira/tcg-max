package seed

import (
	"sort"
	"strings"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/domain"
	"github.com/MaxMLira/tcg-max/internal/story"
)

// deckSignature is the sorted list of card ids of a deck, as a stable
// comparison key.
func deckSignature(deck domain.DeckDefinition) string {
	ids := make([]string, 0, len(deck.Entries()))
	for _, entry := range deck.Entries() {
		ids = append(ids, string(entry.Card.ID()))
	}
	sort.Strings(ids)
	return strings.Join(ids, ",")
}

// campaignSignature is the stable structure of the whole campaign: the
// characters, the decks and the rewards of every duel, in order.
func campaignSignature(c story.Campaign) string {
	parts := make([]string, 0, c.DuelCount())
	for _, level := range c.Levels() {
		for _, duel := range level.Duels() {
			rewards := make([]string, 0, 3)
			for _, card := range duel.Rewards() {
				rewards = append(rewards, string(card.ID()))
			}
			parts = append(parts, strings.Join([]string{
				duel.Character().ID(),
				duel.Character().Name(),
				deckSignature(duel.Character().Deck()),
				strings.Join(rewards, "+"),
			}, "|"))
		}
	}
	return strings.Join(parts, ";")
}

func TestCampaignStructure(t *testing.T) {
	c, err := Campaign()
	if err != nil {
		t.Fatalf("Campaign(): %v", err)
	}
	if c.LevelCount() != 5 {
		t.Fatalf("LevelCount() = %d, want 5", c.LevelCount())
	}
	if c.DuelCount() != 20 {
		t.Fatalf("DuelCount() = %d, want 20", c.DuelCount())
	}

	characters := map[string]struct{}{}
	decks := map[string]struct{}{}
	for levelNumber := 1; levelNumber <= 5; levelNumber++ {
		level, ok := c.Level(levelNumber)
		if !ok {
			t.Fatalf("Level(%d) not found", levelNumber)
		}
		// The initial Life Points of the level are the domain rule
		// (game-rules §3.1): 4000 for levels 1 to 4, 8000 for level 5.
		wantLP := 4000
		if levelNumber == 5 {
			wantLP = 8000
		}
		lp, err := level.InitialLifePoints()
		if err != nil {
			t.Fatalf("level %d InitialLifePoints(): %v", levelNumber, err)
		}
		if lp != wantLP {
			t.Fatalf("level %d InitialLifePoints() = %d, want %d", levelNumber, lp, wantLP)
		}
		for duelNumber := 1; duelNumber <= story.DuelsPerLevel; duelNumber++ {
			duel, ok := level.Duel(duelNumber)
			if !ok {
				t.Fatalf("level %d Duel(%d) not found", levelNumber, duelNumber)
			}
			char := duel.Character()
			if _, dup := characters[char.ID()]; dup {
				t.Fatalf("character %s appears in more than one duel", char.ID())
			}
			characters[char.ID()] = struct{}{}
			sig := deckSignature(char.Deck())
			if _, dup := decks[sig]; dup {
				t.Fatalf("level %d duel %d: deck composition is repeated", levelNumber, duelNumber)
			}
			decks[sig] = struct{}{}
			rewards := duel.Rewards()
			if len(rewards) != 3 {
				t.Fatalf("level %d duel %d: %d rewards, want 3", levelNumber, duelNumber, len(rewards))
			}
			for _, card := range rewards {
				if card == nil {
					t.Fatalf("level %d duel %d: nil reward", levelNumber, duelNumber)
				}
			}
		}
	}
	if len(characters) != 20 {
		t.Fatalf("%d distinct characters, want 20", len(characters))
	}
	if len(decks) != 20 {
		t.Fatalf("%d distinct decks, want 20", len(decks))
	}
}

func TestCampaignHasNoTrapBurnOrPowerUp(t *testing.T) {
	c, err := Campaign()
	if err != nil {
		t.Fatalf("Campaign(): %v", err)
	}
	mvpSubtypes := []domain.MonsterSubtype{
		domain.SubtypeDragon(),
		domain.SubtypeWarrior(),
		domain.SubtypeMage(),
		domain.SubtypeAquatic(),
	}
	accepted := map[domain.MonsterSubtype]struct{}{}
	equipment := 0
	seen := map[domain.CardID]struct{}{}
	inspect := func(card *domain.Card) {
		if _, dup := seen[card.ID()]; dup {
			return
		}
		seen[card.ID()] = struct{}{}
		switch card.Type() {
		case domain.CardTypeMonster:
			if card.Subtype().IsZero() {
				t.Fatalf("monster %s has no subtype", card.ID())
			}
			for _, s := range mvpSubtypes {
				if card.Subtype().Equals(s) {
					return
				}
			}
			t.Fatalf("monster %s has subtype %q, want one of the four MVP subtypes", card.ID(), card.Subtype())
		case domain.CardTypeMagic:
			// The only Magic of the MVP is the four Magic Equipment
			// (game-rules §35.2.1); any other Magic (Burn, Power Up)
			// declares no accepted subtype and is out of scope.
			if card.AcceptedMonsterSubtype().IsZero() {
				t.Fatalf("magic %s is not a Magic Equipment: Burn and Power Up are out of the MVP", card.ID())
			}
			equipment++
			accepted[card.AcceptedMonsterSubtype()] = struct{}{}
		case domain.CardTypeTrap:
			t.Fatalf("trap %s: Traps are out of the MVP", card.ID())
		}
	}
	for _, level := range c.Levels() {
		for _, duel := range level.Duels() {
			for _, entry := range duel.Character().Deck().Entries() {
				inspect(entry.Card)
			}
			for _, card := range duel.Rewards() {
				inspect(card)
			}
		}
	}
	if equipment != 4 {
		t.Fatalf("%d Magic Equipment in the campaign, want the four of game-rules §35.2.1", equipment)
	}
	for _, s := range mvpSubtypes {
		if _, ok := accepted[s]; !ok {
			t.Fatalf("no Magic Equipment accepting subtype %q", s)
		}
	}
}

func TestStarterDeck(t *testing.T) {
	deck, err := StarterDeck()
	if err != nil {
		t.Fatalf("StarterDeck(): %v", err)
	}
	if err := deck.Validate(); err != nil {
		t.Fatalf("StarterDeck() is invalid: %v", err)
	}
	if got := deck.TotalCards(); got != domain.DeckSize {
		t.Fatalf("TotalCards() = %d, want %d", got, domain.DeckSize)
	}
	base := map[string]struct{}{}
	for _, id := range deckStarter {
		base[string(id)] = struct{}{}
	}
	entries := deck.Entries()
	if len(entries) != 15 {
		t.Fatalf("%d entries, want 15 distinct cards", len(entries))
	}
	for _, entry := range entries {
		if entry.Card.Type() != domain.CardTypeMonster {
			t.Fatalf("starter deck card %s is not a Monster", entry.Card.ID())
		}
		if _, ok := base[string(entry.Card.ID())]; !ok {
			t.Fatalf("starter deck card %s is not a base-tier card", entry.Card.ID())
		}
	}
}

func TestSeedDeterministic(t *testing.T) {
	first, err := Campaign()
	if err != nil {
		t.Fatalf("Campaign(): %v", err)
	}
	second, err := Campaign()
	if err != nil {
		t.Fatalf("Campaign() again: %v", err)
	}
	if campaignSignature(first) != campaignSignature(second) {
		t.Fatal("two builds of Campaign() differ: the seed is not deterministic")
	}
	starterFirst, err := StarterDeck()
	if err != nil {
		t.Fatalf("StarterDeck(): %v", err)
	}
	starterSecond, err := StarterDeck()
	if err != nil {
		t.Fatalf("StarterDeck() again: %v", err)
	}
	if deckSignature(starterFirst) != deckSignature(starterSecond) {
		t.Fatal("two builds of StarterDeck() differ: the seed is not deterministic")
	}
}
