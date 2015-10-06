package avalon

import (
	"testing"
)

// avalon.NumEvils is not tested because it is a wrapper around NumEvils in util

func TestNumPlayers(t *testing.T) {
	var tests = []struct {
		players []string
		want    int
	}{
		{
			nil, 0,
		},
		{
			[]string{}, 0,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G"},
			7,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"},
			11,
		},
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.Players = test.players

		n := avalon.NumPlayers()
		if n != test.want {
			t.Errorf("wanted %d, got %d", test.want, n)
		}
	}
}

func TestNumGoods(t *testing.T) {
	var tests = []struct {
		players []string
		want    int
	}{
		{
			[]string{"A", "B", "C", "D", "E"},
			3,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			4,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G"},
			4,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H"},
			5,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			6,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
			6,
		},
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.Players = test.players

		n := avalon.NumGoods()

		if n != test.want {
			t.Errorf("wanted %d for %d players, got %d", test.want, avalon.NumPlayers(), n)
		}
	}
}

func TestPlayerExists(t *testing.T) {
	var tests = []struct {
		players []string
		search  string
		want    bool
	}{
		{
			nil,
			"Justin",
			false,
		},
		{
			[]string{},
			"Justin",
			false,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			"Z",
			false,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			"A",
			true,
		},
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.Players = test.players

		res := avalon.PlayerExists(test.search)
		if res != test.want {
			t.Errorf("wanted %t, got %t", test.want, res)
		}
	}
}

func TestAddPlayer(t *testing.T) {
	avalon := NewAvalon()

	if avalon.NumPlayers() != 0 {
		t.Errorf("expected 0 players, got %d", avalon.NumPlayers())
	}

	_ = avalon.AddPlayer("Justin")
	if avalon.NumPlayers() != 1 {
		t.Errorf("expected 1 players, got %d", avalon.NumPlayers())
	}

	if !avalon.PlayerExists("Justin") {
		t.Errorf("expected true for PlayerExists, got false")
	}

	err := avalon.AddPlayer("Justin")
	if err == nil {
		t.Error("expected error, got no error")
	}

	avalon.Players = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	err = avalon.AddPlayer("Justin")
	if err == nil {
		t.Error("expected error, got no error")
	}
}

func TestEvilsWithoutSpecial(t *testing.T) {
	var tests = []struct {
		evils    []string
		specials map[string]string
		enabled  map[string]bool
		without  string
		want     []string
	}{
		{
			[]string{},
			map[string]string{},
			map[string]bool{},
			"mordred",
			[]string{},
		},
		{
			[]string{"A"},
			map[string]string{"mordred": "A"},
			map[string]bool{"mordred": true},
			"mordred",
			[]string{},
		},
		{
			[]string{"A", "B"},
			map[string]string{"mordred": "A"},
			map[string]bool{"mordred": true},
			"mordred",
			[]string{"B"},
		},
		{
			[]string{"A", "B"},
			map[string]string{"mordred": "A", "oberon": "B"},
			map[string]bool{"mordred": true, "oberon": true},
			"oberon",
			[]string{"A"},
		},
		{
			[]string{"A", "B"},
			map[string]string{"mordred": "A", "oberon": "B"},
			map[string]bool{"mordred": true, "oberon": true},
			"morgana",
			[]string{"A", "B"},
		},
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.Evils = test.evils
		avalon.Specials = test.specials
		avalon.OptionsEnabled = test.enabled

		list := avalon.EvilsWithoutSpecial(test.without)
		if !setsEqual(list, test.want) {
			t.Errorf("wanted %v, got %v", test.want, list)
			t.Error(test)
		}
	}
}
