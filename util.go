package avalon

func OptionExists(target string) bool {
	for _, option := range AvailableOptions {
		if target == option {
			return true
		}
	}

	return false
}
