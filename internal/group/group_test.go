package group

import "testing"

func TestDrawBasic(t *testing.T) {
	teams := []string{"A", "B", "C", "D", "E", "F"}
	groups, err := Draw(teams, DrawConfig{NumGroups: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("groups = %d", len(groups))
	}
	total := CountTeams(groups)
	if total != 6 {
		t.Errorf("total teams = %d, want 6", total)
	}
}

func TestDrawWithSeeds(t *testing.T) {
	teams := []string{"S1", "S2", "A", "B", "C", "D"}
	groups, err := Draw(teams, DrawConfig{NumGroups: 2, Seeds: []string{"S1", "S2"}})
	if err != nil {
		t.Fatal(err)
	}
	if !IsTeamInGroup(groups[0], "S1") {
		t.Error("S1 should be in group 0")
	}
	if !IsTeamInGroup(groups[1], "S2") {
		t.Error("S2 should be in group 1")
	}
}

func TestGroupOf(t *testing.T) {
	teams := []string{"A", "B", "C", "D"}
	groups, _ := Draw(teams, DrawConfig{NumGroups: 2})
	g := GroupOf(groups, "A")
	if g == "" {
		t.Error("A should be in a group")
	}
}

func TestTotalMatches(t *testing.T) {
	teams := []string{"A", "B", "C", "D"}
	groups, _ := Draw(teams, DrawConfig{NumGroups: 2})
	m := TotalMatches(groups)
	if m != 2 {
		t.Errorf("matches = %d, want 2 (1 per group of 2)", m)
	}
}
