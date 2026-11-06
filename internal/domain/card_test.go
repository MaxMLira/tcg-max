package domain

import (
	"errors"
	"strings"
	"testing"
)

// mustCard builds a valid card definition or fails the test. Tests never
// construct a Card by bypassing NewCard's validation.
func mustCard(t *testing.T, id CardID, name string, typ CardType, atk, def, weight int, subtype MonsterSubtype) *Card {
	t.Helper()

	card, err := NewCard(id, name, typ, atk, def, weight, subtype)
	if err != nil {
		t.Fatalf("NewCard(%q, %q, %v, %d, %d, %d, %q) failed: %v", id, name, typ, atk, def, weight, subtype, err)
	}
	return card
}

func TestNewCard_CreatesDefinitionForEachCardType(t *testing.T) {
	tests := []struct {
		name    string
		typ     CardType
		subtype MonsterSubtype
	}{
		{name: "Monster", typ: CardTypeMonster, subtype: SubtypeDragon()},
		{name: "Trap", typ: CardTypeTrap, subtype: MonsterSubtype{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card, err := NewCard(CardID("CARD-"+tc.name), "Card "+tc.name, tc.typ, 1500, 1200, 4, tc.subtype)

			// Assert
			if err != nil {
				t.Fatalf("NewCard returned unexpected error: %v", err)
			}
			if card == nil {
				t.Fatal("NewCard returned a nil card without an error")
			}
			if card.ID() != CardID("CARD-"+tc.name) {
				t.Errorf("ID() = %q, want %q", card.ID(), "CARD-"+tc.name)
			}
			if card.Name() != "Card "+tc.name {
				t.Errorf("Name() = %q, want %q", card.Name(), "Card "+tc.name)
			}
			if card.Type() != tc.typ {
				t.Errorf("Type() = %v, want %v", card.Type(), tc.typ)
			}
			if card.ATK() != 1500 {
				t.Errorf("ATK() = %d, want 1500", card.ATK())
			}
			if card.DEF() != 1200 {
				t.Errorf("DEF() = %d, want 1200", card.DEF())
			}
			if card.Weight() != 4 {
				t.Errorf("Weight() = %d, want 4", card.Weight())
			}
			if !card.Subtype().Equals(tc.subtype) {
				t.Errorf("Subtype() = %q, want %q", card.Subtype(), tc.subtype)
			}
		})
	}
}

func TestNewCard_MagicHasNoCombatStatsOrWeight(t *testing.T) {
	// Arrange

	// Act
	card, err := NewCard("MAGIC-1", "Burn", CardTypeMagic, 0, 0, 0, MonsterSubtype{})

	// Assert
	if err != nil {
		t.Fatalf("NewCard returned unexpected error: %v", err)
	}
	if card.Type() != CardTypeMagic || card.ATK() != 0 || card.DEF() != 0 || card.Weight() != 0 {
		t.Errorf("card = (%v, %d, %d, %d), want (Magic, 0, 0, 0)",
			card.Type(), card.ATK(), card.DEF(), card.Weight())
	}
}

func TestNewCard_RejectsMagicWithMonsterAttributes(t *testing.T) {
	tests := []struct {
		name   string
		atk    int
		def    int
		weight int
	}{
		{name: "atk", atk: 1, def: 0, weight: 0},
		{name: "def", atk: 0, def: 1, weight: 0},
		{name: "weight", atk: 0, def: 0, weight: 1},
		{name: "negative atk", atk: -1, def: 0, weight: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card, err := NewCard("MAGIC-1", "Burn", CardTypeMagic, tc.atk, tc.def, tc.weight, MonsterSubtype{})

			// Assert
			if !errors.Is(err, ErrInvalidCard) {
				t.Fatalf("error = %v, want ErrInvalidCard", err)
			}
			if card != nil {
				t.Errorf("card = %v, want nil", card)
			}
		})
	}
}

