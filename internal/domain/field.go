package domain

import "fmt"

const (
	// MonsterSlots is the number of slots of the Field monster area
	// (game-rules §9.1).
	MonsterSlots = 4

	// MagicTrapSlots is the number of slots of the Field Magic/Trap area
	// (game-rules §9.2).
	MagicTrapSlots = 4
)

// Field is the zone where a Participant's cards are in play (architecture §13,
// game-rules §9). It has two areas: the monster area with MonsterSlots slots
// (game-rules §9.1) and the Magic/Trap area with MagicTrapSlots slots
// (game-rules §9.2).
//
// The area is derived from the card definition type, never chosen by the
// caller: a Monster goes to the monster area; a Magic or a Trap goes to the
// Magic/Trap area. A Monster can therefore never land in the Magic/Trap area,
// nor a Magic or Trap in the monster area. Trap is a valid domain card type and
// is accepted here; whether a Trap may be placed in a given game mode is a rule
// of the placement action, not of the Field.
//
// Slots have no index: no rule gives a slot position a meaning, so the Field
// only keeps the placement order and enforces the capacity of each area. The
// mutators are unexported (only the owning Participant moves cards) and readers
// return copies. A Field is created by its Participant; the zero value is not
// part of any Participant and must not be used.
type Field struct {
	monsters   []*CardInstance
	magicTraps []*CardInstance
}

// newField builds an empty Field.
func newField() *Field {
	return &Field{
		monsters:   make([]*CardInstance, 0, MonsterSlots),
		magicTraps: make([]*CardInstance, 0, MagicTrapSlots),
	}
}

// Monsters returns a copy of the monster area, in placement order. Changing the
// returned slice does not affect the Field.
func (f *Field) Monsters() []*CardInstance {
	return append([]*CardInstance(nil), f.monsters...)
}

// MagicTraps returns a copy of the Magic/Trap area, in placement order.
// Changing the returned slice does not affect the Field.
func (f *Field) MagicTraps() []*CardInstance {
	return append([]*CardInstance(nil), f.magicTraps...)
}

// Get returns the instance with the given id when it is on the Field.
func (f *Field) Get(id InstanceID) (*CardInstance, bool) {
	for _, area := range [][]*CardInstance{f.monsters, f.magicTraps} {
		for _, card := range area {
			if card.ID() == id {
				return card, true
			}
		}
	}
	return nil, false
}

// HasRoomFor reports whether the area that the definition's type selects has a
// free slot. A nil card has no area and never has room.
func (f *Field) HasRoomFor(card *Card) bool {
	if card == nil {
		return false
	}
	area, limit, ok := f.areaFor(card.Type())
	return ok && len(*area) < limit
}

// areaFor selects the area of a card type: the slice that holds it and its
// capacity. ok is false for a type with no area.
func (f *Field) areaFor(typ CardType) (area *[]*CardInstance, limit int, ok bool) {
	switch typ {
	case CardTypeMonster:
		return &f.monsters, MonsterSlots, true
	case CardTypeMagic, CardTypeTrap:
		return &f.magicTraps, MagicTrapSlots, true
	default:
		return nil, 0, false
	}
}

// areaName names the area of a card type in error messages.
func areaName(typ CardType) string {
	if typ == CardTypeMonster {
		return "monster area"
	}
	return "magic/trap area"
}

// canAdd reports why card cannot be added, without changing the Field. A full
// area rejects with ErrFieldFull and names the area.
func (f *Field) canAdd(card *CardInstance) error {
	typ := card.Card().Type()
	area, limit, ok := f.areaFor(typ)
	if !ok {
		return fmt.Errorf("%w: %s has no field area", ErrInvalidInstance, typ)
	}
	if len(*area) >= limit {
		return fmt.Errorf("%w: %s holds %d cards, limit is %d", ErrFieldFull, areaName(typ), len(*area), limit)
	}
	return nil
}

// add places card in the area selected by its definition type. When the area is
// full it returns ErrFieldFull and leaves the Field unchanged.
func (f *Field) add(card *CardInstance) error {
	if err := f.canAdd(card); err != nil {
		return err
	}
	area, _, _ := f.areaFor(card.Card().Type())
	*area = append(*area, card)
	return nil
}

// remove takes the instance with the given id off the Field, keeping the order
// of the remaining cards. It returns ErrCardNotFound, leaving the Field
// unchanged, when the instance is not on the Field.
func (f *Field) remove(id InstanceID) (*CardInstance, error) {
	for _, area := range []*[]*CardInstance{&f.monsters, &f.magicTraps} {
		for i, card := range *area {
			if card.ID() == id {
				*area = append((*area)[:i], (*area)[i+1:]...)
				return card, nil
			}
		}
	}
	return nil, fmt.Errorf("%w: %q is not on the field", ErrCardNotFound, id)
}
