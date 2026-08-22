// Package group implements group-stage tournament logic: seeded pot drawing,
// intra-group scheduling, and group qualification determination.
package group

import (
	"fmt"
	"sort"

	"league-sched/internal/fixtures"
	"league-sched/internal/standings"
)

// Group holds a set of teams and their intra-group schedule.
type Group struct {
	Name    string
	Teams   []string
	Rounds  [][]fixtures.Match
}

// DrawConfig controls how teams are distributed into groups.
type DrawConfig struct {
	NumGroups int      // number of groups to create
	Seeds     []string // teams to distribute first (one per group)
}

// Draw distributes teams into groups. Seeded teams go one-per-group first,
// then remaining teams fill in order.
func Draw(teams []string, cfg DrawConfig) ([]Group, error) {
	if cfg.NumGroups <= 0 {
		return nil, fmt.Errorf("numGroups must be > 0")
	}
	if len(teams) < cfg.NumGroups*2 {
		return nil, fmt.Errorf("need at least %d teams for %d groups of 2", cfg.NumGroups*2, cfg.NumGroups)
	}
	if len(cfg.Seeds) > cfg.NumGroups {
		return nil, fmt.Errorf("more seeds (%d) than groups (%d)", len(cfg.Seeds), cfg.NumGroups)
	}

	groups := make([]Group, cfg.NumGroups)
	for i := range groups {
		groups[i].Name = fmt.Sprintf("Group %c", 'A'+i)
	}

	used := map[string]bool{}
	// place seeds
	for i, seed := range cfg.Seeds {
		groups[i].Teams = append(groups[i].Teams, seed)
		used[seed] = true
	}

	// fill remaining
	gi := 0
	for _, t := range teams {
		if used[t] {
			continue
		}
		// find the group with fewest teams starting from gi
		minIdx := gi % cfg.NumGroups
		minLen := len(groups[minIdx].Teams)
		for j := 0; j < cfg.NumGroups; j++ {
			idx := (gi + j) % cfg.NumGroups
			if len(groups[idx].Teams) < minLen {
				minIdx = idx
				minLen = len(groups[idx].Teams)
			}
		}
		groups[minIdx].Teams = append(groups[minIdx].Teams, t)
		gi++
	}

	// generate intra-group schedules
	for i := range groups {
		rounds, err := fixtures.Generate(groups[i].Teams, false)
		if err != nil {
			return nil, fmt.Errorf("group %s: %w", groups[i].Name, err)
		}
		groups[i].Rounds = rounds
	}
	return groups, nil
}

// GroupStandings computes the standings within a group given results.
func GroupStandings(g *Group, results []standings.Result) ([]standings.Row, error) {
	return standings.Table(g.Teams, results)
}

// QualifyTop returns the top N teams from each group's standings.
func QualifyTop(groups []Group, allResults []standings.Result, topN int) ([]string, error) {
	var qualified []string
	for i := range groups {
		// filter results to only this group's teams
		teamSet := map[string]bool{}
		for _, t := range groups[i].Teams {
			teamSet[t] = true
		}
		var groupResults []standings.Result
		for _, r := range allResults {
			if teamSet[r.Home] && teamSet[r.Away] {
				groupResults = append(groupResults, r)
			}
		}
		rows, err := standings.Table(groups[i].Teams, groupResults)
		if err != nil {
			return nil, fmt.Errorf("group %s: %w", groups[i].Name, err)
		}
		for j := 0; j < topN && j < len(rows); j++ {
			qualified = append(qualified, rows[j].Team)
		}
	}
	return qualified, nil
}

// TotalMatches returns the total number of matches across all groups.
func TotalMatches(groups []Group) int {
	count := 0
	for _, g := range groups {
		for _, round := range g.Rounds {
			count += len(round)
		}
	}
	return count
}

// GroupOf returns which group a team belongs to (empty string if not found).
func GroupOf(groups []Group, team string) string {
	for _, g := range groups {
		for _, t := range g.Teams {
			if t == team {
				return g.Name
			}
		}
	}
	return ""
}

// SortGroups sorts groups by name alphabetically.
func SortGroups(groups []Group) {
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
}
