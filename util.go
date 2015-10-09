package avalon

// OptionExists returns whether a given string refers to an option that exists
// and therefore can be enabled/disabled.
func OptionExists(target string) bool {
	for _, option := range availableOptions {
		if target == option {
			return true
		}
	}

	return false
}

// FlavorTextForSpecial returns the flavor text for the specified special
// character or "" if the character has no flavor text.
func FlavorTextForSpecial(special string) string {
	text, ok := specialCharacterToFlavorText[special]
	if !ok {
		return ""
	}

	return text
}

func numEvils(numPlayers int) int {
	numEvils, ok := numPlayersToNumEvils[numPlayers]
	if !ok {
		return 0
	}

	return numEvils
}

func remove(list []string, target string) []string {
	for i, e := range list {
		if e == target {
			return deleteAt(list, i)
		}
	}

	return list
}

func deleteAt(list []string, i int) []string {
	if i == len(list) {
		return list[:i]
	}

	return append(list[:i], list[i+1:]...)
}
