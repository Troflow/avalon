package avalon

import (
	"reflect"
	"strings"
	"testing"
)

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

func TestNumEvils(t *testing.T) {
	var tests = []struct {
		players []string
		want    int
	}{
		{
			[]string{"A", "B", "C", "D", "E"},
			2,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			2,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G"},
			3,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H"},
			3,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			3,
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
			4,
		},
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.Players = test.players

		n, err := avalon.NumEvils()
		if err != nil {
			t.Errorf("wanted no error, got: %v", err)
		}

		if n != test.want {
			t.Errorf("wanted %d for %d players, got %d", test.want, avalon.NumPlayers(), n)
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

		n, err := avalon.NumGoods()
		if err != nil {
			t.Errorf("wanted no error, got: %v", err)
		}

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
			[]string{"No options are enabled"},
		},
		{
			map[string]bool{"mordred": true, "lake": true, "oberon": true},
			[]string{"mordred", "lake", "oberon"},
		},
		{
			map[string]bool{"mordred": true, "lake": true, "oberon": false},
			[]string{"mordred", "lake"},
		},
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.EnabledOptions = test.options
		res := avalon.ListEnabledOptions()

		for _, want := range test.wants {
			if !strings.Contains(res, want) {
				t.Errorf("expected %s, but was absent from %s", want, res)
			}
		}
	}
}

func TestEnableOptions(t *testing.T) {
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
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.EnabledOptions = test.enabled
		avalon.EnableOptions(test.options)
		if !reflect.DeepEqual(avalon.EnabledOptions, test.want) {
			t.Errorf("expected %v, got %v", test.want, avalon.EnabledOptions)
		}
	}
}
