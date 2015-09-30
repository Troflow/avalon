// Package avalon provides an abstraction for creating and running a game of
// Avalon.
package avalon

import (
	"errors"
)

var (
	ErrPlayerExists   = errors.New("avalon: player is already in this game")
	ErrTooManyPlayers = errors.New("avalon: there are already 10 players")
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
	*AvalonConfig

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
	av := &Avalon{
		AvalonConfig: NewAvalonConfig(),
		Specials:     make(map[string]string),
	}

	return av
}

// NumPlayers returns the number of players in the game.
func (av *Avalon) NumPlayers() int {
	return len(av.Players)
}

// NumEvils returns the number of evil characters based on the total number of
// players.
func (av *Avalon) NumEvils() int {
	return NumEvils(av.NumPlayers())
}

// NumGoods returns the number of good characters based on the total number of
// players.
func (av *Avalon) NumGoods() int {
	numEvils := av.NumEvils()
	if numEvils == 0 {
		return 0
	}

	return av.NumPlayers() - numEvils
}

// PlayerExists checks the game's current players to see if a given nick is
// already registered.
func (av *Avalon) PlayerExists(nick string) bool {
	for _, player := range av.Players {
		if player == nick {
			return true
		}
	}

	return false
}

// AddPlayer attempts to add a new player to the list of players. Errors if the
// player already exists or they are too many players.
func (av *Avalon) AddPlayer(nick string) error {
	if av.NumPlayers() >= 10 {
		return ErrTooManyPlayers
	}

	if av.PlayerExists(nick) {
		return ErrPlayerExists
	}

	av.Players = append(av.Players, nick)
	return nil
}

// IsValid overrides the AvalonConfig IsValid and does not require a numPlayers
// to be passed in.
func (av *Avalon) IsValid() error {
	return av.AvalonConfig.IsValid(av.NumPlayers())
}
