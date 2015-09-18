package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	cmodels "github.com/justinkim/yullibot/commands/models"
	"github.com/justinkim/yullibot/irc"
	"github.com/justinkim/yullibot/models"
)

var (
	alreadyJoinedString            = "You have already joined this game."
	availableOptionsString         = "Available subsubcommands: mordred, morganapercival, oberon, lake"
	configurationHelpString        = "Use `!avalon enable` and `!avalon disable` to configure the game. Use !avalon start again when you're done."
	gameImplodedString             = "The game has been imploded by the lobby leader."
	gameInProgressString           = "There is already a game in progress."
	helpString                     = "!avalon subcommands are: disable, enable, help, implode, info, init, join, start, status"
	infoHelpString                 = "Info subsubcommands are: game, players."
	lobbyClosedString              = "The lobby leader has already closed the lobby."
	newGameString                  = "New game created. Players can join using !avalon join."
	noActiveGameString             = "No game is currently active. Start one with !avalon init."
	noDMCommandString              = "You cannot use this command through DM."
	onlyLobbyLeaderCanDoThatString = "Only the lobby leader can do that."
	tooManyPlayersString           = "There are already 10 players in this game."

	closingLobbyFormatString         = "Closing the lobby with %d players."
	configOptionsEnabledFormatString = "Config options enabled: %s."
	correctChannelFormatString       = "You must do this in %s."
	currentLeaderFormatString        = "The current leader for Quest %d Party %d is %s"
	evilTeamFormatString             = "You see those allied with Mordred: %s."
	joinProgressFormatString         = "%d players have joined so far."
	infoGameFormatString             = "%s is running a game in %s."
	infoPlayersFormatString          = "This game has %d player%s: %s."
	lakeFormatString                 = "The Lady of the Lake is %s."
	percivalRevealFormatString       = "Merlin and Morgana are %s and %s. You are not sure which is which."
	readyToConfigFormatString        = "%s: There are enough players to begin the game. Use !avalon start to start or wait for more players."
	startingGameFormatString         = "%s: Starting game with: %s."
)

type AvalonCommand struct {
	triggerStrings []string

	game *cmodels.AvalonGame

	// The channel the game will take place in.
	channel string

	// The player who initiated the game and controls when to begin.
	// Currently serves no purpose other than starting the game with the
	// correct number of players and deciding which special characters to
	// play with.
	lobbyLeader string
}

func NewAvalonCommand() *AvalonCommand {
	return &AvalonCommand{
		triggerStrings: []string{
			"avalon",
			"a",
		},
	}
}

// Jfc.
func (avalon *AvalonCommand) TriggerStrings() []string {
	return avalon.triggerStrings
}

func (avalon *AvalonCommand) Process(bot *models.Bot, msg *irc.Message) {
	if !models.CommandTriggered(avalon, msg) {
		return
	}

	subcommand := msg.Subcommand()
	switch subcommand {
	case "disable":
		avalon.DoDisable(bot, msg)
	case "enable":
		avalon.DoEnable(bot, msg)
	case "help", "":
		avalon.DoHelp(bot, msg)
	case "implode":
		avalon.DoImplode(bot, msg)
	case "info":
		avalon.DoInfo(bot, msg)
	case "init":
		avalon.DoInit(bot, msg)
	case "join":
		avalon.DoJoin(bot, msg)
	case "start":
		avalon.DoStart(bot, msg)
	case "status":
		avalon.DoStatus(bot, msg)
	}
}

func (avalon *AvalonCommand) DoDisable(bot *models.Bot, msg *irc.Message) {
	err := avalon.IssuedByLobbyLeaderInGameChannel(msg)
	if err != nil {
		bot.Connection.RespondTo(msg, err.Error())
		return
	}

	option := msg.Subsubcommand()
	if option == "" {
		bot.Connection.RespondTo(msg, availableOptionsString)
		return
	}

	err = avalon.game.DisableOption(option)
	if err != nil {
		bot.Connection.RespondTo(msg, err.Error())
	}

	resp := fmt.Sprintf(configOptionsEnabledFormatString, avalon.game.ConfigStatus())
	bot.Connection.RespondTo(msg, resp)
}

func (avalon *AvalonCommand) DoEnable(bot *models.Bot, msg *irc.Message) {
	err := avalon.IssuedByLobbyLeaderInGameChannel(msg)
	if err != nil {
		bot.Connection.RespondTo(msg, err.Error())
		return
	}

	option := msg.Subsubcommand()
	if option == "" {
		bot.Connection.RespondTo(msg, availableOptionsString)
		return
	}

	err = avalon.game.EnableOption(option)
	if err != nil {
		bot.Connection.RespondTo(msg, err.Error())
	}

	resp := fmt.Sprintf(configOptionsEnabledFormatString, avalon.game.ConfigStatus())
	bot.Connection.RespondTo(msg, resp)
}

