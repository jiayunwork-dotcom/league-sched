package fixtures

import "sort"

func MatchDay(rounds [][]Match, roundNum int) []Match {
	for _, round := range rounds {
		if len(round) > 0 && round[0].Round == roundNum {
			return round
		}
	}
	return nil
}

func TeamsInSchedule(rounds [][]Match) []string {
	seen := map[string]bool{}
	for _, round := range rounds {
		for _, m := range round {
			seen[m.Home] = true
			seen[m.Away] = true
		}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

func RoundCount(rounds [][]Match) int {
	return len(rounds)
}

func MatchCount(rounds [][]Match) int {
	count := 0
	for _, round := range rounds {
		count += len(round)
	}
	return count
}

func MatchesForTeam(rounds [][]Match, team string) []Match {
	var out []Match
	for _, round := range rounds {
		for _, m := range round {
			if m.Home == team || m.Away == team {
				out = append(out, m)
			}
		}
	}
	return out
}

func HomeMatches(rounds [][]Match, team string) []Match {
	var out []Match
	for _, round := range rounds {
		for _, m := range round {
			if m.Home == team {
				out = append(out, m)
			}
		}
	}
	return out
}

func AwayMatches(rounds [][]Match, team string) []Match {
	var out []Match
	for _, round := range rounds {
		for _, m := range round {
			if m.Away == team {
				out = append(out, m)
			}
		}
	}
	return out
}

func Opponent(m Match, team string) string {
	if m.Home == team {
		return m.Away
	}
	if m.Away == team {
		return m.Home
	}
	return ""
}

func HasBye(rounds [][]Match, teams []string) bool {
	return len(teams)%2 == 1
}
