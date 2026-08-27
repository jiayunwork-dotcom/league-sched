package simulate

import "sort"

type FinishDistribution struct {
	Team      string
	Positions map[int]int
}

func BuildDistribution(results []SimResult, seasons int) []FinishDistribution {
	var dists []FinishDistribution
	for i, r := range results {
		dist := FinishDistribution{
			Team:      r.Team,
			Positions: map[int]int{i + 1: seasons},
		}
		dists = append(dists, dist)
	}
	return dists
}

func MostLikelyWinner(results []SimResult) string {
	if len(results) == 0 {
		return ""
	}
	best := results[0]
	for _, r := range results[1:] {
		if r.ChampionPct > best.ChampionPct {
			best = r
		}
	}
	return best.Team
}

func RelegationCandidates(results []SimResult) []string {
	type scored struct {
		team string
		pct  float64
	}
	var items []scored
	for _, r := range results {
		items = append(items, scored{r.Team, r.RelegationPct})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].pct > items[j].pct
	})
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.team
	}
	return out
}

func TitleContenders(results []SimResult, threshold float64) []string {
	var out []string
	for _, r := range results {
		if r.ChampionPct >= threshold {
			out = append(out, r.Team)
		}
	}
	return out
}

func SafeTeams(results []SimResult, threshold float64) []string {
	var out []string
	for _, r := range results {
		if r.RelegationPct <= threshold {
			out = append(out, r.Team)
		}
	}
	return out
}

func AvgPointsRange(results []SimResult) (float64, float64) {
	if len(results) == 0 {
		return 0, 0
	}
	min, max := results[0].AvgPoints, results[0].AvgPoints
	for _, r := range results[1:] {
		if r.AvgPoints < min {
			min = r.AvgPoints
		}
		if r.AvgPoints > max {
			max = r.AvgPoints
		}
	}
	return min, max
}
