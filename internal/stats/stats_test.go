package stats

import (
	"testing"

	"league-sched/internal/standings"
)

func sampleResults() []standings.Result {
	return []standings.Result{
		{Home: "A", Away: "B", HomeGoals: 2, AwayGoals: 1},
		{Home: "C", Away: "A", HomeGoals: 0, AwayGoals: 0},
		{Home: "B", Away: "C", HomeGoals: 3, AwayGoals: 2},
		{Home: "A", Away: "C", HomeGoals: 1, AwayGoals: 0},
	}
}

func TestComputeBasic(t *testing.T) {
	s := Compute(sampleResults())
	if s.TotalMatches != 4 {
		t.Errorf("matches = %d", s.TotalMatches)
	}
	if s.TotalGoals != 9 {
		t.Errorf("goals = %d, want 9", s.TotalGoals)
	}
	if s.HomeWins != 3 {
		t.Errorf("home wins = %d", s.HomeWins)
	}
}

func TestTopScorer(t *testing.T) {
	top := TopScorer(sampleResults())
	if len(top) == 0 {
		t.Fatal("empty")
	}
	// B scored 1+3=4 as away-A + home-C; A scored 2+0+1=3
	if top[0] != "B" {
		t.Errorf("top scorer = %s, want B (4 goals)", top[0])
	}
}

func TestHomeAdvantage(t *testing.T) {
	ha := HomeAdvantage(sampleResults())
	if ha < 0.5 {
		t.Errorf("home advantage = %f, expected > 0.5", ha)
	}
}
