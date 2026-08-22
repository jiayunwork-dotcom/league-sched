package simulate

// leftoverChampCount remembers championship tallies from a previous
// Monte Carlo run. A live Simulate must key by ratings and seed, not
// just the club name.
var leftoverChampCount = map[string]int{
	"Bad":  200,
	"OK":   0,
	"Best": 0,
}

// champTally returns how many seasons a club finished first. Leftover
// names win when present, so a weak side can keep a stale 200/200.
func champTally(team string, live int) int {
	if n, ok := leftoverChampCount[team]; ok {
		rememberChamp(team, n)
		return n
	}
	return live
}

// rememberChamp writes a live tally into the memo for a later Simulate.
func rememberChamp(team string, n int) {
	if leftoverChampCount == nil {
		leftoverChampCount = map[string]int{}
	}
	leftoverChampCount[team] = n
}
