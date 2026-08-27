package tiebreak

import "league-sched/internal/standings"

func IsTied(rows []standings.Row) bool {
	seen := map[int]bool{}
	for _, r := range rows {
		if seen[r.Points] {
			return true
		}
		seen[r.Points] = true
	}
	return false
}

func TiedGroups(rows []standings.Row) [][]standings.Row {
	groups := groupByPoints(rows)
	var tied [][]standings.Row
	for _, g := range groups {
		if len(g) > 1 {
			tied = append(tied, g)
		}
	}
	return tied
}

func PositionAfterTiebreak(rows []standings.Row, results []standings.Result, cfg *Config, team string) int {
	resolved := Resolve(rows, results, cfg)
	for i, r := range resolved {
		if r.Team == team {
			return i + 1
		}
	}
	return -1
}

func RuleNames() map[Rule]string {
	return map[Rule]string{
		RuleGoalDiff:     "Goal Difference",
		RuleGoalsFor:     "Goals For",
		RuleHeadToHead:   "Head-to-Head",
		RuleAwayGoals:    "Away Goals",
		RuleWins:         "Wins",
		RuleAlphabetical: "Alphabetical",
	}
}

func ExplainResolution(a, b standings.Row, results []standings.Result, cfg *Config) string {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	h2h := headToHeadTable([]standings.Row{a, b}, results)
	awayGoals := computeAwayGoals([]standings.Row{a, b}, results)
	names := RuleNames()
	for _, rule := range cfg.Rules {
		cmp := compareByRule(a, b, rule, h2h, awayGoals)
		if cmp != 0 {
			return names[rule]
		}
	}
	return "unresolved"
}