func TestNewCard_AcceptsBoundaryValues(t *testing.T) {
	tests := []struct {
		name   string
		atk    int
		def    int
		weight int
	}{
		{name: "zero attack and defense", atk: 0, def: 0, weight: 1},
		{name: "weight above tribute tiers", atk: 3000, def: 2500, weight: 13},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card, err := NewCard("CARD-1", "Card", CardTypeMonster, tc.atk, tc.def, tc.weight, SubtypeDragon())

			// Assert
			if err != nil {
				t.Fatalf("NewCard returned unexpected error: %v", err)
			}
			if card.ATK() != tc.atk || card.DEF() != tc.def || card.Weight() != tc.weight {
				t.Errorf("values = (%d, %d, %d), want (%d, %d, %d)",
					card.ATK(), card.DEF(), card.Weight(), tc.atk, tc.def, tc.weight)
			}
		})
	}
}

func TestNewCard_RejectsInvalidDefinition(t *testing.T) {
	tests := []struct {
		name   string
		id     CardID
		cname  string
		typ    CardType
		atk    int
		def    int
		weight int
	}{
		{name: "empty id", id: "", cname: "Card", typ: CardTypeMonster, atk: 1000, def: 1000, weight: 4},
		{name: "blank id with spaces", id: "   ", cname: "Card", typ: CardTypeMonster, atk: 1000, def: 1000, weight: 4},
		{name: "blank id with tab", id: "\t", cname: "Card", typ: CardTypeMonster, atk: 1000, def: 1000, weight: 4},
		{name: "blank name with spaces", id: "CARD-1", cname: "  ", typ: CardTypeMonster, atk: 1000, def: 1000, weight: 4},
		{name: "blank name with newline", id: "CARD-1", cname: "\n ", typ: CardTypeMonster, atk: 1000, def: 1000, weight: 4},
		{name: "empty name", id: "CARD-1", cname: "", typ: CardTypeMonster, atk: 1000, def: 1000, weight: 4},
		{name: "unknown card type", id: "CARD-1", cname: "Card", typ: CardType(99), atk: 1000, def: 1000, weight: 4},
		{name: "zero card type", id: "CARD-1", cname: "Card", typ: CardType(0), atk: 1000, def: 1000, weight: 4},
		{name: "negative attack", id: "CARD-1", cname: "Card", typ: CardTypeMonster, atk: -1, def: 1000, weight: 4},
		{name: "negative defense", id: "CARD-1", cname: "Card", typ: CardTypeMonster, atk: 1000, def: -1, weight: 4},
		{name: "zero weight", id: "CARD-1", cname: "Card", typ: CardTypeMonster, atk: 1000, def: 1000, weight: 0},
		{name: "negative weight", id: "CARD-1", cname: "Card", typ: CardTypeMonster, atk: 1000, def: 1000, weight: -5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card, err := NewCard(tc.id, tc.cname, tc.typ, tc.atk, tc.def, tc.weight, SubtypeDragon())

			// Assert
			if !errors.Is(err, ErrInvalidCard) {
				t.Fatalf("error = %v, want ErrInvalidCard", err)
			}
			if card != nil {
				t.Errorf("card = %v, want nil", card)
			}
		})
	}
}

func TestNewCard_RejectsMonsterWithoutSubtype(t *testing.T) {
	// Arrange

	// Act
	card, err := NewCard("CARD-1", "Card", CardTypeMonster, 1000, 1000, 4, MonsterSubtype{})

	// Assert
	if !errors.Is(err, ErrInvalidCard) {
		t.Fatalf("error = %v, want ErrInvalidCard", err)
	}
	if card != nil {
		t.Errorf("card = %v, want nil", card)
	}
	if !strings.Contains(err.Error(), "subtype") {
		t.Errorf("error message %q does not mention the subtype", err)
	}
}

func TestNewCard_MonsterStoresAndReturnsItsSubtype(t *testing.T) {
	tests := []struct {
		name    string
		subtype MonsterSubtype
	}{
		{name: "Dragon", subtype: SubtypeDragon()},
		{name: "Warrior", subtype: SubtypeWarrior()},
		{name: "Mage", subtype: SubtypeMage()},
		{name: "Aquatic", subtype: SubtypeAquatic()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card := mustCard(t, "CARD-1", "Card", CardTypeMonster, 1000, 1000, 4, tc.subtype)

			// Assert
			if card.Subtype().IsZero() || !card.Subtype().Equals(tc.subtype) {
				t.Errorf("Subtype() = %q, want %q", card.Subtype(), tc.subtype)
			}
		})
	}
}

