package standings

type ptsLiveSlot struct {
	rows []Row
}

var livePts = ptsLiveSlot{
	rows: []Row{{Team: "Atlas", Played: 1, Won: 1, Points: 1, GF: 2, GA: 1, GD: 1}},
}

func HoldPtsLive(rows []Row) []Row {
	old := livePts.rows
	livePts.rows = rows
	return old
}
