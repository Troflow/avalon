package avalon

import (
	"reflect"
	"strings"
	"testing"
)

func TestIsOptionEnabled(t *testing.T) {
	tests := []struct {
		enabled map[string]bool
		option  string
		want    bool
	}{
		{
			map[string]bool{},
			"oberon",
			false,
		},
		{
			map[string]bool{"oberon": true},
			"oberon",
			true,
		},
		{
			map[string]bool{"oberon": true},
			"mordred",
			false,
		},
		{
			map[string]bool{"oberon": true, "mordred": true},
			"mordred",
			true,
		},
		{
			map[string]bool{"oberon": true, "mordred": true},
			"lake",
			false,
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.enabled
		res := config.IsOptionEnabled(test.option)

		if res != test.want {
			t.Errorf("expected %t, but got %t", test.want, res)
		}
	}
}

func TestListEnabledOptions(t *testing.T) {
	tests := []struct {
		options map[string]bool
		wants   []string
	}{
		{
			map[string]bool{"oberon": true},
			[]string{"oberon"},
		},
		{
			map[string]bool{},
			[]string{"none"},
		},
		{
			map[string]bool{"mordred": true, "lake": true, "oberon": true},
			[]string{"mordred", "lake", "oberon"},
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.options
		res := config.ListEnabledOptions()

		for _, want := range test.wants {
			if !strings.Contains(res, want) {
				t.Errorf("expected %s, but was absent from %s", want, res)
			}
		}
	}
}

func TestEnableOption(t *testing.T) {
	tests := []struct {
		enabled map[string]bool
		option  string
		want    map[string]bool
	}{
		{
			map[string]bool{},
			"mordred",
			map[string]bool{"mordred": true},
		},
		{
			map[string]bool{"oberon": true},
			"oberon",
			map[string]bool{"oberon": true},
		},
		{
			map[string]bool{"oberon": true},
			"dingleberry",
			map[string]bool{"oberon": true},
		},
		{
			map[string]bool{"oberon": true},
			"MoRdReD",
			map[string]bool{"mordred": true, "oberon": true},
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.enabled
		config.EnableOption(test.option)
		if !reflect.DeepEqual(config.OptionsEnabled, test.want) {
			t.Errorf("expected %v, got %v", test.want, config.OptionsEnabled)
		}
	}
}

func TestEnableMany(t *testing.T) {
	tests := []struct {
		enabled map[string]bool
		options []string
		want    map[string]bool
	}{
		{
			map[string]bool{},
			[]string{"mordred"},
			map[string]bool{"mordred": true},
		},
		{
			map[string]bool{"oberon": true},
			[]string{"oberon"},
			map[string]bool{"oberon": true},
		},
		{
			map[string]bool{"oberon": true},
			[]string{"dingleberry", "mordred"},
			map[string]bool{"oberon": true, "mordred": true},
		},
		{
			map[string]bool{"oberon": true},
			[]string{"mordred", "lake"},
			map[string]bool{"mordred": true, "oberon": true, "lake": true},
		},
		{
			map[string]bool{"oberon": false},
			[]string{"mordred", "lake"},
			map[string]bool{"mordred": true, "oberon": false, "lake": true},
		},
		{
			map[string]bool{"oberon": false},
			[]string{"oberon", "lake"},
			map[string]bool{"oberon": true, "lake": true},
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.enabled
		config.EnableMany(test.options)
		if !reflect.DeepEqual(config.OptionsEnabled, test.want) {
			t.Errorf("expected %v, got %v", test.want, config.OptionsEnabled)
		}
	}
}

func TestDisableOption(t *testing.T) {
	tests := []struct {
		enabled map[string]bool
		option  string
		want    map[string]bool
	}{
		{
			map[string]bool{},
			"mordred",
			map[string]bool{},
		},
		{
			map[string]bool{"mordred": true},
			"oberon",
			map[string]bool{"mordred": true},
		},
		{
			map[string]bool{"lake": true, "mordred": true},
			"dingleberry",
			map[string]bool{"lake": true, "mordred": true},
		},
		{
			map[string]bool{"oberon": true},
			"OBerON",
			map[string]bool{},
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.enabled
		config.DisableOption(test.option)
		if !reflect.DeepEqual(config.OptionsEnabled, test.want) {
			t.Errorf("expected %v, got %v", test.want, config.OptionsEnabled)
		}
	}
}

func TestDisableMany(t *testing.T) {
	tests := []struct {
		enabled map[string]bool
		options []string
		want    map[string]bool
	}{
		{
			map[string]bool{},
			[]string{"mordred"},
			map[string]bool{},
		},
		{
			map[string]bool{"mordred": true},
			[]string{"oberon"},
			map[string]bool{"mordred": true},
		},
		{
			map[string]bool{"lake": true, "mordred": true},
			[]string{"dingleberry", "mordred"},
			map[string]bool{"lake": true},
		},
		{
			map[string]bool{"oberon": true},
			[]string{"mordred", "oberon"},
			map[string]bool{},
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.enabled
		config.DisableMany(test.options)
		if !reflect.DeepEqual(config.OptionsEnabled, test.want) {
			t.Errorf("expected %v, got %v", test.want, config.OptionsEnabled)
		}
	}
}

func TestNumEvilSpecials(t *testing.T) {
	var tests = []struct {
		options map[string]bool
		want    int
	}{
		{
			map[string]bool{},
			0,
		},
		{
			map[string]bool{"mordred": true},
			1,
		},
		{
			map[string]bool{"mordred": true, "morganapercival": true},
			2,
		},
		{
			map[string]bool{"mordred": false, "oberon": true},
			1,
		},
		{
			map[string]bool{"mordred": false, "morganapercival": true, "oberon": true},
			2,
		},
		{
			map[string]bool{"mordred": true, "morganapercival": true, "oberon": true},
			3,
		},
		{
			map[string]bool{"mordred": false, "morganapercival": true, "lake": true},
			1,
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.options

		n := config.NumEvilSpecials()

		if n != test.want {
			t.Errorf("wanted %d for options %v, got %d", test.want, test.options, n)
		}
	}
}

// * Lake requires 7 players
// * Oberon requires 10 players
// * Must reserve at least one evil slot for assassin
func TestIsValid(t *testing.T) {
	var tests = []struct {
		numPlayers int
		options    map[string]bool
		want       bool
	}{
		{
			6,
			map[string]bool{"lake": true},
			false,
		},
		{
			7,
			map[string]bool{"lake": true},
			true,
		},
		{
			9,
			map[string]bool{"oberon": true},
			false,
		},
		{
			10,
			map[string]bool{"oberon": true},
			true,
		},
		{
			6,
			map[string]bool{"mordred": true, "morganapercival": true},
			false,
		},
		{
			7,
			map[string]bool{"mordred": true, "morganapercival": true},
			true,
		},
		{
			9,
			map[string]bool{"mordred": true, "morganapercival": true},
			true,
		},
		{
			10,
			map[string]bool{"mordred": true, "morganapercival": true, "oberon": true},
			true,
		},
	}

	for _, test := range tests {
		config := NewAvalonConfig()
		config.OptionsEnabled = test.options

		res := config.IsValid(test.numPlayers)
		if res != test.want {
			t.Errorf("expected %t for %d players, options %v, got %t", test.want, test.numPlayers, test.options, res)
		}
	}
}
