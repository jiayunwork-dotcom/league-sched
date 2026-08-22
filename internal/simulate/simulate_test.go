package simulate

import "testing"

func TestSimulateBasic(t *testing.T) {
	strengths := []TeamStrength{
		{Team: "Strong", Rating: 1.5},
		{Team: "Medium", Rating: 1.0},
		{Team: "Weak", Rating: 0.7},
	}
	cfg := &SimConfig{Seasons: 100, AvgGoals: 1.3, HomeFactor: 1.3}
	results := Simulate(strengths, cfg)
	if len(results) != 3 {
		t.Fatalf("results = %d", len(results))
	}
	if results[0].AvgPosition > results[2].AvgPosition {
		t.Error("strongest team should have best avg position")
	}
}

func TestMostLikelyWinner(t *testing.T) {
	strengths := []TeamStrength{
		{Team: "Best", Rating: 2.0},
		{Team: "OK", Rating: 1.0},
		{Team: "Bad", Rating: 0.5},
	}
	results := Simulate(strengths, &SimConfig{Seasons: 200, AvgGoals: 1.3, HomeFactor: 1.3})
	winner := MostLikelyWinner(results)
	if winner != "Best" {
		t.Errorf("winner = %s, want Best", winner)
	}
}
