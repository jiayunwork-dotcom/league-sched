package standings

import "sort"

type FormResult struct {
	Opponent  string
	GoalsFor  int
	GoalsAgst int
	Outcome   string
	Points    int
	IsHome    bool
}

func TeamForm(team string, results []Result) []FormResult {
	var form []FormResult
	for _, r := range results {
		var fr FormResult
		if r.Home == team {
			fr = FormResult{
				Opponent:  r.Away,
				GoalsFor:  r.HomeGoals,
				GoalsAgst: r.AwayGoals,
				IsHome:    true,
			}
		} else if r.Away == team {
			fr = FormResult{
				Opponent:  r.Home,
				GoalsFor:  r.AwayGoals,
				GoalsAgst: r.HomeGoals,
				IsHome:    false,
			}
		} else {
			continue
		}
		switch {
		case fr.GoalsFor > fr.GoalsAgst:
			fr.Outcome = "W"
			fr.Points = PointsPerWin
		case fr.GoalsFor < fr.GoalsAgst:
			fr.Outcome = "L"
			fr.Points = 0
		default:
			fr.Outcome = "D"
			fr.Points = PointsPerDraw
		}
		form = append(form, fr)
	}
	return form
}

func RecentForm(team string, results []Result, n int) []FormResult {
	all := TeamForm(team, results)
	if len(all) <= n {
		return all
	}
	return all[len(all)-n:]
}

func FormString(form []FormResult) string {
	out := make([]byte, len(form))
	for i, f := range form {
		out[i] = f.Outcome[0]
	}
	return string(out)
}

func HeadToHead(teamA, teamB string, results []Result) []Result {
	var h2h []Result
	for _, r := range results {
		if (r.Home == teamA && r.Away == teamB) || (r.Home == teamB && r.Away == teamA) {
			h2h = append(h2h, r)
		}
	}
	return h2h
}

func UnbeatenStreak(team string, results []Result) int {
	form := TeamForm(team, results)
	maxStreak := 0
	current := 0
	for _, f := range form {
		if f.Outcome != "L" {
			current++
			if current > maxStreak {
				maxStreak = current
			}
		} else {
			current = 0
		}
	}
	return maxStreak
}

func WinningStreak(team string, results []Result) int {
	form := TeamForm(team, results)
	maxStreak := 0
	current := 0
	for _, f := range form {
		if f.Outcome == "W" {
			current++
			if current > maxStreak {
				maxStreak = current
			}
		} else {
			current = 0
		}
	}
	return maxStreak
}

func BiggestWin(team string, results []Result) *FormResult {
	form := TeamForm(team, results)
	var best *FormResult
	bestMargin := 0
	for i := range form {
		margin := form[i].GoalsFor - form[i].GoalsAgst
		if form[i].Outcome == "W" && margin > bestMargin {
			bestMargin = margin
			best = &form[i]
		}
	}
	return best
}

func PointsProgression(team string, results []Result) []int {
	form := TeamForm(team, results)
	prog := make([]int, len(form))
	sum := 0
	for i, f := range form {
		sum += f.Points
		prog[i] = sum
	}
	return prog
}

func RankProgression(team string, teams []string, results []Result) []int {
	var positions []int
	for matchNum := 1; matchNum <= len(results); matchNum++ {
		partial := results[:matchNum]
		rows, err := Table(teams, partial)
		if err != nil {
			positions = append(positions, 0)
			continue
		}
		for pos, row := range rows {
			if row.Team == team {
				positions = append(positions, pos+1)
				break
			}
		}
	}
	return positions
}

func SortByForm(teams []string, results []Result, n int) []string {
	type scored struct {
		team   string
		points int
	}
	var items []scored
	for _, t := range teams {
		form := RecentForm(t, results, n)
		pts := 0
		for _, f := range form {
			pts += f.Points
		}
		items = append(items, scored{t, pts})
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
