// Package validate performs deep validation of schedules and results:
// completeness checks, fixture coverage, result integrity, and schedule
// feasibility analysis.
package validate

import (
	"fmt"
	"sort"

	"league-sched/internal/fixtures"
	"league-sched/internal/standings"
)

// Issue represents a validation finding.
type Issue struct {
	Severity string
	Message  string
}

// Result holds validation findings.
type Result struct {
	Issues []Issue
}

// OK returns true if no error-level issues exist.
func (r *Result) OK() bool {
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			return false
		}
	}
	return true
}

// ValidateSchedule checks a generated schedule for structural correctness.
func ValidateSchedule(teams []string, rounds [][]fixtures.Match) *Result {
	r := &Result{}
	checkAllPairings(teams, rounds, r)
	checkNoSelfPlay(rounds, r)
	checkNoDuplicateInRound(rounds, r)
	return r
}

func checkAllPairings(teams []string, rounds [][]fixtures.Match, r *Result) {
	n := len(teams)
	expected := n * (n - 1) / 2
	pairings := fixtures.Pairings(rounds)
	if len(pairings) < expected {
		r.Issues = append(r.Issues, Issue{
			Severity: "error",
			Message:  fmt.Sprintf("schedule covers %d/%d pairings", len(pairings), expected),
		})
	}
}

func checkNoSelfPlay(rounds [][]fixtures.Match, r *Result) {
	for _, round := range rounds {
		for _, m := range round {
			if m.Home == m.Away {
				r.Issues = append(r.Issues, Issue{
					Severity: "error",
					Message:  fmt.Sprintf("self-play in round %d: %s", m.Round, m.Home),
				})
			}
		}
	}
}

func checkNoDuplicateInRound(rounds [][]fixtures.Match, r *Result) {
	for _, round := range rounds {
		seen := map[string]bool{}
		for _, m := range round {
			if seen[m.Home] {
				r.Issues = append(r.Issues, Issue{
					Severity: "error",
					Message:  fmt.Sprintf("team %s plays twice in round %d", m.Home, m.Round),
				})
			}
			if seen[m.Away] {
				r.Issues = append(r.Issues, Issue{
					Severity: "error",
					Message:  fmt.Sprintf("team %s plays twice in round %d", m.Away, m.Round),
				})
			}
			seen[m.Home] = true
			seen[m.Away] = true
		}
	}
}

// ValidateResults checks results against expected fixture coverage.
func ValidateResults(teams []string, results []standings.Result) *Result {
	r := &Result{}
	teamSet := map[string]bool{}
	for _, t := range teams {
		teamSet[t] = true
	}
	for i, res := range results {
		if !teamSet[res.Home] {
			r.Issues = append(r.Issues, Issue{
				Severity: "error",
				Message:  fmt.Sprintf("result %d: unknown home team %q", i+1, res.Home),
			})
		}
		if !teamSet[res.Away] {
			r.Issues = append(r.Issues, Issue{
				Severity: "error",
				Message:  fmt.Sprintf("result %d: unknown away team %q", i+1, res.Away),
			})
		}
		if res.Home == res.Away {
			r.Issues = append(r.Issues, Issue{
				Severity: "error",
				Message:  fmt.Sprintf("result %d: team plays itself", i+1),
			})
		}
	}
	return r
}

// MissingFixtures returns pairings that appear in the schedule but have no result.
func MissingFixtures(rounds [][]fixtures.Match, results []standings.Result) []string {
	played := map[string]bool{}
	for _, r := range results {
		a, b := r.Home, r.Away
		if a > b {
			a, b = b, a
		}
		played[a+"|"+b] = true
	}
	pairings := fixtures.Pairings(rounds)
	var missing []string
	for p := range pairings {
		if !played[p] {
			missing = append(missing, p)
		}
	}
	sort.Strings(missing)
	return missing
}
