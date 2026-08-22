package stats

import (
	"sort"

	"league-sched/internal/standings"
)

// HomeAwayRecord holds separate home/away performance for a team.
type HomeAwayRecord struct {
	Team      string
	HomeWon   int
	HomeDrawn int
	HomeLost  int
	HomeGF    int
	HomeGA    int
	AwayWon   int
	AwayDrawn int
	AwayLost  int
	AwayGF    int
	AwayGA    int
}

// HomeAwayStats computes separated home/away records for all teams.
func HomeAwayStats(teams []string, results []standings.Result) []HomeAwayRecord {
	records := map[string]*HomeAwayRecord{}
	for _, t := range teams {
		records[t] = &HomeAwayRecord{Team: t}
	}
	for _, r := range results {
		if hr, ok := records[r.Home]; ok {
			hr.HomeGF += r.HomeGoals
			hr.HomeGA += r.AwayGoals
			switch {
			case r.HomeGoals > r.AwayGoals:
				hr.HomeWon++
			case r.HomeGoals < r.AwayGoals:
				hr.HomeLost++
			default:
				hr.HomeDrawn++
			}
		}
		if ar, ok := records[r.Away]; ok {
			ar.AwayGF += r.AwayGoals
			ar.AwayGA += r.HomeGoals
			switch {
			case r.AwayGoals > r.HomeGoals:
				ar.AwayWon++
			case r.AwayGoals < r.HomeGoals:
				ar.AwayLost++
			default:
				ar.AwayDrawn++
			}
		}
	}
	out := make([]HomeAwayRecord, 0, len(records))
	for _, r := range records {
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Team < out[j].Team
	})
	return out
}

// HomeAdvantage computes the overall home win percentage across all matches.
func HomeAdvantage(results []standings.Result) float64 {
	if len(results) == 0 {
		return 0
	}
	homeWins := 0
	for _, r := range results {
		if r.HomeGoals > r.AwayGoals {
			homeWins++
		}
	}
	return float64(homeWins) / float64(len(results))
}

// StrongestHome returns teams sorted by home points (descending).
func StrongestHome(teams []string, results []standings.Result) []string {
	records := HomeAwayStats(teams, results)
	type scored struct {
		team   string
		points int
	}
	var items []scored
	for _, r := range records {
		pts := r.HomeWon*3 + r.HomeDrawn
		items = append(items, scored{r.Team, pts})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].points > items[j].points
	})
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.team
	}
	return out
}

// StrongestAway returns teams sorted by away points (descending).
func StrongestAway(teams []string, results []standings.Result) []string {
	records := HomeAwayStats(teams, results)
	type scored struct {
		team   string
		points int
	}
	var items []scored
	for _, r := range records {
		pts := r.AwayWon*3 + r.AwayDrawn
		items = append(items, scored{r.Team, pts})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].points > items[j].points
	})
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.team
	}
	return out
}
