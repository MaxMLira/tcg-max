package seed

import (
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/domain"
	"github.com/MaxMLira/tcg-max/internal/story"
)

// duelSpec is the declarative description of one campaign duel: the
// opponent, its deck (catalog ids, 2 copies each) and the three predefined
// victory rewards (game-rules §58).
type duelSpec struct {
	charID   string
	charName string
	deck     []domain.CardID
	rewards  [3]domain.CardID
}

// levelSpec is the declarative description of one campaign level: its four
// duels, in the order they are played (game-rules §55 and §57).
type levelSpec struct {
	duels [story.DuelsPerLevel]duelSpec
}

// campaignSpecs is the campaign of the Story Mode MVP: 5 levels, 4 duels
// each (docs/mvp-story-mode.md implements levels 1 to 5 only), 20 distinct
// Characters and 20 distinct decks.
//
// The rewards follow the same curve as the decks: level 1 grants the strong
// base-tier cards, level 2 the tier B cards, level 3 the first Magic
// Equipment together with tier B, level 4 the tier C cards and level 5 the
// tier D cards (game-rules §58). Cards may repeat across different reward
// trios: the rule requires three cards per victory, not three new cards.
func campaignSpecs() [5]levelSpec {
	return [5]levelSpec{
		// Level 1: base tier only, no Equipment.
		{duels: [story.DuelsPerLevel]duelSpec{
			{charID: "char-richard", charName: "Richard", deck: deckL1D1,
				rewards: [3]domain.CardID{"drake-rider", "shield-captain", "hexweaver"}},
			{charID: "char-mira", charName: "Mira", deck: deckL1D2,
				rewards: [3]domain.CardID{"deep-current", "pit-spawn", "frost-drake"}},
			{charID: "char-tove", charName: "Tove", deck: deckL1D3,
				rewards: [3]domain.CardID{"knight", "enchantress", "fin-raider"}},
			{charID: "char-brann", charName: "Brann", deck: deckL1D4,
				rewards: [3]domain.CardID{"halberdier", "rune-master", "coral-guard"}},
		}},
		// Level 2: tier B (one tribute) enters the decks.
		{duels: [story.DuelsPerLevel]duelSpec{
			{charID: "char-kael", charName: "Kael", deck: deckL2D5,
				rewards: [3]domain.CardID{"molten-wyrm", "warlord", "storm-caller"}},
			{charID: "char-wren", charName: "Wren", deck: deckL2D6,
				rewards: [3]domain.CardID{"leviathan-hatchling", "ironclad-drake", "drake-rider"}},
			{charID: "char-dorian", charName: "Dorian", deck: deckL2D7,
				rewards: [3]domain.CardID{"warlord", "rune-master", "hexweaver"}},
			{charID: "char-sable", charName: "Sable", deck: deckL2D8,
				rewards: [3]domain.CardID{"storm-caller", "pit-spawn", "leviathan-hatchling"}},
		}},
		// Level 3: Magic Equipment enters the opponent decks.
		{duels: [story.DuelsPerLevel]duelSpec{
			{charID: "char-fen", charName: "Fen", deck: deckL3D9,
				rewards: [3]domain.CardID{"dragon-equipment", "molten-wyrm", "fin-raider"}},
			{charID: "char-lira", charName: "Lira", deck: deckL3D10,
				rewards: [3]domain.CardID{"warrior-equipment", "knight", "coral-guard"}},
			{charID: "char-orin", charName: "Orin", deck: deckL3D11,
				rewards: [3]domain.CardID{"mage-equipment", "enchantress", "shield-captain"}},
			{charID: "char-maud", charName: "Maud", deck: deckL3D12,
				rewards: [3]domain.CardID{"aquatic-equipment", "halberdier", "ironclad-drake"}},
		}},
		// Level 4: tier C (two tributes) enters the decks.
		{duels: [story.DuelsPerLevel]duelSpec{
			{charID: "char-cedric", charName: "Cedric", deck: deckL4D13,
				rewards: [3]domain.CardID{"elder-drake", "imperial-knight", "grand-sorcerer"}},
			{charID: "char-vesper", charName: "Vesper", deck: deckL4D14,
				rewards: [3]domain.CardID{"storm-dragon", "dragon-equipment", "abyss-koi"}},
			{charID: "char-haldor", charName: "Haldor", deck: deckL4D15,
				rewards: [3]domain.CardID{"archmage", "colossus", "warrior-equipment"}},
			{charID: "char-nyx", charName: "Nyx", deck: deckL4D16,
				rewards: [3]domain.CardID{"leviathan", "grand-sorcerer", "mage-equipment"}},
		}},
		// Level 5: tier D (three tributes), the strongest cards.
		{duels: [story.DuelsPerLevel]duelSpec{
			{charID: "char-magnus", charName: "Magnus", deck: deckL5D17,
				rewards: [3]domain.CardID{"eternal-flame-dragon", "obliterator-knight", "void-emperor"}},
			{charID: "char-serilda", charName: "Serilda", deck: deckL5D18,
				rewards: [3]domain.CardID{"world-ender", "titan-guard", "abyss-leviathan"}},
			{charID: "char-corvus", charName: "Corvus", deck: deckL5D19,
				rewards: [3]domain.CardID{"world-ender", "void-emperor", "aquatic-equipment"}},
			{charID: "char-alaric", charName: "Alaric", deck: deckL5D20,
				rewards: [3]domain.CardID{"titan-guard", "abyss-leviathan", "warrior-equipment"}},
		}},
	}
}

