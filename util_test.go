package avalon

import (
	"testing"
)

func TestOptionExists(t *testing.T) {
	var tests = []struct {
		option string
		want   bool
	}{
		{"justin", false},
		{"lake", true},
		{"mordred", true},
		{"morganapercival", true},
		{"oberon", true},
	}

	for _, test := range tests {
		res := OptionExists(test.option)
		if res != test.want {
			t.Errorf("expected %t, got %t", test.want, res)
		}
	}
}

// * Lake requires 7 players
// * Oberon requires 10 players
// * Must reserve at least one evil slot for assassin
func TestOptionsValid(t *testing.T) {
	var tests = []struct {
		players []string
		options map[string]bool
		want    bool
	}{
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			map[string]bool{"lake": true},
			false,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G"},
			map[string]bool{"lake": true},
			true,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			map[string]bool{"oberon": true},
			false,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
			map[string]bool{"oberon": true},
			true,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			map[string]bool{"mordred": true, "morganapercival": true},
			false,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G"},
			map[string]bool{"mordred": true, "morganapercival": true},
			true,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			map[string]bool{"mordred": true, "morganapercival": true},
			true,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
			map[string]bool{"mordred": true, "morganapercival": true, "oberon": true},
			true,
		},
	}

	for _, test := range tests {
		av := NewAvalon()
		av.Players = test.players
		av.EnabledOptions = test.options

		res := OptionsValid(av)
		if res != test.want {
			t.Errorf("expected %t for players %v, options %v, got %t", test.want, test.players, test.options, res)
		}
	}
}

func TestNumEvils(t *testing.T) {
	var tests = []struct {
		numPlayers int
		want       int
	}{
		{0, 0},
		{1, 0},
		{4, 0},
		{5, 2},
		{6, 2},
		{7, 3},
		{8, 3},
		{9, 3},
		{10, 4},
		{11, 0},
	}

	for _, test := range tests {
		n := NumEvils(test.numPlayers)

		if n != test.want {
			t.Errorf("wanted %d for %d players, got %d", test.want, test.numPlayers, n)
		}
	}
}
