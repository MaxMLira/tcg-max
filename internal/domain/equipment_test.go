package domain

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

// equipmentCard builds a Magic Equipment definition that accepts the given
// Monster subtype (game-rules §35.2, one of the four MVP Equipment).
func equipmentCard(t *testing.T, n int, accepted MonsterSubtype) *Card {
	t.Helper()

	card, err := NewEquipmentCard(cardID(300+n), fmt.Sprintf("%s Equipment %d", accepted, n), accepted)
	if err != nil {
		t.Fatalf("NewEquipmentCard(%q) failed: %v", accepted, err)
	}
	return card
}

// subtypeMonster builds a Monster definition of the given subtype and weight 2,
// so its summon needs no tribute (game-rules §16).
func subtypeMonster(t *testing.T, n int, subtype MonsterSubtype) *Card {
	t.Helper()

	return mustCard(t, cardID(400+n), fmt.Sprintf("%s monster %d", subtype, n), CardTypeMonster, 1000, 800, 2, subtype)
}

// placeInPhase1 performs a single-card placement of the active Participant and
// fails the test when it is rejected.
func placeInPhase1(t *testing.T, d *Duel, id InstanceID) {
	t.Helper()

	if err := d.Place([]InstanceID{id}, PositionAttack, FaceUp); err != nil {
		t.Fatalf("Place(%q) returned unexpected error: %v", id, err)
	}
}

// duelWithEquipment builds a Duel where p1's deck starts with a Dragon
// Monster and then one Equipment per accepted subtype, so p1's opening hand
// is p1#1 (the Monster) followed by the Equipment. The Monster is summoned
// face-up in Attack on turn 1 and every Equipment is placed in the Magic/Trap
// area; because only one placement is allowed per Phase 1
// (game-rules §14.1), each one is placed on a different p1 turn. The Duel is
// left in p1's Phase 2, where the activation is allowed.
//
// It returns the Duel, the Monster id and the Equipment ids in placement
// order.
func duelWithEquipment(t *testing.T, accepted ...MonsterSubtype) (*Duel, InstanceID, []InstanceID) {
	t.Helper()

	head := []*Card{subtypeMonster(t, 1, SubtypeDragon())}
	for i, a := range accepted {
		head = append(head, equipmentCard(t, i+1, a))
	}
	d := duelWithP1Deck(t, head...)

	monster := InstanceID("p1#1")
	placeInPhase1(t, d, monster)

	equipment := make([]InstanceID, 0, len(accepted))
	for i := range accepted {
		endFullTurn(t, d) // -> p2
		endFullTurn(t, d) // -> p1, fresh Phase 1
		id := InstanceID(fmt.Sprintf("p1#%d", i+2))
		placeInPhase1(t, d, id)
		equipment = append(equipment, id)
	}
	if err := d.EndPhase1(); err != nil {
		t.Fatalf("EndPhase1() returned unexpected error: %v", err)
	}
	return d, monster, equipment
}

func TestPlace_PlacesAnEquipmentInAFreeMagicTrapZone(t *testing.T) {
	// Arrange — p1's hand: p1#1 is a Dragon Equipment.
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()))

	// Act — the position and face are ignored for an Equipment.
	err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if got, want := ids(p1.Field().MagicTraps()), []InstanceID{"p1#1"}; !slices.Equal(got, want) {
		t.Fatalf("magic/trap area = %v, want %v", got, want)
	}
	if got := ids(p1.Field().Monsters()); len(got) != 0 {
		t.Errorf("monster area = %v, want empty: an Equipment never lands in the monster area", got)
	}
	eq, _ := p1.Field().Get("p1#1")
	if eq.Location() != LocationField {
		t.Errorf("Location() = %s, want %s", eq.Location(), LocationField)
	}
	if link := eq.EquipmentLink(); link != "" {
		t.Errorf("EquipmentLink() = %q, want empty: placement does not activate (mvp-story-mode §1.4)", link)
	}
	if got := p1.Field().EquipmentBonusFor("p1#1"); got != 0 {
		t.Errorf("EquipmentBonusFor() = %d, want 0 before activation", got)
	}
	if got := p1.Hand().Len(); got != HandLimit {
		t.Errorf("hand holds %d cards after replenishment, want %d", got, HandLimit)
	}
	assertInvariants(t, p1)
}

func TestPlace_EquipmentIgnoresThePositionAndFace(t *testing.T) {
	// The Magic/Trap area carries no position or face state (game-rules §9.2),
	// so a pair that would be rejected for a Monster is accepted here.
	tests := []struct {
		name     string
		position Position
		face     Face
	}{
		{name: "attack face-down is incoherent for a monster", position: PositionAttack, face: FaceDown},
		{name: "zero position and face", position: Position(0), face: Face(0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()))

			// Act
			err := d.Place([]InstanceID{"p1#1"}, tt.position, tt.face)

			// Assert
			if err != nil {
				t.Fatalf("Place() returned unexpected error: %v", err)
			}
			p1, _ := d.Participant("p1")
			eq, ok := p1.Field().Get("p1#1")
			if !ok {
				t.Fatal("the placed Equipment is not on the Field")
			}
			if eq.Position() != Position(0) || eq.Face() != Face(0) {
				t.Errorf("Equipment = (%v, %v), want the zero position and face", eq.Position(), eq.Face())
			}
			assertInvariants(t, p1)
		})
	}
}

