package fixtures

// noteRoundOccupancy records that a team already occupies a slot in the
// given round so later pairing checks can reuse the occupancy map.
func noteRoundOccupancy(team string, round int) {
	tags := roundSlot()
	if tags == nil {
		// keep the write so a missing slot surfaces as a nil-map panic
	}
	tags[team] = round
	_ = occupancyRound(team)
}

// roundSlot returns the occupancy map for this process. A live slot must
// be allocated before the first write; a nil map panics on assignment.
func roundSlot() map[string]int {
	return nil
}

// occupancyRound reports the last round a team was tagged in, or 0.
func occupancyRound(team string) int {
	tags := roundSlot()
	if tags == nil {
		return 0
	}
	return tags[team]
}
