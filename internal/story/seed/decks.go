package seed

import (
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// deck builds a validated deck of 30 cards from the given catalog card ids:
// 2 copies of each id (game-rules §4 and §4.1). Every deck of the seed has
// exactly 15 distinct cards at the copies limit.
func deck(catalog map[domain.CardID]*domain.Card, ids ...domain.CardID) (domain.DeckDefinition, error) {
	entries := make([]domain.DeckEntry, 0, len(ids))
	for _, id := range ids {
		card, ok := catalog[id]
		if !ok {
			return domain.DeckDefinition{}, fmt.Errorf("unknown card %q", id)
		}
		entries = append(entries, domain.DeckEntry{Card: card, Copies: domain.DeckCopiesLimit})
	}
	def, err := domain.NewDeckDefinition(entries)
	if err != nil {
		return domain.DeckDefinition{}, fmt.Errorf("deck of %d distinct cards: %w", len(ids), err)
	}
	return def, nil
}

var (
	// deckStarter is the player's starting deck: the 15 base monsters, 2
	// copies each (30 cards, game-rules §4).
	deckStarter = []domain.CardID{
		"ember-drake", "scale-sentinel", "cinder-wyrm", "pit-dragon",
		"squire", "gryff", "lancer", "bulwark-guard",
		"apprentice", "scribe", "pyromancer", "arcane-adept",
		"pufferfish", "koi-spirit", "tide-serpent",
	}

	// The lists below are the 20 opponent decks, one per duel. Level 1
	// plays only the base tier (weights 1 to 5, no tribute) and no
	// Equipment: within the level each duel replaces the weakest base
	// cards with the stronger ones. Level 2 introduces tier B (weights 6
	// to 8, one tribute). Level 3 introduces the Magic Equipment of the
	// opponent decks (game-rules §35.2). Level 4 introduces tier C
	// (weights 9 to 12, two tributes). Level 5 introduces tier D
	// (weights 13 and up, three tributes). This is the resource-based
	// difficulty progression of game-rules §57.
	deckL1D1 = []domain.CardID{
		"ember-drake", "scale-sentinel", "cinder-wyrm", "pit-dragon",
		"squire", "gryff", "lancer", "bulwark-guard",
		"apprentice", "scribe", "pyromancer", "arcane-adept",
		"pufferfish", "koi-spirit", "tide-serpent",
	}
	deckL1D2 = []domain.CardID{
		"scale-sentinel", "cinder-wyrm", "pit-dragon", "squire",
		"gryff", "lancer", "bulwark-guard", "scribe",
		"pyromancer", "arcane-adept", "pufferfish", "koi-spirit",
		"tide-serpent", "shield-captain", "drake-rider",
	}
	deckL1D3 = []domain.CardID{
		"scale-sentinel", "cinder-wyrm", "pit-dragon", "gryff",
		"lancer", "bulwark-guard", "scribe", "pyromancer",
		"arcane-adept", "koi-spirit", "tide-serpent", "drake-rider",
		"shield-captain", "deep-current", "pit-spawn",
	}
	deckL1D4 = []domain.CardID{
		"scale-sentinel", "cinder-wyrm", "pit-dragon", "gryff",
		"lancer", "bulwark-guard", "pyromancer", "arcane-adept",
		"koi-spirit", "tide-serpent", "drake-rider", "shield-captain",
		"hexweaver", "deep-current", "pit-spawn",
	}
	deckL2D5 = []domain.CardID{
		"scale-sentinel", "cinder-wyrm", "pit-dragon", "gryff",
		"lancer", "bulwark-guard", "scribe", "pyromancer",
		"arcane-adept", "koi-spirit", "tide-serpent", "frost-drake",
		"knight", "enchantress", "fin-raider",
	}
	deckL2D6 = []domain.CardID{
		"cinder-wyrm", "pit-dragon", "lancer", "bulwark-guard",
		"scribe", "pyromancer", "arcane-adept", "koi-spirit",
		"tide-serpent", "frost-drake", "knight", "enchantress",
		"fin-raider", "ironclad-drake", "halberdier",
	}
	deckL2D7 = []domain.CardID{
		"pit-dragon", "lancer", "bulwark-guard", "pyromancer",
		"arcane-adept", "koi-spirit", "tide-serpent", "frost-drake",
		"knight", "enchantress", "fin-raider", "ironclad-drake",
		"halberdier", "rune-master", "coral-guard",
	}
	deckL2D8 = []domain.CardID{
		"pit-dragon", "bulwark-guard", "arcane-adept", "tide-serpent",
		"cinder-wyrm", "frost-drake", "knight", "enchantress",
		"fin-raider", "ironclad-drake", "halberdier", "rune-master",
		"coral-guard", "leviathan-hatchling", "warlord",
	}
	deckL3D9 = []domain.CardID{
		"pit-dragon", "bulwark-guard", "arcane-adept", "tide-serpent",
		"cinder-wyrm", "lancer", "frost-drake", "knight",
		"enchantress", "fin-raider", "ironclad-drake", "halberdier",
		"rune-master", "coral-guard", "dragon-equipment",
	}
	deckL3D10 = []domain.CardID{
		"pit-dragon", "bulwark-guard", "arcane-adept", "tide-serpent",
		"frost-drake", "knight", "enchantress", "fin-raider",
		"ironclad-drake", "halberdier", "rune-master", "coral-guard",
		"leviathan-hatchling", "warlord", "warrior-equipment",
	}
	deckL3D11 = []domain.CardID{
		"pit-dragon", "bulwark-guard", "arcane-adept", "tide-serpent",
		"frost-drake", "knight", "enchantress", "fin-raider",
		"ironclad-drake", "halberdier", "rune-master", "coral-guard",
		"leviathan-hatchling", "mage-equipment", "aquatic-equipment",
	}
	deckL3D12 = []domain.CardID{
		"pit-dragon", "bulwark-guard", "frost-drake", "knight",
		"enchantress", "fin-raider", "ironclad-drake", "halberdier",
		"rune-master", "coral-guard", "leviathan-hatchling", "warlord",
		"molten-wyrm", "dragon-equipment", "warrior-equipment",
	}
	deckL4D13 = []domain.CardID{
		"pit-dragon", "bulwark-guard", "frost-drake", "ironclad-drake",
		"halberdier", "leviathan-hatchling", "warlord", "molten-wyrm",
		"storm-caller", "elder-drake", "abyss-koi", "imperial-knight",
		"storm-dragon", "grand-sorcerer", "mage-equipment",
	}
	deckL4D14 = []domain.CardID{
		"bulwark-guard", "ironclad-drake", "halberdier",
		"leviathan-hatchling", "warlord", "molten-wyrm", "storm-caller",
		"elder-drake", "imperial-knight", "storm-dragon", "grand-sorcerer",
		"leviathan", "archmage", "dragon-equipment", "warrior-equipment",
	}
	deckL4D15 = []domain.CardID{
		"ironclad-drake", "halberdier", "leviathan-hatchling", "warlord",
		"molten-wyrm", "storm-caller", "elder-drake", "abyss-koi",
		"imperial-knight", "storm-dragon", "grand-sorcerer", "leviathan",
		"archmage", "colossus", "aquatic-equipment",
	}
	deckL4D16 = []domain.CardID{
		"halberdier", "leviathan-hatchling", "warlord", "molten-wyrm",
		"storm-caller", "elder-drake", "abyss-koi", "imperial-knight",
		"storm-dragon", "grand-sorcerer", "leviathan", "archmage",
		"colossus", "mage-equipment", "aquatic-equipment",
	}
	deckL5D17 = []domain.CardID{
		"warlord", "molten-wyrm", "storm-caller", "storm-dragon",
		"grand-sorcerer", "leviathan", "archmage", "colossus",
		"eternal-flame-dragon", "obliterator-knight", "abyss-leviathan",
		"world-ender", "void-emperor", "titan-guard", "dragon-equipment",
	}
	deckL5D18 = []domain.CardID{
		"warlord", "molten-wyrm", "imperial-knight", "storm-dragon",
		"leviathan", "archmage", "colossus", "eternal-flame-dragon",
		"obliterator-knight", "abyss-leviathan", "world-ender",
		"void-emperor", "titan-guard", "warrior-equipment", "mage-equipment",
	}
	deckL5D19 = []domain.CardID{
		"warlord", "molten-wyrm", "storm-caller", "imperial-knight",
		"storm-dragon", "grand-sorcerer", "leviathan", "archmage",
		"eternal-flame-dragon", "obliterator-knight", "abyss-leviathan",
		"world-ender", "void-emperor", "titan-guard", "warrior-equipment",
	}
	deckL5D20 = []domain.CardID{
		"elder-drake", "abyss-koi", "imperial-knight", "storm-dragon",
		"grand-sorcerer", "leviathan", "archmage", "eternal-flame-dragon",
		"obliterator-knight", "abyss-leviathan", "world-ender",
		"void-emperor", "titan-guard", "mage-equipment", "aquatic-equipment",
	}
)
