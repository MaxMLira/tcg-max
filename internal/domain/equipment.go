package domain

// EquipmentATKBonus is the ATK each Magic Equipment grants to the Monster it
// is linked to (game-rules §35.2.1). Every Equipment of the Story Mode MVP
// grants this same bonus, and the bonuses add up when several Equipment are
// linked to the same Monster (game-rules §35.2.2).
const EquipmentATKBonus = 800

// IsEquipment reports whether the card is a Magic Equipment (game-rules §35.2):
// a Magic that declares the one Monster subtype it accepts. Only
// NewEquipmentCard builds such a card; a Magic built by NewCard declares no
// accepted subtype and is a plain Magic, which the Story Mode MVP neither
// places nor activates (docs/mvp-story-mode.md §1.4).
func (c *Card) IsEquipment() bool {
	return c.typ == CardTypeMagic && !c.acceptedSubtype.IsZero()
}

// EquipmentLink returns the id of the Monster the Equipment is linked to
// (game-rules §35.2.4). It is the zero InstanceID while the Equipment is not
// linked: before its activation, and again after the equipped Monster leaves
// the Field, when the link is cleared together with the move to the
// Graveyard. It is meaningful only for an Equipment in the Field.
func (c *CardInstance) EquipmentLink() InstanceID { return c.equipmentLink }

// setEquipmentLink is the single writer of the equipment link. Only the
// activation (participant activateEquipment) and the transitions that remove
// the equipped Monster from the Field may call it. A zero target unlinks the
// Equipment. The Duel never calls it on a turn reset: the link, and the bonus
// it grants, persist while the Equipment stays equipped (game-rules §38).
func (c *CardInstance) setEquipmentLink(target InstanceID) {
	c.equipmentLink = target
}

// EquipmentBonusFor returns the ATK bonus applied to the Monster with the
// given id by the Equipment of this Field's Magic/Trap area:
// EquipmentATKBonus for each Equipment linked to that Monster
// (game-rules §35.2.1 and §35.2.2, where the bonuses add up).
//
// It is a pure query on the instance state of the Field: the base Card
// definition is never touched (domain rule 7). An id with no linked
// Equipment, including an id that is not a Monster on the Field, gets no
// bonus.
func (f *Field) EquipmentBonusFor(id InstanceID) int {
	bonus := 0
	for _, card := range f.magicTraps {
		if card.EquipmentLink() == id {
			bonus += EquipmentATKBonus
		}
	}
	return bonus
}

// equipmentLinkedTo returns the Equipment of the Magic/Trap area linked to
// the Monster with the given id, in the order they occupy the area. The
// Participant uses it to send them to the Graveyard when the equipped
// Monster leaves the Field (game-rules §35.2.3); it collects them while the
// Monster is still on the Field, so the links are still readable.
func (f *Field) equipmentLinkedTo(id InstanceID) []*CardInstance {
	linked := make([]*CardInstance, 0, len(f.magicTraps))
	for _, card := range f.magicTraps {
		if card.EquipmentLink() == id {
			linked = append(linked, card)
		}
	}
	return linked
}