func TestNewCard_MonsterAcceptsANewSubtypeWithoutDomainChanges(t *testing.T) {
	// Arrange
	beast, err := NewMonsterSubtype("Beast")
	if err != nil {
		t.Fatalf("NewMonsterSubtype failed: %v", err)
	}

	// Act
	card := mustCard(t, "CARD-1", "Card", CardTypeMonster, 1000, 1000, 4, beast)

	// Assert
	if !card.Subtype().Equals(beast) {
		t.Errorf("Subtype() = %q, want %q", card.Subtype(), beast)
	}
}

func TestNewCard_RejectsMagicAndTrapWithASubtype(t *testing.T) {
	tests := []struct {
		name   string
		typ    CardType
		atk    int
		def    int
		weight int
	}{
		{name: "Magic", typ: CardTypeMagic, atk: 0, def: 0, weight: 0},
		{name: "Trap", typ: CardTypeTrap, atk: 1000, def: 1000, weight: 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card, err := NewCard("CARD-1", "Card", tc.typ, tc.atk, tc.def, tc.weight, SubtypeDragon())

			// Assert
			if !errors.Is(err, ErrInvalidCard) {
				t.Fatalf("error = %v, want ErrInvalidCard", err)
			}
			if card != nil {
				t.Errorf("card = %v, want nil", card)
			}
			if !strings.Contains(err.Error(), "subtype") {
				t.Errorf("error message %q does not mention the subtype", err)
			}
		})
	}
}

func TestNewCard_MagicAndTrapWithoutSubtypeHaveZeroSubtype(t *testing.T) {
	tests := []struct {
		name   string
		typ    CardType
		atk    int
		def    int
		weight int
	}{
		{name: "Magic", typ: CardTypeMagic, atk: 0, def: 0, weight: 0},
		{name: "Trap", typ: CardTypeTrap, atk: 1000, def: 1000, weight: 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card := mustCard(t, "CARD-1", "Card", tc.typ, tc.atk, tc.def, tc.weight, MonsterSubtype{})

			// Assert
			if !card.Subtype().IsZero() {
				t.Errorf("Subtype() = %q, want the zero value", card.Subtype())
			}
		})
	}
}

func TestNewCard_ErrorDescribesTheViolation(t *testing.T) {
	// Arrange

	// Act
	_, err := NewCard("", "Card", CardTypeMonster, 1000, 1000, 4, SubtypeDragon())

	// Assert
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !strings.Contains(err.Error(), "id") {
		t.Errorf("error message %q does not mention the offending field", err)
	}
}

