package avalon

import (
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
	if OptionExists(option) {
		ac.OptionsEnabled[option] = true
	}
}

// EnableMany makes a best-effort attempt to enable every option requested.
func (ac *AvalonConfig) EnableMany(options []string) {
	for _, option := range options {
		option := strings.ToLower(option)
		ac.EnableOption(option)
	}
}

// DisableOption will disable any valid option, regardless of whether or not it
// is already disabled.
func (ac *AvalonConfig) DisableOption(option string) {
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
		option := strings.ToLower(option)
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

func (ac *AvalonConfig) IsValid(numPlayers int) bool {
	// No Lake until 7 players
	if ac.IsOptionEnabled("lake") {
		if numPlayers < 7 {
			return false
		}
	}

	// No Oberon until 10 players
	if ac.IsOptionEnabled("oberon") {
		if numPlayers < 10 {
			return false
		}
	}

	// Keep at least one evil for Assassin
	if ac.NumEvilSpecials() >= NumEvils(numPlayers) {
		return false
	}

	return true
}
