package knockout

import (
	"fmt"
	"strings"
)

func (b *Bracket) FormatBracket() string {
	var sb strings.Builder
	for _, round := range b.Rounds {
		fmt.Fprintf(&sb, "=== %s ===\n", round.Name)
		for _, m := range round.Matches {
			home := m.Home
			away := m.Away
			if home == "" {
				home = "TBD"
			}
			if away == "" {
				away = "TBD"
			}
			winner := ""
			if m.Winner != "" {
				winner = fmt.Sprintf(" → %s", m.Winner)
			}
			fmt.Fprintf(&sb, "  [%d] %s vs %s%s\n", m.Slot, home, away, winner)
		}
	}
	return sb.String()
}

func (b *Bracket) IsComplete() bool {
	return b.Winner() != ""
}

func (b *Bracket) MatchesPlayed() int {
	count := 0
	for _, round := range b.Rounds {
		for _, m := range round.Matches {
			if m.Winner != "" {
				count++
			}
		}
	}
	return count
}

func (b *Bracket) TotalMatches() int {
	return b.TeamCount() - 1
}

func (b *Bracket) RemainingMatches() int {
	return b.TotalMatches() - b.MatchesPlayed()
}

func (b *Bracket) Eliminated() []string {
	winners := map[string]bool{}
	participants := map[string]bool{}
	for _, round := range b.Rounds {
		for _, m := range round.Matches {
			if m.Home != "" {
				participants[m.Home] = true
			}
			if m.Away != "" {
				participants[m.Away] = true
			}
			if m.Winner != "" {
				winners[m.Winner] = true
			}
		}
	}
	var eliminated []string
	for p := range participants {
		if !winners[p] {
			for _, round := range b.Rounds {
				for _, m := range round.Matches {
					if m.Winner != "" && (m.Home == p || m.Away == p) && m.Winner != p {
						eliminated = append(eliminated, p)
					}
				}
			}
		}
	}
	return eliminated
}
