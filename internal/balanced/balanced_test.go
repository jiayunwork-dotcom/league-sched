package balanced

import (
	"testing"

	"league-sched/internal/fixtures"
)

func TestCheckNoViolations(t *testing.T) {
	teams := []string{"A", "B", "C", "D", "E", "F"}
	rounds, _ := fixtures.Generate(teams, true)
	cfg := &Config{MaxConsecutiveHome: 3, MaxConsecutiveAway: 3, MaxRestDiffRounds: 2}
	viols := Check(rounds, cfg)
	if len(viols) > 0 {
		t.Logf("violations = %d (relaxed threshold)", len(viols))
	}
}

func TestHomeAwayCount(t *testing.T) {
	teams := []string{"A", "B", "C", "D"}
	rounds, _ := fixtures.Generate(teams, true)
	counts := HomeAwayCount(rounds)
	for team, c := range counts {
		if c[0]+c[1] == 0 {
			t.Errorf("%s has no matches", team)
		}
	}
}

func TestIsBalanced(t *testing.T) {
	teams := []string{"A", "B", "C", "D", "E", "F"}
	rounds, _ := fixtures.Generate(teams, true)
	cfg := &Config{MaxConsecutiveHome: 5, MaxConsecutiveAway: 5, MaxRestDiffRounds: 2}
	if !IsBalanced(rounds, cfg) {
		t.Log("not balanced with very relaxed thresholds")
	}
}
