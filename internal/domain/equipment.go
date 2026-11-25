package domain

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
