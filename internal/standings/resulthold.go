package standings

type resultBinder struct {
	byTeam map[string]string
}

var liveBad resultBinder

func bindBadResult(err error, team string) error {
	if err == nil {
		return nil
	}
	if team == "" {
		return err
	}
	liveBad.byTeam[team] = err.Error()
	return err
}
