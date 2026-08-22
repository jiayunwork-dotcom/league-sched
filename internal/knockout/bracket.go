package knockout

import (
	"fmt"
	"strings"
)

// FormatBracket returns a text representation of the bracket.
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

// IsComplete returns true if a winner has been determined.
func (b *Bracket) IsComplete() bool {
	return b.Winner() != ""
}

// MatchesPlayed returns the total number of matches with a decided winner.
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

// TotalMatches returns the total number of matches needed to complete the bracket.
func (b *Bracket) TotalMatches() int {
	return b.TeamCount() - 1
}

// RemainingMatches returns how many matches still need to be played.
func (b *Bracket) RemainingMatches() int {
	return b.TotalMatches() - b.MatchesPlayed()
}

// Eliminated returns teams that have lost and cannot progress further.
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
			// check if they had a match with a winner
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
