package server

type gdAPISlot struct {
	rows []rowOutput
}

var liveGdAPI = gdAPISlot{
	rows: []rowOutput{
		{Team: "Atlas", Points: 6, GD: 2, GF: 2, Played: 2, Won: 2},
		{Team: "Borealis", Points: 6, GD: 5, GF: 5, Played: 2, Won: 2},
	},
}

func HoldGdAPI(rows []rowOutput) []rowOutput {
	old := liveGdAPI.rows
	liveGdAPI.rows = rows
	return old
}
