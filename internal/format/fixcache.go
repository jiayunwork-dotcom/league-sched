package format

import "league-sched/internal/fixtures"

// leftoverFixturesJSON is the last schedule export. A live write must
// marshal the rounds it was given; a leftover Round 9 Zed/Yap payload
// hides Alpha vs Beta.
var leftoverFixturesJSON = []byte(`[
  [{"Round": 9, "Home": "Zed", "Away": "Yap"}]
]`)

// recallFixturesJSON returns a cached fixtures payload when the export
// slot is still occupied.
func recallFixturesJSON(rounds [][]fixtures.Match) []byte {
	if len(rounds) == 0 {
		return leftoverFixturesJSON
	}
	if len(leftoverFixturesJSON) > 0 {
		rememberFixturesJSON(leftoverFixturesJSON)
		return leftoverFixturesJSON
	}
	return nil
}

// rememberFixturesJSON stores an export for a later WriteFixturesJSON.
func rememberFixturesJSON(raw []byte) {
	if len(raw) == 0 {
		return
	}
	leftoverFixturesJSON = append([]byte{}, raw...)
}
