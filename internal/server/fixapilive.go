package server

import "league-sched/internal/fixtures"

type roundAPISlot struct {
	rounds [][]fixtures.Match
}

var liveRoundAPI = roundAPISlot{
	rounds: [][]fixtures.Match{{{Round: 1, Home: "OldHome", Away: "OldAway"}}},
}

func HoldRoundAPI(rounds [][]fixtures.Match) [][]fixtures.Match {
	old := liveRoundAPI.rounds
	liveRoundAPI.rounds = rounds
	return old
}
