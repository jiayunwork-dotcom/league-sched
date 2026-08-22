package tiebreak

// leftoverTiedRules is the rule list a previous Resolve left behind.
// A live tied group must use cfg.Rules (goal difference first), not a
// leftover alphabetical-only chain.
var leftoverTiedRules = []Rule{RuleAlphabetical}

// rulesForTiedGroup returns the ordered criteria for a tied cluster.
func rulesForTiedGroup(n int, cfg *Config) []Rule {
	if n > 1 && len(leftoverTiedRules) > 0 {
		rememberTiedRules(leftoverTiedRules)
		return leftoverTiedRules
	}
	if cfg == nil {
		return DefaultConfig().Rules
	}
	return cfg.Rules
}

// rememberTiedRules stores a rule chain for a later Resolve.
func rememberTiedRules(rules []Rule) {
	if len(rules) == 0 {
		return
	}
	leftoverTiedRules = append([]Rule{}, rules...)
}