func (avalon *AvalonCommand) DoHelp(bot *models.Bot, msg *irc.Message) {
	bot.Connection.RespondTo(msg, helpString)
}

func (avalon *AvalonCommand) DoImplode(bot *models.Bot, msg *irc.Message) {
	err := avalon.IssuedByLobbyLeaderInGameChannel(msg)
	if err != nil {
		bot.Connection.RespondTo(msg, err.Error())
		return
	}

	resp := fmt.Sprintf("%s: %s", avalon.game.ListPlayers(), gameImplodedString)
	bot.Connection.Say(msg.Channel(), resp)

	avalon.channel = ""
	avalon.game = nil
}

func (avalon *AvalonCommand) DoInfo(bot *models.Bot, msg *irc.Message) {
	if avalon.game == nil {
		bot.Connection.RespondTo(msg, noActiveGameString)
		return
	}

	subsub := msg.Subsubcommand()
	switch subsub {
	case "":
		bot.Connection.RespondTo(msg, infoHelpString)
	case "game":
		resp := fmt.Sprintf(infoGameFormatString, avalon.lobbyLeader, avalon.channel)
		bot.Connection.RespondTo(msg, resp)
	case "players":
		plural := ""
		if avalon.game.NumPlayers() > 1 {
			plural = "s"
		}

		resp := fmt.Sprintf(infoPlayersFormatString, avalon.game.NumPlayers(), plural, avalon.game.ListPlayers())
		bot.Connection.RespondTo(msg, resp)
	}
}

func (avalon *AvalonCommand) DoInit(bot *models.Bot, msg *irc.Message) {
	if msg.IsDirectMessage() {
		bot.Connection.RespondTo(msg, noDMCommandString)
		return
	}

	if avalon.game != nil {
		bot.Connection.RespondTo(msg, gameInProgressString)
		return
	}

	avalon.channel = msg.Channel()
	avalon.lobbyLeader = msg.SenderNick()

	avalon.game = cmodels.NewAvalonGame()
	avalon.game.AddPlayer(avalon.lobbyLeader)

	bot.Connection.RespondTo(msg, newGameString)
}

func (avalon *AvalonCommand) DoJoin(bot *models.Bot, msg *irc.Message) {
	if msg.IsDirectMessage() {
		bot.Connection.RespondTo(msg, noDMCommandString)
		return
	}

	if avalon.game == nil {
		bot.Connection.RespondTo(msg, noActiveGameString)
		return
	}

	if avalon.game.GamePhase != cmodels.PhaseWaitingForPlayers {
		bot.Connection.RespondTo(msg, lobbyClosedString)
		return
	}

	if avalon.channel == msg.Channel() {
		err := avalon.game.AddPlayer(msg.SenderNick())
		if err != nil {
			if err == cmodels.ErrPlayerExists {
				bot.Connection.RespondTo(msg, alreadyJoinedString)
			} else if err == cmodels.ErrTooManyPlayers {
				bot.Connection.RespondTo(msg, tooManyPlayersString)
			}
			return
		}
		bot.Connection.RespondTo(msg, playerAddedString)

		joinProgress := fmt.Sprintf(joinProgressFormatString, avalon.game.NumPlayers())
		bot.Connection.Say(msg.Channel(), joinProgress)

		if avalon.game.CanStartConfig() {
			resp := fmt.Sprintf(readyToConfigFormatString, avalon.lobbyLeader)
			bot.Connection.Say(msg.Channel(), resp)
		}
	}
}

func (avalon *AvalonCommand) DoStart(bot *models.Bot, msg *irc.Message) {
	err := avalon.IssuedByLobbyLeaderInGameChannel(msg)
	if err != nil {
		bot.Connection.RespondTo(msg, err.Error())
		return
	}

	if avalon.game.GamePhase == cmodels.PhaseWaitingForPlayers {
		err := avalon.game.StartConfig()
		if err != nil {
			bot.Connection.RespondTo(msg, err.Error())
			return
		}

		resp := fmt.Sprintf(closingLobbyFormatString, avalon.game.NumPlayers())
		bot.Connection.RespondTo(msg, resp)

		bot.Connection.RespondTo(msg, configurationHelpString)
	} else if avalon.game.GamePhase == cmodels.PhaseConfiguration {
		resp := fmt.Sprintf(startingGameFormatString, avalon.game.ListPlayers(), avalon.game.ConfigStatus())
		bot.Connection.Say(msg.Channel(), resp)

		avalon.game.AssignPlayers()
		avalon.Reveal(bot, msg)
		avalon.StartFirstQuest()
	}
}

