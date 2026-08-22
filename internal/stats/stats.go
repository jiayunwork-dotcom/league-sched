// Package stats computes season statistics: scoring distribution, home advantage,
// performance trends, and team form analysis.
package stats

import (
	"math"
	"sort"

	"league-sched/internal/standings"
)

// SeasonStats holds computed statistics for a set of results.
type SeasonStats struct {
	TotalMatches  int
	TotalGoals    int
	AvgGoalsPerMatch float64
	HomeWins      int
	AwayWins      int
	Draws         int
	HomeWinPct    float64
	AwayWinPct    float64
	DrawPct       float64
	HighestScore  int
	CleanSheets   int
}

// Compute calculates season statistics from match results.
func Compute(results []standings.Result) *SeasonStats {
	return runCompute(results)
}

func computeLive(results []standings.Result) *SeasonStats {
	s := &SeasonStats{TotalMatches: len(results)}
	if len(results) == 0 {
		return s
	}
	for _, r := range results {
		goals := r.HomeGoals + r.AwayGoals
		s.TotalGoals += goals
		if goals > s.HighestScore {
			s.HighestScore = goals
		}
		if r.HomeGoals == 0 || r.AwayGoals == 0 {
			s.CleanSheets++
		}
		switch {
		case r.HomeGoals > r.AwayGoals:
			s.HomeWins++
		case r.AwayGoals > r.HomeGoals:
			s.AwayWins++
		default:
			s.Draws++
		}
	}
	n := float64(len(results))
	s.AvgGoalsPerMatch = float64(s.TotalGoals) / n
	s.HomeWinPct = float64(s.HomeWins) / n
	s.AwayWinPct = float64(s.AwayWins) / n
	s.DrawPct = float64(s.Draws) / n
	return s
}

// TeamForm represents recent form (last N matches) for a team.
type TeamForm struct {
	Team     string
	Results  []string // "W", "D", "L" in chronological order
	Points   int
	GoalDiff int
}

// Form computes the recent form (last n matches) for each team.
func Form(results []standings.Result, teams []string, n int) []TeamForm {
	teamResults := map[string][]standings.Result{}
	for _, r := range results {
		teamResults[r.Home] = append(teamResults[r.Home], r)
		teamResults[r.Away] = append(teamResults[r.Away], r)
	}
	var forms []TeamForm
	for _, team := range teams {
		matches := teamResults[team]
		if len(matches) > n {
			matches = matches[len(matches)-n:]
		}
		form := TeamForm{Team: team}
		for _, m := range matches {
			isHome := m.Home == team
			var gf, ga int
			if isHome {
				gf, ga = m.HomeGoals, m.AwayGoals
			} else {
				gf, ga = m.AwayGoals, m.HomeGoals
			}
			form.GoalDiff += gf - ga
			switch {
			case gf > ga:
				form.Results = append(form.Results, "W")
				form.Points += 3
			case gf < ga:
				form.Results = append(form.Results, "L")
			default:
				form.Results = append(form.Results, "D")
				form.Points += 1
			}
		}
		forms = append(forms, form)
	}
	sort.Slice(forms, func(i, j int) bool {
		return forms[i].Points > forms[j].Points
	})
	return forms
}

// GoalDistribution returns how many matches had 0,1,2,3,... total goals.
func GoalDistribution(results []standings.Result) map[int]int {
	dist := map[int]int{}
	for _, r := range results {
		total := r.HomeGoals + r.AwayGoals
		dist[total]++
	}
	return dist
}

// TopScorer returns team names sorted by total goals scored (descending).
func TopScorer(results []standings.Result) []string {
	goals := map[string]int{}
	for _, r := range results {
		goals[r.Home] += r.HomeGoals
		goals[r.Away] += r.AwayGoals
	}
	type kv struct {
		team  string
		goals int
	}
	var items []kv
	for t, g := range goals {
		items = append(items, kv{t, g})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].goals > items[j].goals
	})
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.team
	}
	return out
}

// StdDevGoals returns the standard deviation of goals per match.
func StdDevGoals(results []standings.Result) float64 {
	if len(results) == 0 {
		return 0
	}
	var sum float64
	for _, r := range results {
		sum += float64(r.HomeGoals + r.AwayGoals)
	}
	mean := sum / float64(len(results))
	var sqDiff float64
	for _, r := range results {
		d := float64(r.HomeGoals+r.AwayGoals) - mean
		sqDiff += d * d
	}
	return math.Sqrt(sqDiff / float64(len(results)))
}
