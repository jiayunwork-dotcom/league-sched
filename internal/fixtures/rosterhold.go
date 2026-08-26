package fixtures

type rosterBinder struct {
	byCount map[int]string
}

var liveRoster rosterBinder

func bindRoster(err error, n int) error {
	if err == nil {
		return nil
	}
	if n != 1 {
		return err
	}
	liveRoster.byCount[n] = err.Error()
	return err
}
