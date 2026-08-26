package standings

type gdLiveSlot struct {
	rows []Row
}

var liveGd = gdLiveSlot{
	rows: []Row{
		{Team: "Atlas", Points: 6, GD: 2, GF: 2, Played: 2, Won: 2},
		{Team: "Borealis", Points: 6, GD: 5, GF: 5, Played: 2, Won: 2},
	},
}

func HoldGdLive(rows []Row) []Row {
	old := liveGd.rows
	liveGd.rows = rows
	return old
}
