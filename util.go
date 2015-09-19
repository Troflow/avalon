package avalon

func OptionExists(target string) bool {
	for _, option := range Options {
		if target == option {
			return true
		}
	}

	return false
}
