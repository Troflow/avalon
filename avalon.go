// Package avalon provides an abstraction for creating and running a game of
// Avalon.
package avalon

import (
	"errors"
)

// | Players | Evils | Q1 | Q2 | Q3 | Q4 | Q5 |
// |---------|-------|----|----|----|----|----|
// |       5 |     2 |  2 |  3 |  2 |  3 |  3 |
// |       6 |     2 |  2 |  3 |  4 |  3 |  4 |
// |       7 |     3 |  2 |  3 |  3 |  4 |  4 |
// |       8 |     3 |  3 |  4 |  4 |  5 |  5 |
// |       9 |     3 |  3 |  4 |  4 |  5 |  5 |
// |      10 |     4 |  3 |  4 |  4 |  5 |  5 |

// With 7+ players, Quest 4 requires two fails to fail.

// Avalon represents underlying game state and facilitates changes to the game
// as it progresses and win conditions.
type Avalon struct {
	// Players and their roles
	Players  []string
	Goods    []string
	Evils    []string
	Specials map[string]string

	// Game state information that changes throughout the game's lifecycle
	CurrentQuest         int
	CurrentLake          string
	CurrentLeader        string
	CurrentProposedParty []string
	VoteTrack            int

	// Information about past quests
	QuestSuccesses   []bool
	PastQuestParties []string
}

func NewAvalon() *Avalon {
	return &Avalon{
		Specials: make(map[string]string),
	}
}

// NumPlayers returns the number of players in the game.
func (av *Avalon) NumPlayers() int {
	return len(av.Players)
}

// NumEvils returns the number of evil characters based on the total number of
// players.
func (av *Avalon) NumEvils() (int, error) {
	numEvils, ok := numPlayersToNumEvils[av.NumPlayers()]
	if !ok {
		return 0, errors.New("avalon: not enough players to determine evils")
	}

	return numEvils, nil
}

// NumGoods returns the number of good characters based on the total number of
// players.
func (av *Avalon) NumGoods() (int, error) {
	numEvils, err := av.NumEvils()
	if err != nil {
		return 0, err
	}

	return av.NumPlayers() - numEvils, nil
}
