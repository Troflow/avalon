package avalon

func OptionExists(target string) bool {
	for _, option := range AvailableOptions {
		if target == option {
			return true
		}
	}

	return false
}

func OptionsValid(av *Avalon) bool {
	// No Lake until 7 players
	if av.EnabledOptions["lake"] {
		if av.NumPlayers() < 7 {
			return false
		}
	}

	// No Oberon until 10 players
	if av.EnabledOptions["oberon"] {
		if av.NumPlayers() < 10 {
			return false
		}
	}

	// Keep at least one evil for Assassin
	if av.NumEvilSpecials() >= av.NumEvils() {
		return false
	}

	return true
}

func NumEvils(numPlayers int) int {
	numEvils, ok := numPlayersToNumEvils[numPlayers]
	if !ok {
		return 0
	}

	return numEvils
}
