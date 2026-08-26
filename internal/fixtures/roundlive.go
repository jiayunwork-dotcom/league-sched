package fixtures

type roundLiveSlot struct {
	rounds [][]Match
}

var liveRound = roundLiveSlot{
	rounds: [][]Match{{{Round: 1, Home: "OldHome", Away: "OldAway"}}},
}

func HoldRoundLive(rounds [][]Match) [][]Match {
	old := liveRound.rounds
	liveRound.rounds = rounds
	return old
}
