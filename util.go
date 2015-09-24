package avalon

func OptionExists(target string) bool {
	for _, option := range AvailableOptions {
		if target == option {
			return true
		}
	}

	return false
}

func CanEnable(option string, av *Avalon) bool {
	switch option {
	case "lake":
		// No Lake until at least 7 players
		if av.NumPlayers() < 7 {
			return false
		}
	case "mordred", "morganapercival", "oberon":
		// No Oberon until at least 10 players
		if option == "oberon" && av.NumPlayers() < 10 {
			return false
		}

		// Keep at least one spare evil for Assassin
		if av.NumEvilSpecials() >= av.NumEvils()-1 {
			return false
		}
	}

	return true
}