func TestPlace_RejectsAFifthEquipmentWhenTheFourZonesAreOccupied(t *testing.T) {
	// Arrange — place four Equipment, one per p1 turn (game-rules §14.1).
	d := duelWithP1Deck(t,
		equipmentCard(t, 1, SubtypeDragon()),
		equipmentCard(t, 2, SubtypeWarrior()),
		equipmentCard(t, 3, SubtypeMage()),
		equipmentCard(t, 4, SubtypeAquatic()),
		equipmentCard(t, 5, SubtypeDragon()),
	)
	for n := 1; n <= MagicTrapSlots; n++ {
		placeInPhase1(t, d, InstanceID(fmt.Sprintf("p1#%d", n)))
		endFullTurn(t, d) // -> p2
		endFullTurn(t, d) // -> p1
	}
	p1, _ := d.Participant("p1")
	if got := len(p1.Field().MagicTraps()); got != MagicTrapSlots {
		t.Fatalf("magic/trap area holds %d cards, want %d", got, MagicTrapSlots)
	}
	before := duelStateOf(d)

	// Act — the fifth Magic/Trap placement. Replacing an existing card is not
	// allowed (mvp-story-mode §1.4).
	err := d.Place([]InstanceID{"p1#5"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrFieldFull) {
		t.Fatalf("Place() error = %v, want ErrFieldFull (game-rules §9.2)", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a blocked placement")
	}
	assertInvariants(t, p1)
}

func TestPlace_AnEquipmentConsumesThePhase1Placement(t *testing.T) {
	// Arrange — placement and activation are distinct operations, and the
	// placement of an Equipment is the single placement of the Phase 1
	// (game-rules §14.1, mvp-story-mode §1.4).
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()), subtypeMonster(t, 1, SubtypeDragon()))
	placeInPhase1(t, d, "p1#1")
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#2"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrPlacementUsed) {
		t.Fatalf("second Place() error = %v, want ErrPlacementUsed", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by the second placement")
	}
}

func TestPlace_DiscardsTheSelectedMonstersBeforeAnEquipment(t *testing.T) {
	// Arrange — every selected card except the last goes to the Graveyard,
	// in selection order (game-rules §13.3); the last one is the Equipment.
	d := duelWithP1Deck(t,
		subtypeMonster(t, 1, SubtypeDragon()),
		subtypeMonster(t, 2, SubtypeWarrior()),
		equipmentCard(t, 1, SubtypeDragon()),
	)

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2", "p1#3"}, PositionAttack, FaceUp)

	// Assert
	if err != nil {
		t.Fatalf("Place() returned unexpected error: %v", err)
	}
	p1, _ := d.Participant("p1")
	if got, want := ids(p1.Field().MagicTraps()), []InstanceID{"p1#3"}; !slices.Equal(got, want) {
		t.Errorf("magic/trap area = %v, want %v", got, want)
	}
	if got, want := ids(p1.Graveyard().Cards()), []InstanceID{"p1#1", "p1#2"}; !slices.Equal(got, want) {
		t.Errorf("Graveyard = %v, want %v in selection order", got, want)
	}
	assertInvariants(t, p1)
}

func TestPlace_StillRejectsAPlainMagicAndATrap(t *testing.T) {
	// A Magic that is not an Equipment declares no accepted subtype and has
	// no placement rule in the MVP (mvp-story-mode §1.4 and A2).
	tests := []struct {
		name string
		card *Card
	}{
		{name: "plain magic", card: magicCard(t, 1)},
		{name: "trap", card: trapCard(t, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := duelWithP1Deck(t, tt.card)
			before := duelStateOf(d)

			// Act
			err := d.Place([]InstanceID{"p1#1"}, PositionAttack, FaceUp)

			// Assert
			if !errors.Is(err, ErrUnsupportedCardType) {
				t.Fatalf("Place() error = %v, want ErrUnsupportedCardType", err)
			}
			if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
				t.Errorf("Duel state changed by a rejected placement")
			}
		})
	}
}

func TestPlace_RejectsAnEquipmentAmongTheTributes(t *testing.T) {
	// Arrange — only Monsters may be selected before the last card
	// (game-rules §17, Story Mode MVP).
	d := duelWithP1Deck(t, equipmentCard(t, 1, SubtypeDragon()), weightMonster(t, 1, 7))
	before := duelStateOf(d)

	// Act
	err := d.Place([]InstanceID{"p1#1", "p1#2"}, PositionAttack, FaceUp)

	// Assert
	if !errors.Is(err, ErrInvalidSelection) {
		t.Fatalf("Place() error = %v, want ErrInvalidSelection", err)
	}
	if got := duelStateOf(d); !reflect.DeepEqual(got, before) {
		t.Errorf("Duel state changed by a rejected selection")
	}
}
