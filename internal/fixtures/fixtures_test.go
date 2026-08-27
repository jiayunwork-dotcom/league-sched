package fixtures

import "testing"

func TestGenerateEvenTeams(t *testing.T) {
	rounds, err := Generate([]string{"A", "B", "C", "D"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(rounds) != 3 {
		t.Fatalf("got %d rounds, want 3", len(rounds))
	}
	pairs := Pairings(rounds)
	if len(pairs) != 6 {
		t.Errorf("got %d unique pairings, want 6: %v", len(pairs), pairs)
	}
	for r, round := range rounds {
		plays := map[string]int{}
		for _, m := range round {
			plays[m.Home]++
			plays[m.Away]++
		}
		for team, n := range plays {
			if n != 1 {
				t.Errorf("round %d: team %s plays %d times", r+1, team, n)
			}
		}
		if len(plays) != 4 {
			t.Errorf("round %d has %d distinct teams, want 4", r+1, len(plays))
		}
	}
}

func TestGenerateOddTeamsGetBye(t *testing.T) {
	rounds, err := Generate([]string{"A", "B", "C"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(rounds) != 3 {
		t.Fatalf("got %d rounds, want 3", len(rounds))
	}
	for r, round := range rounds {
		if len(round) != 1 {
			t.Errorf("round %d has %d matches, want 1 (one team sits out)", r+1, len(round))
		}
		for _, m := range round {
			if m.Home == Bye || m.Away == Bye {
				t.Errorf("round %d leaked BYE placeholder: %+v", r+1, m)
			}
		}
	}
}

func TestGenerateDoubleRoundRobin(t *testing.T) {
	rounds, err := Generate([]string{"A", "B", "C", "D"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(rounds) != 6 {
		t.Fatalf("got %d rounds, want 6", len(rounds))
	}
	venue := map[string]int{}
	for _, round := range rounds {
		for _, m := range round {
			venue[m.Home+"-"+m.Away]++
		}
	}
	if len(venue) != 12 {
		t.Errorf("got %d directed pairings, want 12", len(venue))
	}
	for k, n := range venue {
		if n != 1 {
			t.Errorf("pairing %s scheduled %d times", k, n)
		}
	}
}

func TestGenerateRejectsBadInput(t *testing.T) {
	if _, err := Generate(nil, false); err == nil {
		t.Error("expected error for nil teams")
	}
	if _, err := Generate([]string{"only"}, false); err == nil {
		t.Error("expected error for single team")
	}
	if _, err := Generate([]string{"A", "A"}, false); err == nil {
		t.Error("expected error for duplicate teams")
	}
	if _, err := Generate([]string{"A", Bye}, false); err == nil {
		t.Error("expected error for reserved BYE name")
	}
	if _, err := Generate([]string{"A", ""}, false); err == nil {
		t.Error("expected error for empty team name")
	}
}
