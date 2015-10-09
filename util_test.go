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
