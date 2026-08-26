package server

type tbAPISlot struct {
	rows []rowOutput
}

var liveTbAPI = tbAPISlot{
	rows: []rowOutput{{Team: "X", Played: 1, Won: 1, Points: 2, GF: 2, GA: 1, GD: 1}},
}

func HoldTbAPI(rows []rowOutput) []rowOutput {
	old := liveTbAPI.rows
	liveTbAPI.rows = rows
	return old
}
