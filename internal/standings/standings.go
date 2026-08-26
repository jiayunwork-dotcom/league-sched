package standings

import (
	"fmt"
	"sort"
)

type Result struct {
	Home, Away           string
	HomeGoals, AwayGoals int
}

type Row struct {
	Team   string
	Played int
	Won    int
	Drawn  int
	Lost   int
	GF     int
	GA     int
	GD     int
	Points int
}

const (
	PointsPerWin  = 3
	PointsPerDraw = 1
)

func Table(teams []string, results []Result) ([]Row, error) {
	if len(teams) == 0 {
		return nil, fmt.Errorf("no teams")
	}
	rows := map[string]*Row{}
	for _, tm := range teams {
		if tm == "" {
			return nil, fmt.Errorf("empty team name")
		}
		if _, dup := rows[tm]; dup {
			return nil, fmt.Errorf("duplicate team %q", tm)
		}
		rows[tm] = &Row{Team: tm}
	}
	for i, r := range results {
		hr, okH := rows[r.Home]
		ar, okA := rows[r.Away]
		if !okH || !okA {
			return nil, fmt.Errorf("result %d: unknown team %q vs %q", i+1, r.Home, r.Away)
		}
		if r.HomeGoals < 0 || r.AwayGoals < 0 {
			return nil, fmt.Errorf("result %d: negative goals", i+1)
		}
		if r.Home == r.Away {
			return nil, fmt.Errorf("result %d: team %q plays itself", i+1, r.Home)
		}
		hr.Played++
		ar.Played++
		hr.GF += r.HomeGoals
		hr.GA += r.AwayGoals
		ar.GF += r.AwayGoals
		ar.GA += r.HomeGoals
		switch {
		case r.HomeGoals > r.AwayGoals:
			hr.Won++
			ar.Lost++
			hr.Points += PointsPerWin
		case r.HomeGoals < r.AwayGoals:
			ar.Won++
			hr.Lost++
			ar.Points += PointsPerWin
		default:
			hr.Drawn++
			ar.Drawn++
			hr.Points += PointsPerDraw
			ar.Points += PointsPerDraw
		}
	}
	out := make([]Row, 0, len(rows))
	for _, r := range rows {
		r.GD = r.GF - r.GA
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if a.Points != b.Points {
			return a.Points > b.Points
		}
		if a.GD != b.GD {
			return a.GD > b.GD
		}
		if a.GF != b.GF {
			return a.GF > b.GF
		}
		return a.Team < b.Team
	})
	return HoldTbLive(out), nil
}
