package avalon

func OptionExists(target string) bool {
	for _, option := range AvailableOptions {
		if target == option {
			return true
		}
	}

	return false
}

func NumEvils(numPlayers int) int {
	numEvils, ok := numPlayersToNumEvils[numPlayers]
	if !ok {
		return 0
	}

	return numEvils
}
