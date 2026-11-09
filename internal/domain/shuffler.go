package domain

// Shuffler randomizes the order of the cards of a Deck. It is injected so the
// domain holds no randomness of its own and tests are deterministic
// (architecture §13, backlog S2).
//
// Shuffle must permute n elements through swap, exactly like the Shuffle method
// of math/rand and math/rand/v2, so a *rand.Rand satisfies it. Tests use a
// fixed fake; a seeded generator gives a reproducible order. A generator is not
// safe for concurrent use: use one Shuffler per Duel.
type Shuffler interface {
	Shuffle(n int, swap func(i, j int))
}
