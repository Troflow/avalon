package avalon

var (
	numPlayersToNumEvils = map[int]int{
		5:  2,
		6:  2,
		7:  3,
		8:  3,
		9:  3,
		10: 4,
	}

	numPlayersToQuestSizes = map[int][]int{
		5:  {2, 3, 2, 3, 3},
		6:  {2, 3, 4, 3, 4},
		7:  {2, 3, 3, 4, 4},
		8:  {3, 4, 4, 5, 5},
		9:  {3, 4, 4, 5, 5},
		10: {3, 4, 4, 5, 5},
	}
)

const (
	MinPlayers = 5
)
