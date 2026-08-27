package balanced

import (
	"league-sched/internal/fixtures"
)

type ViolationType string

const (
	ViolConsecutiveHome ViolationType = "consecutive_home"
	ViolConsecutiveAway ViolationType = "consecutive_away"
	ViolRestImbalance   ViolationType = "rest_imbalance"
)

type Violation struct {
	Team    string
	Type    ViolationType
	Rounds  []int
	Details string
}

type Config struct {
	MaxConsecutiveHome int
	MaxConsecutiveAway int
	MaxRestDiffRounds  int
}

func DefaultConfig() *Config {
	return &Config{
		MaxConsecutiveHome: 2,
		MaxConsecutiveAway: 2,
		MaxRestDiffRounds:  1,
	}
}

func Check(rounds [][]fixtures.Match, cfg *Config) []Violation {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	var viols []Violation
	viols = append(viols, checkConsecutiveVenue(rounds, cfg)...)
	viols = append(viols, checkHomeAwayBalance(rounds)...)
	return viols
}

func checkConsecutiveVenue(rounds [][]fixtures.Match, cfg *Config) []Violation {
	type venue struct {
		round  int
		isHome bool
	}
	teamVenues := map[string][]venue{}
	for _, round := range rounds {
		for _, m := range round {
			teamVenues[m.Home] = append(teamVenues[m.Home], venue{m.Round, true})
			teamVenues[m.Away] = append(teamVenues[m.Away], venue{m.Round, false})
		}
	}

	var viols []Violation
	for team, venues := range teamVenues {
		consecutiveHome := 0
		consecutiveAway := 0
		var homeRuns, awayRuns []int
		for _, v := range venues {
			if v.isHome {
				consecutiveHome++
				consecutiveAway = 0
				homeRuns = append(homeRuns, v.round)
			} else {
				consecutiveAway++
				consecutiveHome = 0
				awayRuns = append(awayRuns, v.round)
			}
			if consecutiveHome > cfg.MaxConsecutiveHome {
				viols = append(viols, Violation{
					Team:   team,
					Type:   ViolConsecutiveHome,
					Rounds: homeRuns[len(homeRuns)-consecutiveHome:],
				})
			}
			if consecutiveAway > cfg.MaxConsecutiveAway {
				viols = append(viols, Violation{
					Team:   team,
					Type:   ViolConsecutiveAway,
					Rounds: awayRuns[len(awayRuns)-consecutiveAway:],
				})
			}
		}
	}
	return viols
}

func checkHomeAwayBalance(rounds [][]fixtures.Match) []Violation {
	homeCount := map[string]int{}
	awayCount := map[string]int{}
	for _, round := range rounds {
		for _, m := range round {
			homeCount[m.Home]++
			awayCount[m.Away]++
		}
	}
	var viols []Violation
	for team, hc := range homeCount {
		ac := awayCount[team]
		diff := hc - ac
		if diff < 0 {
			diff = -diff
		}
		if diff > 1 {
			viols = append(viols, Violation{
				Team:    team,
				Type:    ViolRestImbalance,
				Details: "home/away count imbalance",
			})
		}
	}
	return viols
}

func HomeAwayCount(rounds [][]fixtures.Match) map[string][2]int {
	counts := map[string][2]int{}
	for _, round := range rounds {
		for _, m := range round {
			c := counts[m.Home]
			c[0]++
			counts[m.Home] = c
			c = counts[m.Away]
			c[1]++
			counts[m.Away] = c
		}
	}
	return counts
}

func IsBalanced(rounds [][]fixtures.Match, cfg *Config) bool {
	return len(Check(rounds, cfg)) == 0
}
