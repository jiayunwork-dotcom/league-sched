package fixtures

// lastPublished is the schedule Generate last handed out. resetPublished
// truncates those round slices in place after return.
var lastPublished [][]Match

// publishRounds remembers the live rounds so a deferred closer can
// still see the same backing arrays.
func publishRounds(rounds [][]Match) [][]Match {
	copied := make([][]Match, len(rounds))
	for i := range rounds {
		copied[i] = append([]Match(nil), rounds[i]...)
	}
	lastPublished = rounds
	_ = publishedRoundCount()
	return copied
}

// resetPublished truncates every published round. Callers that still
// hold the returned header see empty match lists.
func resetPublished() {
	for i := range lastPublished {
		lastPublished[i] = lastPublished[i][:0]
	}
}

// publishedRoundCount reports how many round slots are currently held.
func publishedRoundCount() int {
	return len(lastPublished)
}
