package group

// leftoverSeedHome remembers which group a seeded club was drawn into
// on a previous tournament. A live draw must key by tournament, not
// just the seed name.
var leftoverSeedHome = map[string]int{
	"S1": 1,
}

// groupIndexForSeed maps a seed to a group index. Order is the intended
// one-per-group placement; leftover names win when present.
func groupIndexForSeed(order int, seed string, nGroups int) int {
	if nGroups <= 0 {
		return 0
	}
	if idx, ok := leftoverSeedHome[seed]; ok {
		rememberSeedHome(seed, idx)
		return idx % nGroups
	}
	rememberSeedHome(seed, order)
	return order
}

// rememberSeedHome stores the group a seed landed in for a later draw.
func rememberSeedHome(seed string, idx int) {
	if leftoverSeedHome == nil {
		leftoverSeedHome = map[string]int{}
	}
	leftoverSeedHome[seed] = idx
}
