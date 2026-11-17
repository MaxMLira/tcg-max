package story

import (
	"fmt"
	"testing"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// testSubtypes cycles through the four subtypes of the Story Mode MVP.
var testSubtypes = []domain.MonsterSubtype{
	domain.SubtypeDragon(),
	domain.SubtypeWarrior(),
	domain.SubtypeMage(),
	domain.SubtypeAquatic(),
}

// testCard builds a valid monster card with the given index. The id is
// unique per index, so cards of the same index are interchangeable in the
// tests but distinct inside one deck.
func testCard(t *testing.T, i int) *domain.Card {
	t.Helper()
	card, err := domain.NewCard(
		domain.CardID(fmt.Sprintf("card-%03d", i)),
		fmt.Sprintf("Test Card %d", i),
		domain.CardTypeMonster, 100, 50, 1,
		testSubtypes[i%len(testSubtypes)],
	)
	if err != nil {
		t.Fatalf("building test card %d: %v", i, err)
	}
	return card
}

// testDeck builds a valid deck: 15 distinct cards, 2 copies each
// (game-rules §4 and §4.1).
func testDeck(t *testing.T) domain.DeckDefinition {
	t.Helper()
	entries := make([]domain.DeckEntry, 0, 15)
	for i := 0; i < 15; i++ {
		entries = append(entries, domain.DeckEntry{Card: testCard(t, i), Copies: 2})
	}
	deck, err := domain.NewDeckDefinition(entries)
	if err != nil {
		t.Fatalf("building test deck: %v", err)
	}
	return deck
}

// testCharacter builds a valid character with the given id and name.
func testCharacter(t *testing.T, id, name string) Character {
	t.Helper()
	c, err := NewCharacter(id, name, testDeck(t))
	if err != nil {
		t.Fatalf("building test character %q: %v", id, err)
	}
	return c
}

// testDuel builds a valid duel with the given number, character and three
// distinct valid reward cards.
func testDuel(t *testing.T, number int, character Character) Duel {
	t.Helper()
	d, err := NewDuel(number, character, [3]*domain.Card{
		testCard(t, 100+number),
		testCard(t, 101+number),
		testCard(t, 102+number),
	})
	if err != nil {
		t.Fatalf("building test duel %d: %v", number, err)
	}
	return d
}

// testCampaign builds a valid campaign with the given number of levels,
// each with four duels and a distinct character per duel.
func testCampaign(t *testing.T, levels int) Campaign {
	t.Helper()
	built := make([]Level, 0, levels)
	for l := 1; l <= levels; l++ {
		duels := [DuelsPerLevel]Duel{}
		for d := 1; d <= DuelsPerLevel; d++ {
			duels[d-1] = testDuel(t, d, testCharacter(t,
				fmt.Sprintf("char-%d-%d", l, d),
				fmt.Sprintf("Character %d-%d", l, d),
			))
		}
		level, err := NewLevel(l, duels)
		if err != nil {
			t.Fatalf("building test level %d: %v", l, err)
		}
		built = append(built, level)
	}
	campaign, err := NewCampaign(built...)
	if err != nil {
		t.Fatalf("building test campaign: %v", err)
	}
	return campaign
}
