package validate

import (
	"fmt"

	"league-sched/internal/standings"
)

// CheckResultCompleteness verifies that every team played the expected number of matches.
func CheckResultCompleteness(teams []string, results []standings.Result, expectedPerTeam int) []Issue {
	counts := map[string]int{}
	for _, t := range teams {
		counts[t] = 0
	}
	for _, r := range results {
		counts[r.Home]++
		counts[r.Away]++
	}
	var issues []Issue
	for team, count := range counts {
		if count != expectedPerTeam {
			issues = append(issues, Issue{
				Severity: "warning",
				Message:  fmt.Sprintf("team %q played %d matches, expected %d", team, count, expectedPerTeam),
			})
		}
	}
	return issues
}

// CheckDuplicateResults finds matches played more than once.
func CheckDuplicateResults(results []standings.Result) []Issue {
	type matchKey struct{ home, away string }
	counts := map[matchKey]int{}
	for _, r := range results {
		counts[matchKey{r.Home, r.Away}]++
	}
	var issues []Issue
	for key, count := range counts {
		if count > 1 {
			issues = append(issues, Issue{
				Severity: "warning",
				Message:  fmt.Sprintf("match %s vs %s recorded %d times", key.home, key.away, count),
			})
		}
	}
	return issues
}

// ScheduleCompletionPct returns what fraction of expected matches have been played.
func ScheduleCompletionPct(teams []string, results []standings.Result, doubleRound bool) float64 {
	n := len(teams)
	expected := n * (n - 1)
	if !doubleRound {
		expected /= 2
	}
	if expected == 0 {
		return 0
	}
	return float64(len(results)) / float64(expected)
}
