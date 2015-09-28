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
