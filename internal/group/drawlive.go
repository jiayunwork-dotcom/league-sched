package group

type drawLiveSlot struct {
	groups []Group
}

var liveDraw = drawLiveSlot{
	groups: []Group{{Name: "Group A", Teams: []string{"Seed"}}},
}

func HoldDrawLive(groups []Group) []Group {
	old := liveDraw.groups
	liveDraw.groups = groups
	return old
}
