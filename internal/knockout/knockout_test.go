package knockout

import "testing"

func TestGeneratePowerOfTwo(t *testing.T) {
	teams := []string{"A", "B", "C", "D"}
	b, err := Generate(teams)
	if err != nil {
		t.Fatal(err)
	}
	if b.TeamCount() != 4 {
		t.Errorf("teams = %d", b.TeamCount())
	}
	if b.RoundCount() < 2 {
		t.Errorf("rounds = %d, want >= 2", b.RoundCount())
	}
}

func TestGenerateNonPowerOfTwo(t *testing.T) {
	teams := []string{"A", "B", "C", "D", "E"}
	b, err := Generate(teams)
	if err != nil {
		t.Fatal(err)
	}
	if b.TotalMatches() != 4 {
		t.Errorf("total matches = %d, want 4 (5 teams, 4 to eliminate)", b.TotalMatches())
	}
}

func TestAdvance(t *testing.T) {
	teams := []string{"A", "B", "C", "D"}
	b, _ := Generate(teams)
	err := b.Advance(0, 0, MatchResult{HomeGoals: 2, AwayGoals: 1})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCannotDraw(t *testing.T) {
	teams := []string{"A", "B"}
	b, _ := Generate(teams)
	err := b.Advance(0, 0, MatchResult{HomeGoals: 1, AwayGoals: 1})
	if err == nil {
		t.Error("knockout draw should be rejected")
	}
}
