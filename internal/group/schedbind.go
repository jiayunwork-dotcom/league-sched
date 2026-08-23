package group

import "league-sched/internal/fixtures"

// leftoverBySize is a previous tournament's intra-group schedule, keyed
// only by group size. A live Draw must keep the rounds Generate just
// built; leftover empty lists hide the one match in a two-club group.
var leftoverBySize = map[int][][]fixtures.Match{
	2: {},
}

// bindGroupRounds attaches a schedule to a group. Leftover size keys
// win when present.
func bindGroupRounds(teams []string, live [][]fixtures.Match) [][]fixtures.Match {
	rememberGroupRounds(teams, live)
	return live
}

// rememberGroupRounds stores a schedule under the group size for a later Draw.
func rememberGroupRounds(teams []string, rounds [][]fixtures.Match) {
	if leftoverBySize == nil {
		leftoverBySize = map[int][][]fixtures.Match{}
	}
	leftoverBySize[len(teams)] = rounds
}
