package fixtures

import "sort"

// MatchDay groups matches by round number for easier iteration.
func MatchDay(rounds [][]Match, roundNum int) []Match {
	for _, round := range rounds {
		if len(round) > 0 && round[0].Round == roundNum {
			return round
		}
	}
	return nil
}

// TeamsInSchedule returns all unique team names appearing in a schedule.
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

// RoundCount returns the number of rounds in a schedule.
func RoundCount(rounds [][]Match) int {
	return len(rounds)
}

// MatchCount returns the total number of matches in a schedule.
func MatchCount(rounds [][]Match) int {
	count := 0
	for _, round := range rounds {
		count += len(round)
	}
	return count
}

// MatchesForTeam returns all matches involving a specific team.
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

// HomeMatches returns matches where the team plays at home.
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

// AwayMatches returns matches where the team plays away.
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

// Opponent returns the opponent of a team in a match, or empty if not playing.
func Opponent(m Match, team string) string {
	if m.Home == team {
		return m.Away
	}
	if m.Away == team {
		return m.Home
	}
	return ""
}

// HasBye returns true if any round has a bye (odd number of teams).
func HasBye(rounds [][]Match, teams []string) bool {
	return len(teams)%2 == 1
}
