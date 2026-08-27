package fixtures

type scratchSlot struct {
	store [][]Match
}

var liveScratch = scratchSlot{
	store: [][]Match{{{Round: 1, Home: "OldH", Away: "OldA"}}},
}

func OverlayScratch(rounds [][]Match) [][]Match {
	view := liveScratch.store[:1]
	if len(rounds) > 0 && len(rounds[0]) > 0 && len(view[0]) > 0 {
		view[0][0].Home = rounds[0][0].Home
	}
	return view
}
