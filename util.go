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

func FlavorTextForSpecial(special string) string {
	text, ok := specialCharacterToFlavorText[special]
	if !ok {
		return ""
	}

	return text
}

func remove(list []string, target string) {
	for i, e := range list {
		if e == target {
			deleteAt(list, i)
		}
	}
}

func deleteAt(list []string, i int) {
	if i == len(list) {
		list = list[:i]
	}

	list = append(list[:i], list[i+1:]...)
}
