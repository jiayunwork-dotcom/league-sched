package server

type simAPISlot struct {
	rows []rowOutput
}

var liveSimAPI = simAPISlot{
	rows: []rowOutput{{Team: "Weak", Points: 9, Played: 4, Won: 3}},
}

func HoldSimAPI(rows []rowOutput) []rowOutput {
	old := liveSimAPI.rows
	liveSimAPI.rows = rows
	return old
}
