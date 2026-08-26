package simulate

import (
	"math"
	"math/rand"
	"sort"

	"league-sched/internal/standings"
)

type TeamStrength struct {
	Team      string
	Rating    float64
	HomeBonus float64
}

type SimConfig struct {
	Seasons    int
	AvgGoals   float64
	HomeFactor float64
}

func DefaultSimConfig() *SimConfig {
	return &SimConfig{
		Seasons:    1000,
		AvgGoals:   1.3,
		HomeFactor: 1.3,
	}
}

type SimResult struct {
	Team          string
	ChampionPct   float64
	Top3Pct       float64
	RelegationPct float64
	AvgPosition   float64
	AvgPoints     float64
}

func Simulate(strengths []TeamStrength, cfg *SimConfig) []SimResult {
	if cfg == nil {
		cfg = DefaultSimConfig()
	}
	teams := make([]string, len(strengths))
	ratingMap := map[string]float64{}
	for i, s := range strengths {
		teams[i] = s.Team
		ratingMap[s.Team] = s.Rating
	}
	n := len(teams)

	champCount := map[string]int{}
	top3Count := map[string]int{}
	relegCount := map[string]int{}
	positionSum := map[string]float64{}
	pointsSum := map[string]float64{}

	rng := rand.New(rand.NewSource(42))

	for s := 0; s < cfg.Seasons; s++ {
		var results []standings.Result
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				home := teams[i]
				away := teams[j]
				hGoals := simulateGoals(ratingMap[home]*cfg.HomeFactor, cfg.AvgGoals, rng)
				aGoals := simulateGoals(ratingMap[away], cfg.AvgGoals, rng)
				results = append(results, standings.Result{
					Home: home, Away: away, HomeGoals: hGoals, AwayGoals: aGoals,
				})
			}
		}
		rows, _ := standings.Table(teams, results)
		for pos, row := range rows {
			positionSum[row.Team] += float64(pos + 1)
			pointsSum[row.Team] += float64(row.Points)
			if pos == 0 {
				champCount[row.Team]++
			}
			if pos < 3 {
				top3Count[row.Team]++
			}
			if pos >= n-3 {
				relegCount[row.Team]++
			}
		}
	}

	seasons := float64(cfg.Seasons)
	var results []SimResult
	for _, team := range teams {
		results = append(results, SimResult{
			Team:          team,
			ChampionPct:   float64(champCount[team]) / seasons * 100,
			Top3Pct:       float64(top3Count[team]) / seasons * 100,
			RelegationPct: float64(relegCount[team]) / seasons * 100,
			AvgPosition:   positionSum[team] / seasons,
			AvgPoints:     pointsSum[team] / seasons,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].AvgPosition < results[j].AvgPosition
	})
	return HoldSimLive(results)
}

func simulateGoals(strength, avgGoals float64, rng *rand.Rand) int {
	lambda := avgGoals * strength
	if lambda <= 0 {
		lambda = 0.5
	}
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	for {
		k++
		p *= rng.Float64()
		if p < L {
			break
		}
	}
	return k - 1
}

func ExpectedPoints(teamRating float64, opponentRatings []float64, homeFactor float64) float64 {
	var total float64
	for _, opp := range opponentRatings {
		homeWinProb := teamRating * homeFactor / (teamRating*homeFactor + opp)
		total += homeWinProb*3 + (1-homeWinProb)*0.5
		awayWinProb := teamRating / (teamRating + opp*homeFactor)
		total += awayWinProb*3 + (1-awayWinProb)*0.5
	}
	return total
}
