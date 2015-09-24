// Package avalon provides an abstraction for creating and running a game of
// Avalon.
package avalon

import (
	"errors"
	"strings"
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
	// Meta information
	EnabledOptions map[string]bool

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
		EnabledOptions: make(map[string]bool),
		Specials:       make(map[string]string),
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

// ListEnabledOptions returns a human-readable list of all the config options
// enabled for this game or a special message if there are none.
// Example: "Lake, Mordred, Oberon"
func (av *Avalon) ListEnabledOptions() string {
	var enabledList []string

	for option, enabled := range av.EnabledOptions {
		if enabled {
			enabledList = append(enabledList, option)
		}
	}

	if len(enabledList) == 0 {
		return "none"
	}

	return strings.Join(enabledList, ", ")
}

// EnableOptions makes a best-effort attempt to enable every option requested
// and silently fails on options it cannot enable.
func (av *Avalon) EnableOptions(options []string) {
	for _, option := range options {
		option := strings.ToLower(option)
		if OptionExists(option) {
			av.EnabledOptions[option] = true
		}
	}
}
