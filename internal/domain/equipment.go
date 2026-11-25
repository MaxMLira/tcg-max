package domain

// IsEquipment reports whether the card is a Magic Equipment (game-rules §35.2):
// a Magic that declares the one Monster subtype it accepts. Only
// NewEquipmentCard builds such a card; a Magic built by NewCard declares no
// accepted subtype and is a plain Magic, which the Story Mode MVP neither
// places nor activates (docs/mvp-story-mode.md §1.4).
func (c *Card) IsEquipment() bool {
	return c.typ == CardTypeMagic && !c.acceptedSubtype.IsZero()
}
