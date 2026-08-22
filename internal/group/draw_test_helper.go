package group

// CountTeams returns the total number of teams across all groups.
func CountTeams(groups []Group) int {
	count := 0
	for _, g := range groups {
		count += len(g.Teams)
	}
	return count
}

// GroupSizes returns the size of each group.
func GroupSizes(groups []Group) []int {
	sizes := make([]int, len(groups))
	for i, g := range groups {
		sizes[i] = len(g.Teams)
	}
	return sizes
}

// AllTeams returns all teams across all groups.
func AllTeams(groups []Group) []string {
	var out []string
	for _, g := range groups {
		out = append(out, g.Teams...)
	}
	return out
}

// IsTeamInGroup checks if a team is in a specific group.
func IsTeamInGroup(g Group, team string) bool {
	for _, t := range g.Teams {
		if t == team {
			return true
		}
	}
	return false
}

// MaxGroupSize returns the size of the largest group.
func MaxGroupSize(groups []Group) int {
	max := 0
	for _, g := range groups {
		if len(g.Teams) > max {
			max = len(g.Teams)
		}
	}
	return max
}

// MinGroupSize returns the size of the smallest group.
func MinGroupSize(groups []Group) int {
	if len(groups) == 0 {
		return 0
	}
	min := len(groups[0].Teams)
	for _, g := range groups[1:] {
		if len(g.Teams) < min {
			min = len(g.Teams)
		}
	}
	return min
}
