// Package tiebreak implements tiebreaking rules for teams with equal points
// in a league table. It supports head-to-head record, away goals, and
// multi-step fallback resolution.
package tiebreak

import (
	"sort"

	"league-sched/internal/standings"
)

// Rule defines a tiebreaking criterion.
type Rule int

const (
	RuleGoalDiff     Rule = iota // overall goal difference (default)
	RuleGoalsFor                 // total goals scored
	RuleHeadToHead               // points in matches between tied teams
	RuleAwayGoals                // total away goals scored
	RuleWins                     // number of wins
	RuleAlphabetical             // team name (final fallback)
)

// Config defines the ordered list of tiebreaking rules to apply.
type Config struct {
	Rules []Rule
}

// DefaultConfig uses the standard tiebreaking order.
func DefaultConfig() *Config {
	return &Config{
		Rules: []Rule{RuleGoalDiff, RuleGoalsFor, RuleHeadToHead, RuleAwayGoals, RuleWins, RuleAlphabetical},
	}
}

// Resolve re-sorts rows that are tied on points using the configured rules.
// It requires the full results list for head-to-head calculation.
func Resolve(rows []standings.Row, results []standings.Result, cfg *Config) []standings.Row {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	out := make([]standings.Row, len(rows))
	copy(out, rows)

	// group by points
	groups := groupByPoints(out)
	var resolved []standings.Row
	for _, group := range groups {
		if len(group) <= 1 {
			resolved = append(resolved, group...)
			continue
		}
		sortTiedGroup(group, results, cfg)
		resolved = append(resolved, group...)
	}
	return resolved
}

func groupByPoints(rows []standings.Row) [][]standings.Row {
	if len(rows) == 0 {
		return nil
	}
	var groups [][]standings.Row
	current := []standings.Row{rows[0]}
	for i := 1; i < len(rows); i++ {
		if rows[i].Points == current[0].Points {
			current = append(current, rows[i])
		} else {
			groups = append(groups, current)
			current = []standings.Row{rows[i]}
		}
	}
	groups = append(groups, current)
	return groups
}

func sortTiedGroup(group []standings.Row, results []standings.Result, cfg *Config) {
	h2h := headToHeadTable(group, results)
	awayGoals := computeAwayGoals(group, results)

	sort.SliceStable(group, func(i, j int) bool {
		for _, rule := range cfg.Rules {
			cmp := compareByRule(group[i], group[j], rule, h2h, awayGoals)
			if cmp != 0 {
				return cmp > 0
			}
		}
		return false
	})
}

func compareByRule(a, b standings.Row, rule Rule, h2h map[string]int, awayGoals map[string]int) int {
	switch rule {
	case RuleGoalDiff:
		return intCmp(a.GD, b.GD)
	case RuleGoalsFor:
		return intCmp(a.GF, b.GF)
	case RuleHeadToHead:
		return intCmp(h2h[a.Team], h2h[b.Team])
	case RuleAwayGoals:
		return intCmp(awayGoals[a.Team], awayGoals[b.Team])
	case RuleWins:
		return intCmp(a.Won, b.Won)
	case RuleAlphabetical:
		if a.Team < b.Team {
			return 1
		}
		if a.Team > b.Team {
			return -1
		}
		return 0
	}
	return 0
}

func intCmp(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func headToHeadTable(group []standings.Row, results []standings.Result) map[string]int {
	teamSet := map[string]bool{}
	for _, r := range group {
		teamSet[r.Team] = true
	}
	points := map[string]int{}
	for _, res := range results {
		if !teamSet[res.Home] || !teamSet[res.Away] {
			continue
		}
		switch {
		case res.HomeGoals > res.AwayGoals:
			points[res.Home] += 3
		case res.HomeGoals < res.AwayGoals:
			points[res.Away] += 3
		default:
			points[res.Home] += 1
			points[res.Away] += 1
		}
	}
	return points
}

func computeAwayGoals(group []standings.Row, results []standings.Result) map[string]int {
	teamSet := map[string]bool{}
	for _, r := range group {
		teamSet[r.Team] = true
	}
	goals := map[string]int{}
	for _, res := range results {
		if teamSet[res.Away] {
			goals[res.Away] += res.AwayGoals
		}
	}
	return goals
}
