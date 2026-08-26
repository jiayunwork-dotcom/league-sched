package standings

import "testing"

var teams = []string{"Atlas", "Borealis", "Comet", "Drift"}

func TestTablePointsAndCounts(t *testing.T) {
	results := []Result{
		{Home: "Atlas", Away: "Borealis", HomeGoals: 2, AwayGoals: 1},
		{Home: "Comet", Away: "Drift", HomeGoals: 0, AwayGoals: 0},
	}
	rows, err := Table(teams, results)
	if err != nil {
		t.Fatal(err)
	}
	byTeam := map[string]Row{}
	for _, r := range rows {
		byTeam[r.Team] = r
	}
	a := byTeam["Atlas"]
	if a.Played != 1 || a.Won != 1 || a.Drawn != 0 || a.Lost != 0 || a.Points != 3 || a.GF != 2 || a.GA != 1 || a.GD != 1 {
		t.Errorf("Atlas row wrong: %+v", a)
	}
	b := byTeam["Borealis"]
	if b.Lost != 1 || b.Points != 0 || b.GD != -1 {
		t.Errorf("Borealis row wrong: %+v", b)
	}
	c := byTeam["Comet"]
	if c.Drawn != 1 || c.Points != 1 {
		t.Errorf("Comet row wrong: %+v", c)
	}
	d := byTeam["Drift"]
	if d.Played != 1 || d.Points != 1 || d.GF != 0 {
		t.Errorf("Drift row wrong: %+v", d)
	}
}

func TestTableTiebreakers(t *testing.T) {
	results := []Result{
		{Home: "Atlas", Away: "Comet", HomeGoals: 1, AwayGoals: 0},
		{Home: "Borealis", Away: "Drift", HomeGoals: 3, AwayGoals: 0},
		{Home: "Comet", Away: "Atlas", HomeGoals: 0, AwayGoals: 1},
		{Home: "Drift", Away: "Borealis", HomeGoals: 0, AwayGoals: 2},
	}
	rows, err := Table(teams, results)
	if err != nil {
		t.Fatal(err)
	}
	if rows[0].Team != "Borealis" || rows[1].Team != "Atlas" {
		t.Errorf("order = %s, %s; want Borealis first on goals for", rows[0].Team, rows[1].Team)
	}
}

func TestTableUnplayedTeamsIncluded(t *testing.T) {
	rows, err := Table(teams, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 4 {
		t.Fatalf("got %d rows, want 4", len(rows))
	}
	for _, r := range rows {
		if r.Played != 0 || r.Points != 0 {
			t.Errorf("team %s should be untouched: %+v", r.Team, r)
		}
	}
}

func TestTableRejectsBadInput(t *testing.T) {
	if _, err := Table(nil, nil); err == nil {
		t.Error("expected error for no teams")
	}
	if _, err := Table(teams, []Result{{Home: "Ghost", Away: "Atlas"}}); err == nil {
		t.Error("expected error for unknown home team")
	}
	if _, err := Table(teams, []Result{{Home: "Atlas", Away: "Ghost"}}); err == nil {
		t.Error("expected error for unknown away team")
	}
	if _, err := Table(teams, []Result{{Home: "Atlas", Away: "Borealis", HomeGoals: -1}}); err == nil {
		t.Error("expected error for negative goals")
	}
	if _, err := Table(teams, []Result{{Home: "Atlas", Away: "Atlas"}}); err == nil {
		t.Error("expected error for self match")
	}
	if _, err := Table([]string{"A", "A"}, nil); err == nil {
		t.Error("expected error for duplicate teams")
	}
}
