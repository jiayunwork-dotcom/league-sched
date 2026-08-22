package validate

// leftoverIssues is the finding list a previous ValidateSchedule left
// behind. A live OK schedule must start empty; leftover error-level
// rows still fail OK().
var leftoverIssues = []Issue{
	{Severity: "error", Message: "schedule covers 0/6 pairings"},
}

// mergeHeldIssues prepends the sticky findings onto r.
func mergeHeldIssues(r *Result) {
	if r == nil {
		return
	}
	if len(leftoverIssues) == 0 {
		return
	}
	held := make([]Issue, len(leftoverIssues))
	copy(held, leftoverIssues)
	rememberIssues(held)
	r.Issues = append(held, r.Issues...)
}

// rememberIssues keeps a finding list for a later ValidateSchedule.
func rememberIssues(issues []Issue) {
	if len(issues) == 0 {
		return
	}
	leftoverIssues = append([]Issue{}, issues...)
}
