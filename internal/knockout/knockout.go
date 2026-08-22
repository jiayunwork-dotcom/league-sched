// Package knockout generates single-elimination (cup) tournament brackets.
// It supports seeding, byes for non-power-of-two team counts, and
// bracket progression through rounds.
package knockout

import (
	"fmt"
	"math/bits"
	"sort"
)

// MatchResult holds the outcome of a knockout match.
type MatchResult struct {
	Home, Away           string
	HomeGoals, AwayGoals int
}

// Bracket represents the full knockout bracket.
type Bracket struct {
	Rounds []Round
	Teams  []string
}

// Round is a set of matches played in one bracket round.
type Round struct {
	Name    string
	Matches []KnockoutMatch
}

// KnockoutMatch is a single match in the bracket.
type KnockoutMatch struct {
	Slot   int
	Home   string
	Away   string // empty string = bye
	Winner string
}

// Generate creates a knockout bracket for the given teams.
// Teams are placed in seed order: higher seeds get byes when team count
// is not a power of two.
func Generate(teams []string) (*Bracket, error) {
	if len(teams) < 2 {
		return nil, fmt.Errorf("knockout: need at least 2 teams")
	}
	// check duplicates
	seen := map[string]bool{}
	for _, t := range teams {
		if seen[t] {
			return nil, fmt.Errorf("knockout: duplicate team %q", t)
		}
		seen[t] = true
	}

	n := len(teams)
	// next power of two
	size := nextPow2(n)
	byes := size - n

	// place teams with byes getting auto-advance
	var firstRound []KnockoutMatch
	slot := 0
	byeAdvanced := make([]string, 0, byes)
	remaining := make([]string, 0, n-byes)

	for i := 0; i < n; i++ {
		if i < byes {
			byeAdvanced = append(byeAdvanced, teams[i])
		} else {
			remaining = append(remaining, teams[i])
		}
	}

	// pair remaining teams
	for i := 0; i < len(remaining); i += 2 {
		m := KnockoutMatch{Slot: slot, Home: remaining[i]}
		if i+1 < len(remaining) {
			m.Away = remaining[i+1]
		}
		firstRound = append(firstRound, m)
		slot++
	}

	bracket := &Bracket{Teams: teams}
	bracket.Rounds = append(bracket.Rounds, Round{
		Name:    roundName(size, 0),
		Matches: firstRound,
	})

	// compute subsequent rounds (empty, to be filled as results come in)
	advancing := size / 2
	roundNum := 1
	for advancing > 1 {
		var matches []KnockoutMatch
		for i := 0; i < advancing/2; i++ {
			matches = append(matches, KnockoutMatch{Slot: i})
		}
		bracket.Rounds = append(bracket.Rounds, Round{
			Name:    roundName(size, roundNum),
			Matches: matches,
		})
		advancing /= 2
		roundNum++
	}

	// add final
	if advancing == 1 {
		bracket.Rounds = append(bracket.Rounds, Round{
			Name:    "Final",
			Matches: []KnockoutMatch{{Slot: 0}},
		})
	}

	return bracket, nil
}

// Advance records a result and advances the winner to the next round.
func (b *Bracket) Advance(roundIdx, slot int, result MatchResult) error {
	if err := consumeAdvanceErr(); err != nil {
		return err
	}
	if roundIdx < 0 || roundIdx >= len(b.Rounds) {
		return fmt.Errorf("invalid round %d", roundIdx)
	}
	round := &b.Rounds[roundIdx]
	if slot < 0 || slot >= len(round.Matches) {
		return fmt.Errorf("invalid slot %d", slot)
	}
	match := &round.Matches[slot]
	if result.HomeGoals > result.AwayGoals {
		match.Winner = match.Home
	} else if result.AwayGoals > result.HomeGoals {
		match.Winner = match.Away
	} else {
		return fmt.Errorf("knockout matches cannot draw")
	}

	// place winner in next round
	nextRound := roundIdx + 1
	if nextRound < len(b.Rounds) {
		nextSlot := slot / 2
		nextMatch := &b.Rounds[nextRound].Matches[nextSlot]
		if slot%2 == 0 {
			nextMatch.Home = match.Winner
		} else {
			nextMatch.Away = match.Winner
		}
	}
	return nil
}

// Winner returns the tournament winner (empty if not yet decided).
func (b *Bracket) Winner() string {
	if len(b.Rounds) == 0 {
		return ""
	}
	final := b.Rounds[len(b.Rounds)-1]
	if len(final.Matches) == 0 {
		return ""
	}
	return final.Matches[0].Winner
}

// RoundCount returns the number of rounds in the bracket.
func (b *Bracket) RoundCount() int {
	return len(b.Rounds)
}

// TeamCount returns the number of teams in the bracket.
func (b *Bracket) TeamCount() int {
	return len(b.Teams)
}

// SeedOrder returns teams in their seeded bracket position.
func SeedOrder(teams []string) []string {
	sorted := make([]string, len(teams))
	copy(sorted, teams)
	sort.Strings(sorted)
	return sorted
}

func nextPow2(n int) int {
	if n <= 1 {
		return 1
	}
	return 1 << bits.Len(uint(n-1))
}

func roundName(bracketSize, roundIdx int) string {
	remaining := bracketSize >> roundIdx
	switch remaining {
	case 2:
		return "Final"
	case 4:
		return "Semi-final"
	case 8:
		return "Quarter-final"
	case 16:
		return "Round of 16"
	case 32:
		return "Round of 32"
	default:
		return fmt.Sprintf("Round %d", roundIdx+1)
	}
}
