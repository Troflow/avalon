package prev_avalon

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
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

var (
	ErrNoSuchOption              = errors.New("There is no such option.")
	ErrNoSuchSpecialCharacter    = errors.New("That special character is not in this game.")
	ErrNotEnoughEvils            = errors.New("There are not enough evils in the game.")
	ErrNotEnoughPlayersForLake   = errors.New("There must be at least 7 players to enable the Lady of the Lake.")
	ErrNotEnoughPlayersForOberon = errors.New("There must be at least 10 players to enable Oberon.")
	ErrPlayerExists              = errors.New("Player has already joined.")
	ErrTooManyPlayers            = errors.New("Maximum number of players have joined already.")
	errNotEnoughPlayersToStart   = errors.New("There must be at least 5 players to close the lobby.")
)

// All possible GamePhases
const (
	PhaseWaitingForPlayers int = iota
	PhaseConfiguration
	PhaseQuestLake
	PhaseQuestPartyNominate
	PhaseQuestPartyVote
	PhaseQuestParty
	PhaseAssassination
)

// Configuration options: includes special characters, lake
const (
	OptionLadyOfTheLake int = iota
	OptionMordred
	OptionMorganaPercival
	OptionOberon
	NumOptions // Sneaky, sneaky
)

const (
	MinPlayersToStart int = 5
	MaxPlayers        int = 10
)

var (
	NumPlayersToNumEvils = map[int]int{
		5:  2,
		6:  2,
		7:  3,
		8:  3,
		9:  3,
		10: 4,
	}

	NumPlayersToQuestSizes = map[int][]int{
		5:  {2, 3, 2, 3, 3},
		6:  {2, 3, 4, 3, 4},
		7:  {2, 3, 3, 4, 4},
		8:  {3, 4, 4, 5, 5},
		9:  {3, 4, 4, 5, 5},
		10: {3, 4, 4, 5, 5},
	}

	GamePhaseToStatusString = map[int]string{
		PhaseWaitingForPlayers:  "Waiting for players to join.",
		PhaseConfiguration:      "Waiting for lobby leader to configure the game.",
		PhaseQuestLake:          "Waiting for the Lady of the Lake to act.",
		PhaseQuestPartyNominate: "Waiting for the quest leader to nominate a party.",
		PhaseQuestPartyVote:     "Waiting for players to vote on the nominated party.",
		PhaseQuestParty:         "The party is going on a quest.",
		PhaseAssassination:      "Assassin is trying to kill Merlin.",
	}

	StringToConfigOption = map[string]int{
		"mordred":         OptionMordred,
		"morganapercival": OptionMorganaPercival,
		"oberon":          OptionOberon,
		"lake":            OptionLadyOfTheLake,
	}
)

type AvalonGame struct {
	////////////////////////////
	// BASIC GAME INFORMATION //
	////////////////////////////

	// Tracks the progress of the game, from waiting for players to reveal
	// phase to quest voting to questing to assassination phase.
	GamePhase int

	// List of players.
	PlayerNicks []string

	// Assignments for good/evil. Full of nicks.
	GoodCharacters []string
	EvilCharacters []string

	// Map of role to nick. E.g. {"merlin": "Yulli", "assassin": "awe"}
	// Roles that do not appear in this map are not present in this game.
	SpecialCharacters map[string]string

	// Config options enabled or disabled?
	OptionEnabled []bool

	///////////////////////////////
	// CURRENT QUEST INFORMATION //
	///////////////////////////////

	// Current quest (1-5).
	CurrentQuest int

	// Current Lady of the Lake, a nick.
	CurrentLake string

	// Current quest party leader, a nick.
	CurrentLeader string

	// Current quest party leader's tentative party, a bunch of nicks.
	CurrentVoteParty []string

	// The vote track for the current quest (1-5).
	VoteTrack int

	////////////////////////////
	// PAST QUEST INFORMATION //
	////////////////////////////

	// Record of quests that have succeeded/failed.
	QuestSuccesses []bool
	// Record of quest parties.
	QuestParties [][]string
	// Record of leaders that nominated the parties for each quest.
	QuestFinalPartyLeader []string
}

func NewAvalonGame() *AvalonGame {
	ag := &AvalonGame{
		OptionEnabled:     make([]bool, NumOptions),
		SpecialCharacters: make(map[string]string),
	}
	return ag
}

