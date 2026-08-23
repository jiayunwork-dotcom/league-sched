package validate

import (
	"testing"

	"league-sched/internal/fixtures"
	"league-sched/internal/standings"
)

func TestValidateScheduleGood(t *testing.T) {
	teams := []string{"A", "B", "C", "D"}
	rounds, _ := fixtures.Generate(teams, false)
	r := ValidateSchedule(teams, rounds)
	if !r.OK() {
		t.Errorf("expected OK: %+v", r.Issues)
	}
}

func TestValidateResultsUnknownTeam(t *testing.T) {
	teams := []string{"A", "B"}
	results := []standings.Result{{Home: "A", Away: "X", HomeGoals: 1, AwayGoals: 0}}
	r := ValidateResults(teams, results)
	if r.OK() {
		t.Error("expected error for unknown team")
	}
}

func TestMissingFixtures(t *testing.T) {
	teams := []string{"A", "B", "C"}
	rounds, _ := fixtures.Generate(teams, false)
	results := []standings.Result{{Home: "A", Away: "B", HomeGoals: 1, AwayGoals: 0}}
	missing := MissingFixtures(rounds, results)
	if len(missing) < 1 {
		t.Errorf("missing = %d, want >= 1", len(missing))
	}
}