func TestCardType_StringNamesSupportedTypes(t *testing.T) {
	tests := []struct {
		typ  CardType
		want string
	}{
		{typ: CardTypeMonster, want: "Monster"},
		{typ: CardTypeMagic, want: "Magic"},
		{typ: CardTypeTrap, want: "Trap"},
		{typ: CardType(99), want: "CardType(99)"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			// Arrange

			// Act
			got := tc.typ.String()

			// Assert
			if got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNewCard_MagicWithoutAcceptedSubtypeIsValidAndAcceptsNothing(t *testing.T) {
	// Arrange
	card := mustCard(t, "MAGIC-1", "Burn", CardTypeMagic, 0, 0, 0, MonsterSubtype{})

	// Act
	accepted := card.AcceptedMonsterSubtype()

	// Assert
	if !accepted.IsZero() {
		t.Errorf("AcceptedMonsterSubtype() = %q, want the zero value", accepted)
	}
	for _, subtype := range []MonsterSubtype{SubtypeDragon(), SubtypeWarrior(), SubtypeMage(), SubtypeAquatic(), {}} {
		if card.AcceptsMonsterSubtype(subtype) {
			t.Errorf("AcceptsMonsterSubtype(%q) = true, want false", subtype)
		}
	}
}

func TestNewEquipmentCard_BuildsAMagicCardThatDeclaresItsAcceptedSubtype(t *testing.T) {
	// Arrange

	// Act
	card, err := NewEquipmentCard("EQUIP-1", "Dragon Equipment", SubtypeDragon())

	// Assert
	if err != nil {
		t.Fatalf("NewEquipmentCard returned unexpected error: %v", err)
	}
	if card == nil {
		t.Fatal("NewEquipmentCard returned a nil card without an error")
	}
	if card.ID() != "EQUIP-1" || card.Name() != "Dragon Equipment" {
		t.Errorf("identity = (%q, %q), want (EQUIP-1, Dragon Equipment)", card.ID(), card.Name())
	}
	if card.Type() != CardTypeMagic {
		t.Errorf("Type() = %v, want Magic", card.Type())
	}
	if card.ATK() != 0 || card.DEF() != 0 || card.Weight() != 0 {
		t.Errorf("values = (%d, %d, %d), want (0, 0, 0)", card.ATK(), card.DEF(), card.Weight())
	}
	if !card.Subtype().IsZero() {
		t.Errorf("Subtype() = %q, want the zero value", card.Subtype())
	}
	if !card.AcceptedMonsterSubtype().Equals(SubtypeDragon()) {
		t.Errorf("AcceptedMonsterSubtype() = %q, want %q", card.AcceptedMonsterSubtype(), SubtypeDragon())
	}
}

func TestNewEquipmentCard_OwnSubtypeAndAcceptedSubtypeAreDistinctFacts(t *testing.T) {
	// Arrange
	card, err := NewEquipmentCard("EQUIP-1", "Mage Equipment", SubtypeMage())
	if err != nil {
		t.Fatalf("NewEquipmentCard returned unexpected error: %v", err)
	}

	// Act
	own := card.Subtype()
	accepted := card.AcceptedMonsterSubtype()

	// Assert
	if !own.IsZero() {
		t.Errorf("Subtype() = %q, want the zero value", own)
	}
	if accepted.IsZero() {
		t.Error("AcceptedMonsterSubtype() is the zero value, want a subtype")
	}
}

func TestNewEquipmentCard_SubtypeIsTheZeroValue(t *testing.T) {
	accepted := []MonsterSubtype{SubtypeDragon(), SubtypeWarrior(), SubtypeMage(), SubtypeAquatic()}

	for _, subtype := range accepted {
		t.Run(subtype.String(), func(t *testing.T) {
			// Arrange
			card, err := NewEquipmentCard("EQUIP-1", subtype.String()+" Equipment", subtype)
			if err != nil {
				t.Fatalf("NewEquipmentCard returned unexpected error: %v", err)
			}

			// Act
			own := card.Subtype()

			// Assert
			if own != (MonsterSubtype{}) {
				t.Errorf("Subtype() = %q, want the zero value", own)
			}
			if !own.IsZero() {
				t.Errorf("Subtype().IsZero() = false, want true")
			}
		})
	}
}

func TestEquipment_DoesNotTakeAMonsterSubtypeAsItsOwn(t *testing.T) {
	subtypes := []MonsterSubtype{SubtypeDragon(), SubtypeWarrior(), SubtypeMage(), SubtypeAquatic()}

	for _, subtype := range subtypes {
		t.Run(subtype.String(), func(t *testing.T) {
			// Arrange

			// Act
			rejected, err := NewCard("EQUIP-1", subtype.String()+" Equipment", CardTypeMagic, 0, 0, 0, subtype)
			equipment, equipmentErr := NewEquipmentCard("EQUIP-2", subtype.String()+" Equipment", subtype)

			// Assert
			if !errors.Is(err, ErrInvalidCard) {
				t.Fatalf("NewCard(Magic, %q) error = %v, want ErrInvalidCard", subtype, err)
			}
			if rejected != nil {
				t.Errorf("card = %v, want nil", rejected)
			}
			if equipmentErr != nil {
				t.Fatalf("NewEquipmentCard returned unexpected error: %v", equipmentErr)
			}
			if !equipment.Subtype().IsZero() {
				t.Errorf("Subtype() = %q, want the zero value", equipment.Subtype())
			}
			if !equipment.AcceptedMonsterSubtype().Equals(subtype) {
				t.Errorf("AcceptedMonsterSubtype() = %q, want %q", equipment.AcceptedMonsterSubtype(), subtype)
			}
		})
	}
}

func TestNewEquipmentCard_RejectsInvalidDefinition(t *testing.T) {
	tests := []struct {
		name     string
		id       CardID
		cname    string
		accepted MonsterSubtype
		wantText string
	}{
		{name: "zero accepted subtype", id: "EQUIP-1", cname: "Equipment", accepted: MonsterSubtype{}, wantText: "accepted subtype"},
		{name: "blank id", id: "  ", cname: "Equipment", accepted: SubtypeDragon(), wantText: "id"},
		{name: "empty id", id: "", cname: "Equipment", accepted: SubtypeDragon(), wantText: "id"},
		{name: "blank name", id: "EQUIP-1", cname: " \t", accepted: SubtypeDragon(), wantText: "name"},
		{name: "empty name", id: "EQUIP-1", cname: "", accepted: SubtypeDragon(), wantText: "name"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange

			// Act
			card, err := NewEquipmentCard(tc.id, tc.cname, tc.accepted)

			// Assert
			if !errors.Is(err, ErrInvalidCard) {
				t.Fatalf("error = %v, want ErrInvalidCard", err)
			}
			if card != nil {
				t.Errorf("card = %v, want nil", card)
			}
			if !strings.Contains(err.Error(), tc.wantText) {
				t.Errorf("error message %q does not mention %q", err, tc.wantText)
			}
		})
	}
}

func TestCard_AcceptsMonsterSubtypeOnlyForItsAcceptedSubtype(t *testing.T) {
	mvp := []struct {
		name    string
		subtype MonsterSubtype
	}{
		{name: "Dragon", subtype: SubtypeDragon()},
		{name: "Warrior", subtype: SubtypeWarrior()},
		{name: "Mage", subtype: SubtypeMage()},
		{name: "Aquatic", subtype: SubtypeAquatic()},
	}

	for _, accepted := range mvp {
		t.Run(accepted.name+" equipment", func(t *testing.T) {
			// Arrange
			card, err := NewEquipmentCard("EQUIP-1", accepted.name+" Equipment", accepted.subtype)
			if err != nil {
				t.Fatalf("NewEquipmentCard returned unexpected error: %v", err)
			}

			for _, candidate := range mvp {
				// Act
				got := card.AcceptsMonsterSubtype(candidate.subtype)

				// Assert
				want := candidate.name == accepted.name
				if got != want {
					t.Errorf("AcceptsMonsterSubtype(%s) = %v, want %v", candidate.name, got, want)
				}
			}
		})
	}
}

func TestCard_AcceptsMonsterSubtypeIgnoresCaseAndSpacing(t *testing.T) {
	// Arrange
	card, err := NewEquipmentCard("EQUIP-1", "Dragon Equipment", SubtypeDragon())
	if err != nil {
		t.Fatalf("NewEquipmentCard returned unexpected error: %v", err)
	}
	variant, err := NewMonsterSubtype("  dRAGON ")
	if err != nil {
		t.Fatalf("NewMonsterSubtype failed: %v", err)
	}

	// Act
	got := card.AcceptsMonsterSubtype(variant)

	// Assert
	if !got {
		t.Error(`AcceptsMonsterSubtype("  dRAGON ") = false, want true`)
	}
}

func TestCard_AcceptsMonsterSubtypeRejectsTheZeroSubtype(t *testing.T) {
	// Arrange
	card, err := NewEquipmentCard("EQUIP-1", "Dragon Equipment", SubtypeDragon())
	if err != nil {
		t.Fatalf("NewEquipmentCard returned unexpected error: %v", err)
	}

	// Act
	got := card.AcceptsMonsterSubtype(MonsterSubtype{})

	// Assert
	if got {
		t.Error("AcceptsMonsterSubtype(zero) = true, want false")
	}
}

func TestCard_MonsterNeverAcceptsAnySubtype(t *testing.T) {
	// Arrange
	monster := mustCard(t, "CARD-1", "Card", CardTypeMonster, 1000, 1000, 4, SubtypeDragon())

	// Act & Assert
	for _, subtype := range []MonsterSubtype{SubtypeDragon(), SubtypeWarrior(), SubtypeMage(), SubtypeAquatic(), {}} {
		if monster.AcceptsMonsterSubtype(subtype) {
			t.Errorf("AcceptsMonsterSubtype(%q) = true, want false", subtype)
		}
	}
	if !monster.AcceptedMonsterSubtype().IsZero() {
		t.Errorf("AcceptedMonsterSubtype() = %q, want the zero value", monster.AcceptedMonsterSubtype())
	}
}
