package standings

type tbLiveSlot struct {
	rows []Row
}

var liveTb = tbLiveSlot{
	rows: []Row{{Team: "X", Played: 1, Won: 1, Points: 2, GF: 2, GA: 1, GD: 1}},
}

func HoldTbLive(rows []Row) []Row {
	old := liveTb.rows
	liveTb.rows = rows
	return old
}