// Campaign returns the complete campaign of the Story Mode MVP: 5 levels,
// 4 duels each, 20 Characters and 20 predefined decks, with the three
// predefined rewards of each duel (game-rules §55 to §58,
// docs/mvp-story-mode.md).
//
// The build is deterministic: the same validated definitions are produced
// on every call, with no randomness and no time. Every construction error
// is propagated, so the campaign is only returned when the whole seed is
// valid.
func Campaign() (story.Campaign, error) {
	catalog, err := buildCatalog()
	if err != nil {
		return story.Campaign{}, err
	}
	specs := campaignSpecs()
	levels := make([]story.Level, 0, len(specs))
	for i, spec := range specs {
		levelNumber := i + 1
		duels := [story.DuelsPerLevel]story.Duel{}
		for d, ds := range spec.duels {
			duelNumber := d + 1
			opponentDeck, err := deck(catalog, ds.deck...)
			if err != nil {
				return story.Campaign{}, fmt.Errorf("seed: level %d duel %d: %w", levelNumber, duelNumber, err)
			}
			character, err := story.NewCharacter(ds.charID, ds.charName, opponentDeck)
			if err != nil {
				return story.Campaign{}, fmt.Errorf("seed: level %d duel %d: %w", levelNumber, duelNumber, err)
			}
			rewards := [3]*domain.Card{}
			for r, id := range ds.rewards {
				card, ok := catalog[id]
				if !ok {
					return story.Campaign{}, fmt.Errorf("seed: level %d duel %d: unknown reward %q", levelNumber, duelNumber, id)
				}
				rewards[r] = card
			}
			duel, err := story.NewDuel(duelNumber, character, rewards)
			if err != nil {
				return story.Campaign{}, fmt.Errorf("seed: level %d duel %d: %w", levelNumber, duelNumber, err)
			}
			duels[duelNumber-1] = duel
		}
		level, err := story.NewLevel(levelNumber, duels)
		if err != nil {
			return story.Campaign{}, fmt.Errorf("seed: level %d: %w", levelNumber, err)
		}
		levels = append(levels, level)
	}
	return story.NewCampaign(levels...)
}

// StarterDeck returns the player's starting deck of the Story Mode MVP
// (docs/mvp-story-mode.md §1.1: the player starts with a predefined deck;
// the deck builder is out of the MVP): the 15 base monsters, 2 copies of
// each, 30 cards (game-rules §4). The build is deterministic.
func StarterDeck() (domain.DeckDefinition, error) {
	catalog, err := buildCatalog()
	if err != nil {
		return domain.DeckDefinition{}, err
	}
	return deck(catalog, deckStarter...)
}
