package seed

import (
	"fmt"

	"github.com/MaxMLira/tcg-max/internal/domain"
)

// cardSpec is the declarative description of one card of the campaign
// catalog. For a Monster the subtype is the card's own subtype; for the
// Magic Equipment cards it is the accepted Monster subtype.
type cardSpec struct {
	ID      domain.CardID
	name    string
	typ     domain.CardType
	atk     int
	def     int
	weight  int
	subtype domain.MonsterSubtype
}

// catalogSpecs is the complete card catalog of the Story Mode MVP campaign.
//
// The monsters are organized in four weight tiers (tributes per
// game-rules §16). The tiers are the resource the difficulty grows with,
// level by level:
//
//   - base tier, weights 1 to 5, no tribute: the player's starter deck and
//     level 1;
//   - tier B, weights 6 to 8, one tribute: levels 2 and 3;
//   - tier C, weights 9 to 12, two tributes: levels 4 and 5;
//   - tier D, weights 13 and up, three tributes: level 5 only.
//
// The four Magic Equipment are the catalog of game-rules §35.2.1; their
// +800 ATK behavior is Equipment behavior of the Duel Engine (task S5) and
// is not modeled in the definition. There is no other Magic in the MVP:
// Burn and Power Up are out of scope (game-rules §10.2), and no Trap exists
// (game-rules §10.3).
func catalogSpecs() []cardSpec {
	return []cardSpec{
		// Base tier: weights 1 to 5, no tribute (game-rules §16).
		{ID: "ember-drake", name: "Ember Drake", typ: domain.CardTypeMonster, atk: 300, def: 200, weight: 1, subtype: domain.SubtypeDragon()},
		{ID: "scale-sentinel", name: "Scale Sentinel", typ: domain.CardTypeMonster, atk: 800, def: 600, weight: 3, subtype: domain.SubtypeDragon()},
		{ID: "cinder-wyrm", name: "Cinder Wyrm", typ: domain.CardTypeMonster, atk: 1000, def: 400, weight: 4, subtype: domain.SubtypeDragon()},
		{ID: "pit-dragon", name: "Pit Dragon", typ: domain.CardTypeMonster, atk: 1200, def: 500, weight: 5, subtype: domain.SubtypeDragon()},
		{ID: "squire", name: "Squire", typ: domain.CardTypeMonster, atk: 300, def: 300, weight: 1, subtype: domain.SubtypeWarrior()},
		{ID: "gryff", name: "Gryff", typ: domain.CardTypeMonster, atk: 700, def: 500, weight: 3, subtype: domain.SubtypeWarrior()},
		{ID: "lancer", name: "Lancer", typ: domain.CardTypeMonster, atk: 900, def: 400, weight: 4, subtype: domain.SubtypeWarrior()},
		{ID: "bulwark-guard", name: "Bulwark Guard", typ: domain.CardTypeMonster, atk: 1100, def: 800, weight: 5, subtype: domain.SubtypeWarrior()},
		{ID: "apprentice", name: "Apprentice", typ: domain.CardTypeMonster, atk: 200, def: 400, weight: 1, subtype: domain.SubtypeMage()},
		{ID: "scribe", name: "Scribe", typ: domain.CardTypeMonster, atk: 600, def: 500, weight: 3, subtype: domain.SubtypeMage()},
		{ID: "pyromancer", name: "Pyromancer", typ: domain.CardTypeMonster, atk: 950, def: 300, weight: 4, subtype: domain.SubtypeMage()},
		{ID: "arcane-adept", name: "Arcane Adept", typ: domain.CardTypeMonster, atk: 1150, def: 400, weight: 5, subtype: domain.SubtypeMage()},
		{ID: "pufferfish", name: "Pufferfish", typ: domain.CardTypeMonster, atk: 250, def: 400, weight: 1, subtype: domain.SubtypeAquatic()},
		{ID: "koi-spirit", name: "Koi Spirit", typ: domain.CardTypeMonster, atk: 650, def: 550, weight: 3, subtype: domain.SubtypeAquatic()},
		{ID: "tide-serpent", name: "Tide Serpent", typ: domain.CardTypeMonster, atk: 1050, def: 600, weight: 5, subtype: domain.SubtypeAquatic()},
		// Base tier, strongest cards of weights 4 and 5: the first cards
		// that are not in the player's starter deck.
		{ID: "drake-rider", name: "Drake Rider", typ: domain.CardTypeMonster, atk: 1250, def: 500, weight: 4, subtype: domain.SubtypeDragon()},
		{ID: "shield-captain", name: "Shield Captain", typ: domain.CardTypeMonster, atk: 1150, def: 900, weight: 4, subtype: domain.SubtypeWarrior()},
		{ID: "hexweaver", name: "Hexweaver", typ: domain.CardTypeMonster, atk: 1300, def: 400, weight: 5, subtype: domain.SubtypeMage()},
		{ID: "deep-current", name: "Deep Current", typ: domain.CardTypeMonster, atk: 1200, def: 650, weight: 4, subtype: domain.SubtypeAquatic()},
		{ID: "pit-spawn", name: "Pit Spawn", typ: domain.CardTypeMonster, atk: 1400, def: 300, weight: 5, subtype: domain.SubtypeDragon()},
		// Tier B: weights 6 to 8, one tribute (game-rules §16).
		{ID: "frost-drake", name: "Frost Drake", typ: domain.CardTypeMonster, atk: 1400, def: 700, weight: 6, subtype: domain.SubtypeDragon()},
		{ID: "ironclad-drake", name: "Ironclad Drake", typ: domain.CardTypeMonster, atk: 1600, def: 900, weight: 7, subtype: domain.SubtypeDragon()},
		{ID: "molten-wyrm", name: "Molten Wyrm", typ: domain.CardTypeMonster, atk: 1900, def: 800, weight: 8, subtype: domain.SubtypeDragon()},
		{ID: "knight", name: "Knight", typ: domain.CardTypeMonster, atk: 1500, def: 800, weight: 6, subtype: domain.SubtypeWarrior()},
		{ID: "halberdier", name: "Halberdier", typ: domain.CardTypeMonster, atk: 1700, def: 700, weight: 7, subtype: domain.SubtypeWarrior()},
		{ID: "warlord", name: "Warlord", typ: domain.CardTypeMonster, atk: 2000, def: 900, weight: 8, subtype: domain.SubtypeWarrior()},
		{ID: "enchantress", name: "Enchantress", typ: domain.CardTypeMonster, atk: 1450, def: 600, weight: 6, subtype: domain.SubtypeMage()},
		{ID: "rune-master", name: "Rune Master", typ: domain.CardTypeMonster, atk: 1650, def: 700, weight: 7, subtype: domain.SubtypeMage()},
		{ID: "storm-caller", name: "Storm Caller", typ: domain.CardTypeMonster, atk: 1850, def: 600, weight: 8, subtype: domain.SubtypeMage()},
		{ID: "fin-raider", name: "Fin Raider", typ: domain.CardTypeMonster, atk: 1400, def: 800, weight: 6, subtype: domain.SubtypeAquatic()},
		{ID: "coral-guard", name: "Coral Guard", typ: domain.CardTypeMonster, atk: 1550, def: 900, weight: 7, subtype: domain.SubtypeAquatic()},
		{ID: "leviathan-hatchling", name: "Leviathan Hatchling", typ: domain.CardTypeMonster, atk: 1750, def: 800, weight: 8, subtype: domain.SubtypeAquatic()},
		// Tier C: weights 9 to 12, two tributes (game-rules §16).
		{ID: "elder-drake", name: "Elder Drake", typ: domain.CardTypeMonster, atk: 2200, def: 1000, weight: 9, subtype: domain.SubtypeDragon()},
		{ID: "storm-dragon", name: "Storm Dragon", typ: domain.CardTypeMonster, atk: 2400, def: 1100, weight: 10, subtype: domain.SubtypeDragon()},
		{ID: "imperial-knight", name: "Imperial Knight", typ: domain.CardTypeMonster, atk: 2300, def: 1200, weight: 9, subtype: domain.SubtypeWarrior()},
		{ID: "colossus", name: "Colossus", typ: domain.CardTypeMonster, atk: 2800, def: 1400, weight: 12, subtype: domain.SubtypeWarrior()},
		{ID: "grand-sorcerer", name: "Grand Sorcerer", typ: domain.CardTypeMonster, atk: 2500, def: 900, weight: 10, subtype: domain.SubtypeMage()},
		{ID: "archmage", name: "Archmage", typ: domain.CardTypeMonster, atk: 2700, def: 1000, weight: 12, subtype: domain.SubtypeMage()},
		{ID: "abyss-koi", name: "Abyss Koi", typ: domain.CardTypeMonster, atk: 2250, def: 1100, weight: 9, subtype: domain.SubtypeAquatic()},
		{ID: "leviathan", name: "Leviathan", typ: domain.CardTypeMonster, atk: 2600, def: 1300, weight: 12, subtype: domain.SubtypeAquatic()},
		// Tier D: weights 13 and up, three tributes (game-rules §16).
		{ID: "eternal-flame-dragon", name: "Eternal Flame Dragon", typ: domain.CardTypeMonster, atk: 3000, def: 1200, weight: 13, subtype: domain.SubtypeDragon()},
		{ID: "world-ender", name: "World-Ender", typ: domain.CardTypeMonster, atk: 3400, def: 1000, weight: 14, subtype: domain.SubtypeDragon()},
		{ID: "obliterator-knight", name: "Obliterator Knight", typ: domain.CardTypeMonster, atk: 3200, def: 1500, weight: 13, subtype: domain.SubtypeWarrior()},
		{ID: "titan-guard", name: "Titan Guard", typ: domain.CardTypeMonster, atk: 3600, def: 1800, weight: 15, subtype: domain.SubtypeWarrior()},
		{ID: "void-emperor", name: "Void Emperor", typ: domain.CardTypeMonster, atk: 3500, def: 1100, weight: 14, subtype: domain.SubtypeMage()},
		{ID: "abyss-leviathan", name: "Abyss Leviathan", typ: domain.CardTypeMonster, atk: 3300, def: 1300, weight: 14, subtype: domain.SubtypeAquatic()},
		// The four Magic Equipment of game-rules §35.2.1.
		{ID: "dragon-equipment", name: "Dragon Equipment", typ: domain.CardTypeMagic, subtype: domain.SubtypeDragon()},
		{ID: "warrior-equipment", name: "Warrior Equipment", typ: domain.CardTypeMagic, subtype: domain.SubtypeWarrior()},
		{ID: "mage-equipment", name: "Mage Equipment", typ: domain.CardTypeMagic, subtype: domain.SubtypeMage()},
		{ID: "aquatic-equipment", name: "Aquatic Equipment", typ: domain.CardTypeMagic, subtype: domain.SubtypeAquatic()},
	}
}

// buildCatalog builds the validated card definitions of the campaign. The
// map is keyed by the card id so the decks and the rewards can reference
// the same validated definition.
func buildCatalog() (map[domain.CardID]*domain.Card, error) {
	specs := catalogSpecs()
	catalog := make(map[domain.CardID]*domain.Card, len(specs))
	for _, spec := range specs {
		var (
			card *domain.Card
			err  error
		)
		if spec.typ == domain.CardTypeMagic {
			card, err = domain.NewEquipmentCard(spec.ID, spec.name, spec.subtype)
		} else {
			card, err = domain.NewCard(spec.ID, spec.name, spec.typ, spec.atk, spec.def, spec.weight, spec.subtype)
		}
		if err != nil {
			return nil, fmt.Errorf("seed: card %q: %w", spec.ID, err)
		}
		catalog[spec.ID] = card
	}
	return catalog, nil
}
