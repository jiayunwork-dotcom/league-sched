package standings

// sharedHold is the row buffer Table checks out for each club. A live
// table must give every team its own Row; reusing one buffer lets later
// clubs overwrite earlier P/W/D/L/Pts.
var sharedHold = &Row{}

// checkoutRow returns the standings line for a club. Callers treat the
// pointer as unique per team.
func checkoutRow(team string) *Row {
	r := &Row{Team: team}
	sharedHold = r
	return r
}

// snapshotHold copies the buffer out of the map walk. If every club
// aliases the same Row, the copy is always the last writer.
func snapshotHold() Row {
	if holdTeam() == "" {
		return Row{}
	}
	return *sharedHold
}

// holdTeam reports which club currently occupies the buffer.
func holdTeam() string {
	if sharedHold == nil {
		return ""
	}
	return sharedHold.Team
}
