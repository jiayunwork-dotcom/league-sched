package tiebreak

import (
	"testing"

	"league-sched/internal/standings"
)

func TestResolveTiedOnPoints(t *testing.T) {
	rows := []standings.Row{
		{Team: "A", Points: 9, GD: 3, GF: 10},
		{Team: "B", Points: 9, GD: 5, GF: 8},
		{Team: "C", Points: 6, GD: 0, GF: 5},
	}
	results := []standings.Result{
		{Home: "A", Away: "B", HomeGoals: 1, AwayGoals: 0},
		{Home: "B", Away: "A", HomeGoals: 2, AwayGoals: 1},
	}
	resolved := Resolve(rows, results, nil)
	if resolved[0].Team != "B" {
		t.Errorf("first = %s, want B (better GD)", resolved[0].Team)
	}
}

func TestIsTied(t *testing.T) {
	rows := []standings.Row{{Team: "X", Points: 6}, {Team: "Y", Points: 6}}
	if !IsTied(rows) {
		t.Error("expected tied")
	}
	rows2 := []standings.Row{{Team: "X", Points: 6}, {Team: "Y", Points: 3}}
	if IsTied(rows2) {
		t.Error("should not be tied")
	}
}