func (ag *AvalonGame) NumPlayers() int {
	return len(ag.PlayerNicks)
}

func (ag *AvalonGame) NumGoods() int {
	return ag.NumPlayers() - ag.NumEvils()
}

func (ag *AvalonGame) NumEvils() int {
	return NumPlayersToNumEvils[ag.NumPlayers()]
}

func (ag *AvalonGame) QuestPartySize() int {
	return NumPlayersToQuestSizes[ag.NumPlayers()][ag.CurrentQuest-1]
}

func (ag *AvalonGame) GameStatusString() string {
	return GamePhaseToStatusString[ag.GamePhase]
}

func (ag *AvalonGame) AddPlayer(nick string) error {
	if ag.NumPlayers() == MaxPlayers {
		return ErrTooManyPlayers
	}

	if ag.PlayerExists(nick) {
		return ErrPlayerExists
	}

	ag.PlayerNicks = append(ag.PlayerNicks, nick)
	return nil
}

func (ag *AvalonGame) PlayerExists(nick string) bool {
	for _, player := range ag.PlayerNicks {
		if player == nick {
			return true
		}
	}

	return false
}

func (ag *AvalonGame) ListPlayers() string {
	return ListNicks(ag.PlayerNicks)
}

func (ag *AvalonGame) ListGoods() string {
	return ListNicks(ag.GoodCharacters)
}

func (ag *AvalonGame) ListEvils() string {
	return ListNicks(ag.EvilCharacters)
}

func (ag *AvalonGame) CanStartConfig() bool {
	return ag.NumPlayers() >= MinPlayersToStart
}

func (ag *AvalonGame) StartConfig() error {
	if !ag.CanStartConfig() {
		return errNotEnoughPlayersToStart
	}

	ag.GamePhase = PhaseConfiguration
	return nil
}

func (ag *AvalonGame) ConfigStatus() string {
	var enabledList []string

	if ag.OptionEnabled[OptionLadyOfTheLake] {
		enabledList = append(enabledList, "Lake")
	}

	if ag.OptionEnabled[OptionMordred] {
		enabledList = append(enabledList, "Mordred")
	}

	if ag.OptionEnabled[OptionMorganaPercival] {
		enabledList = append(enabledList, "MorganaPercival")
	}

	if ag.OptionEnabled[OptionOberon] {
		enabledList = append(enabledList, "Oberon")
	}

	if len(enabledList) == 0 {
		return "No config options are enabled."
	}

	return strings.Join(enabledList, ", ")
}

func (ag *AvalonGame) DisableOption(optionString string) error {
	option, err := ConvertStringToOption(optionString)
	if err != nil {
		return err
	}

	if !ag.OptionEnabled[option] {
		resp := fmt.Sprintf("%s is already disabled.", optionString)
		return errors.New(resp)
	}

	ag.OptionEnabled[option] = false
	return nil
}

func (ag *AvalonGame) EnableOption(optionString string) error {
	option, err := ConvertStringToOption(optionString)
	if err != nil {
		return err
	}

	if ag.OptionEnabled[option] {
		resp := fmt.Sprintf("%s is already enabled.", optionString)
		return errors.New(resp)
	}

	err = ag.CanEnableOption(option)
	if err != nil {
		return err
	}

	ag.OptionEnabled[option] = true
	return nil
}

// Evil specials must not exceed NumEvils() - 1
func (ag *AvalonGame) CanEnableOption(option int) error {
	switch option {
	case OptionLadyOfTheLake:
		if ag.NumPlayers() < 7 {
			return ErrNotEnoughPlayersForLake
		}
	case OptionMordred, OptionMorganaPercival, OptionOberon:
		if option == OptionOberon {
			if ag.NumPlayers() < 10 {
				return ErrNotEnoughPlayersForOberon
			}
		}

		if ag.NumEvilSpecials() >= ag.NumEvils()-1 {
			return ErrNotEnoughEvils
		}
	}
	return nil
}

func (ag *AvalonGame) NumEvilSpecials() int {
	var count int
	for _, option := range []int{OptionMordred, OptionMorganaPercival, OptionOberon} {
		if ag.OptionEnabled[option] {
			count++
		}
	}
	return count
}

func ConvertStringToOption(optionString string) (int, error) {
	if option, ok := StringToConfigOption[strings.ToLower(optionString)]; ok {
		return option, nil
	}
	return -1, ErrNoSuchOption
}

