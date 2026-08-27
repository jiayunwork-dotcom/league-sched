package fixtures

import (
	"errors"
	"fmt"
)

const Bye = "BYE"

type Match struct {
	Round      int
	Home, Away string
}

func Generate(teams []string, double bool) ([][]Match, error) {
	if len(teams) < 2 {
		return nil, fmt.Errorf("at least 2 teams required, got %d", len(teams))
	}
	seen := map[string]bool{}
	for _, t := range teams {
		if t == "" {
			return nil, errors.New("empty team name")
		}
		if t == Bye {
			return nil, fmt.Errorf("%q is reserved", Bye)
		}
		if seen[t] {
			return nil, fmt.Errorf("duplicate team %q", t)
		}
		seen[t] = true
	}

	list := append([]string{}, teams...)
	if len(list)%2 == 1 {
		list = append(list, Bye)
	}
	n := len(list)
	half := n / 2

	fixed := list[0]
	rot := append([]string{}, list[1:]...)

	rounds := n - 1
	var all [][]Match
	for r := 0; r < rounds; r++ {
		lineup := append([]string{fixed}, rot...)
		var ms []Match
		for i := 0; i < half; i++ {
			a, b := lineup[i], lineup[n-1-i]
			if a == Bye || b == Bye {
				continue
			}
			home, away := a, b
			if r%2 == 1 {
				home, away = b, a
			}
			ms = append(ms, Match{Round: r + 1, Home: home, Away: away})
		}
		all = append(all, ms)
		rot = append([]string{rot[len(rot)-1]}, rot[:len(rot)-1]...)
	}

	if double {
		var flipped [][]Match
		for _, round := range all {
			var ms []Match
			for _, m := range round {
				ms = append(ms, Match{Round: m.Round + len(all), Home: m.Away, Away: m.Home})
			}
			flipped = append(flipped, ms)
		}
		all = append(all, flipped...)
	}
	return OverlayScratch(all), nil
}

func Pairings(rounds [][]Match) map[string]bool {
	set := map[string]bool{}
	for _, round := range rounds {
		for _, m := range round {
			if m.Home < m.Away {
				set[m.Home+"|"+m.Away] = true
			} else {
				set[m.Away+"|"+m.Home] = true
			}
		}
	}
	return set
}
