package avalon

import (
	"errors"
	"fmt"
	"strings"
)

type AvalonConfig struct {
	OptionsEnabled map[string]bool
}

func NewAvalonConfig() *AvalonConfig {
	return &AvalonConfig{
		OptionsEnabled: make(map[string]bool),
	}
}

// IsOptionEnabled returns whether or not the specified option is enabled.
func (ac *AvalonConfig) IsOptionEnabled(option string) bool {
	return ac.OptionsEnabled[option]
}

// ListEnabledOptions returns a human-readable list of enabled ac options.
// Example: "Lake, Mordred, Oberon"
func (ac *AvalonConfig) ListEnabledOptions() string {
	if len(ac.OptionsEnabled) == 0 {
		return "none"
	}

	enabled := ac.allEnabledOptions()
	return strings.Join(enabled, ", ")
}

// EnableOption will enable any valid option, regardless of whether or not it
// is already enabled.
func (ac *AvalonConfig) EnableOption(option string) {
	option = strings.ToLower(option)
	if OptionExists(option) {
		ac.OptionsEnabled[option] = true
	}
}

// EnableMany makes a best-effort attempt to enable every option requested.
func (ac *AvalonConfig) EnableMany(options []string) {
	for _, option := range options {
		ac.EnableOption(option)
	}
}

// DisableOption will disable any valid option, regardless of whether or not it
// is already disabled.
func (ac *AvalonConfig) DisableOption(option string) {
	option = strings.ToLower(option)

	if !ac.IsOptionEnabled(option) {
		return
	}

	if _, ok := ac.OptionsEnabled[option]; ok {
		delete(ac.OptionsEnabled, option)
	}
}

// DisableMany makes a best-effort attempt to disable every option requested.
func (ac *AvalonConfig) DisableMany(options []string) {
	for _, option := range options {
		ac.DisableOption(option)
	}
}

func (ac *AvalonConfig) allEnabledOptions() []string {
	var options []string
	for k := range ac.OptionsEnabled {
		options = append(options, k)
	}

	return options
}

// NumEvilSpecials returns the number of special evil characters enabled.
func (ac *AvalonConfig) NumEvilSpecials() int {
	var count int
	for _, option := range SpecialEvilOptions {
		if ac.IsOptionEnabled(option) {
			count++
		}
	}

	return count
}

// IsValid verifies that the config is valid for the number of players
// specified.
func (ac *AvalonConfig) IsValid(numPlayers int) error {
	var errorStrings []string

	// No Lake until 7 players
	if ac.IsOptionEnabled("lake") {
		if numPlayers < 7 {
			errorStrings = append(errorStrings, "lake requires 7 players")
		}
	}

	// No Oberon until 10 players
	if ac.IsOptionEnabled("oberon") {
		if numPlayers < 10 {
			errorStrings = append(errorStrings, "oberon requires 10 players")
		}
	}

	// Keep at least one evil for Assassin
	numEvils := NumEvils(numPlayers)
	numSpecials := ac.NumEvilSpecials()
	if numSpecials >= numEvils {
		n := numSpecials - numEvils + 1
		errorStrings = append(errorStrings, fmt.Sprintf("you have %d too many evils", n))
	}

	if len(errorStrings) > 0 {
		return errors.New(strings.Join(errorStrings, "; "))
	}

	return nil
}
