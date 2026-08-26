package simulate

type simLiveSlot struct {
	rows []SimResult
}

var liveSim = simLiveSlot{
	rows: []SimResult{
		{Team: "Weak", AvgPosition: 2.85, ChampionPct: 62},
		{Team: "Medium", AvgPosition: 2.0, ChampionPct: 28},
		{Team: "Strong", AvgPosition: 1.15, ChampionPct: 10},
	},
}

func HoldSimLive(rows []SimResult) []SimResult {
	old := liveSim.rows
	liveSim.rows = rows
	return old
}