func (avalon *AvalonCommand) IssuedByLobbyLeaderInGameChannel(msg *irc.Message) error {
	if msg.IsDirectMessage() {
		return errors.New(noDMCommandString)
	}

	if avalon.game == nil {
		return errors.New(noActiveGameString)
	}

	if avalon.lobbyLeader != msg.SenderNick() {
		return errors.New(onlyLobbyLeaderCanDoThatString)
	}

	if avalon.channel != msg.Channel() {
		resp := fmt.Sprintf(correctChannelFormatString, avalon.channel)
		return errors.New(resp)
	}

	return nil
}

func (avalon *AvalonCommand) DoStatus(bot *models.Bot, msg *irc.Message) {
	if avalon.game == nil {
		bot.Connection.RespondTo(msg, noActiveGameString)
		return
	}

	bot.Connection.RespondTo(msg, avalon.game.GameStatusString())
}

func (avalon *AvalonCommand) Reveal(bot *models.Bot, msg *irc.Message) {
	avalon.revealToSelves(bot)
	avalon.revealToEvils(bot)
	avalon.revealToMerlin(bot)
	avalon.revealToPercival(bot)
	avalon.revealLake(bot)
}

func (avalon *AvalonCommand) revealToSelves(bot *models.Bot) {
	for _, nick := range avalon.game.GoodCharacters {
		bot.Connection.Say(nick, "You are on the side of good, allied with Merlin.")
	}

	for _, nick := range avalon.game.EvilCharacters {
		bot.Connection.Say(nick, "You are on the side of evil, allied with Mordred.")
	}

	for special, nick := range avalon.game.SpecialCharacters {
		resp := fmt.Sprintf("You are %s. %s", capitalize(special), avalon.game.FlavorText(special))
		bot.Connection.Say(nick, resp)
	}
}

func capitalize(str string) string {
	if len(str) == 0 {
		return ""
	}

	if len(str) == 1 {
		return strings.ToUpper(str)
	}

	first := str[:1]
	first = strings.ToUpper(first)
	return first + str[1:]
}

// Evils see each other, except for Oberon
func (avalon *AvalonCommand) revealToEvils(bot *models.Bot) {
	visibleEvils := avalon.game.EvilCharactersWithoutOberon()
	resp := fmt.Sprintf(evilTeamFormatString, cmodels.ListNicks(visibleEvils))

	for _, evil := range visibleEvils {
		bot.Connection.Say(evil, resp)
	}
}

// Reveal evils to Merlin, but not Mordred
func (avalon *AvalonCommand) revealToMerlin(bot *models.Bot) {
	merlin := avalon.game.SpecialCharacters["merlin"]
	merlinSees := avalon.game.EvilCharactersWithoutMordred()

	resp := fmt.Sprintf("You see the Minions of Mordred: %s.", cmodels.ListNicks(merlinSees))
	bot.Connection.Say(merlin, resp)
}

// Reveal Merlin and Morgana to Percival
func (avalon *AvalonCommand) revealToPercival(bot *models.Bot) {
	if avalon.game.OptionEnabled[cmodels.OptionMorganaPercival] {
		merlinNick, err := avalon.game.SpecialCharacterNick("merlin")
		if err != nil {
			return
		}

		morganaNick, err := avalon.game.SpecialCharacterNick("morgana")
		if err != nil {
			return
		}

		percivalNick, err := avalon.game.SpecialCharacterNick("percival")
		if err != nil {
			return
		}

		// Randomly swap Merlin and Morgana so Percival can't tell
		// which is which.
		if rand.Int()%2 == 0 {
			merlinNick, morganaNick = morganaNick, merlinNick
		}

		resp := fmt.Sprintf(percivalRevealFormatString, merlinNick, morganaNick)
		bot.Connection.Say(percivalNick, resp)
	}
}

// Tell everyone who is the first Lady of the Lake
func (avalon *AvalonCommand) revealLake(bot *models.Bot) {
	if avalon.game.OptionEnabled[cmodels.OptionLadyOfTheLake] {
		resp := fmt.Sprintf(lakeFormatString, avalon.game.CurrentLake)
		bot.Connection.Say(avalon.channel, resp)
	}
}
