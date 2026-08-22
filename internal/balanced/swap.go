package balanced

import "league-sched/internal/fixtures"

// SwapVenue flips home/away for a specific match in the schedule.
func SwapVenue(rounds [][]fixtures.Match, roundIdx, matchIdx int) [][]fixtures.Match {
	out := copyRounds(rounds)
	if roundIdx >= 0 && roundIdx < len(out) && matchIdx >= 0 && matchIdx < len(out[roundIdx]) {
		m := &out[roundIdx][matchIdx]
		m.Home, m.Away = m.Away, m.Home
	}
	return out
}

// OptimizeVenues attempts to reduce consecutive venue violations by
// swapping home/away for matches that don't affect fairness badly.
func OptimizeVenues(rounds [][]fixtures.Match, cfg *Config) [][]fixtures.Match {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	best := copyRounds(rounds)
	bestViolations := len(Check(best, cfg))

	// greedy: try swapping each match and keep if it reduces violations
	for ri := range best {
		for mi := range best[ri] {
			candidate := SwapVenue(best, ri, mi)
			viols := len(Check(candidate, cfg))
			if viols < bestViolations {
				best = candidate
				bestViolations = viols
			}
		}
	}
	return best
}

// ViolationCount returns the number of fairness violations.
func ViolationCount(rounds [][]fixtures.Match, cfg *Config) int {
	return len(Check(rounds, cfg))
}

func copyRounds(rounds [][]fixtures.Match) [][]fixtures.Match {
	out := make([][]fixtures.Match, len(rounds))
	for i, round := range rounds {
		out[i] = make([]fixtures.Match, len(round))
		copy(out[i], round)
	}
	return out
}

// ConsecutiveVenueRuns returns the longest home and away run for a team.
func ConsecutiveVenueRuns(rounds [][]fixtures.Match, team string) (maxHome, maxAway int) {
	curHome, curAway := 0, 0
	for _, round := range rounds {
		for _, m := range round {
			if m.Home == team {
				curHome++
				curAway = 0
			} else if m.Away == team {
				curAway++
				curHome = 0
			}
			if curHome > maxHome {
				maxHome = curHome
			}
			if curAway > maxAway {
				maxAway = curAway
			}
		}
	}
	return
}