// Handles assigning players to good/evil, special roles, first quest leader,
// and Lady of the Lake.
func (ag *AvalonGame) AssignPlayers() {
	ag.assignGoodEvil()
	ag.assignSpecialCharacters()
	ag.assignFirstLeaderAndLake()
}

func (ag *AvalonGame) assignGoodEvil() {
	// Drop players into good/bad buckets.
	randomOrder := rand.Perm(len(ag.PlayerNicks))

	for i, n := range randomOrder {
		nick := ag.PlayerNicks[n]
		if i < ag.NumGoods() {
			ag.GoodCharacters = append(ag.GoodCharacters, nick)
		} else {
			ag.EvilCharacters = append(ag.EvilCharacters, nick)
		}
	}
}

func (ag *AvalonGame) assignSpecialCharacters() {
	// === Good ===
	randomOrder := rand.Perm(ag.NumGoods())
	// Merlin
	ag.SpecialCharacters["merlin"] = ag.GoodCharacters[randomOrder[0]]
	// Percival
	if ag.OptionEnabled[OptionMorganaPercival] {
		ag.SpecialCharacters["percival"] = ag.GoodCharacters[randomOrder[1]]
	}

	// === Evil ===
	randIndex := 1
	randomOrder = rand.Perm(ag.NumEvils())
	// Assassin
	ag.SpecialCharacters["assassin"] = ag.EvilCharacters[randomOrder[0]]
	// Mordred
	if ag.OptionEnabled[OptionMordred] {
		ag.SpecialCharacters["mordred"] = ag.EvilCharacters[randomOrder[randIndex]]
		randIndex++
	}
	// Morgana
	if ag.OptionEnabled[OptionMorganaPercival] {
		ag.SpecialCharacters["morgana"] = ag.EvilCharacters[randomOrder[randIndex]]
		randIndex++
	}
	// Oberon
	if ag.OptionEnabled[OptionOberon] {
		ag.SpecialCharacters["oberon"] = ag.EvilCharacters[randomOrder[randIndex]]
		randIndex++
	}
}

func (ag *AvalonGame) assignFirstLeaderAndLake() {
	randLeader := rand.Intn(ag.NumPlayers())
	ag.CurrentLeader = ag.PlayerNicks[randLeader]

	// Lady of the Lake, cannot be the first quest leader
	if ag.OptionEnabled[OptionLadyOfTheLake] {
		randLake := rand.Intn(ag.NumPlayers())
		for randLake != randLeader {
			randLake = rand.Intn(ag.NumPlayers())
		}

		ag.CurrentLake = ag.PlayerNicks[randLake]
	}
}

func (ag *AvalonGame) EvilCharactersWithoutMordred() []string {
	return ag.evilCharactersExcluding("mordred")
}

func (ag *AvalonGame) EvilCharactersWithoutOberon() []string {
	return ag.evilCharactersExcluding("oberon")
}

func (ag *AvalonGame) evilCharactersExcluding(character string) []string {
	evils := make([]string, len(ag.EvilCharacters))
	copy(evils, ag.EvilCharacters)

	option := StringToConfigOption[strings.ToLower(character)]
	if ag.OptionEnabled[option] {
		exclude := ag.SpecialCharacters[character]
		i := indexOf(evils, exclude)
		return deleteAt(evils, i)
	}

	return evils
}
func (ag *AvalonGame) SpecialCharacterNick(character string) (string, error) {
	nick, ok := ag.SpecialCharacters[character]
	if !ok {
		return "", ErrNoSuchSpecialCharacter
	}

	return nick, nil
}

func (ag *AvalonGame) FlavorText(character string) string {
	switch strings.ToLower(character) {
	case "mordred":
		return "You remain unknown to Merlin."
	case "morgana":
		return "You appear as Merlin to Percival."
	case "percival":
		return "You see Merlin's identity. Morgana attempts to trick you."
	case "oberon":
		return "You are unknown to the other evils, nor do you know them."
	default:
		return ""
	}
}

func indexOf(ary []string, item string) int {
	for i, e := range ary {
		if e == item {
			return i
		}
	}
	return -1
}

func deleteAt(ary []string, i int) []string {
	if i == len(ary) {
		return ary[:i]
	}

	return append(ary[:i], ary[i+1:]...)
}

func ListNicks(list []string) string {
	return strings.Join(list, ", ")
}
