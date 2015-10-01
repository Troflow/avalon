package avalon

import (
	"sort"
	"testing"
)

func TestAssign(t *testing.T) {
	tests := []struct {
		players []string
		options map[string]bool
	}{
		{
			[]string{"A", "B", "C", "D", "E"},
			map[string]bool{},
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"},
			map[string]bool{"lake": true, "mordred": true, "morganapercival": true, "oberon": true},
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"},
			map[string]bool{"lake": true, "mordred": true, "morganapercival": true},
		},
		{
			[]string{"A", "B", "C", "D", "E", "F", "G"},
			map[string]bool{"lake": true, "mordred": true, "morganapercival": true},
		},
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			map[string]bool{"mordred": true},
		},
		{
			[]string{"A", "B", "C", "D", "E", "F"},
			map[string]bool{"morganapercival": true},
		},
	}

	for _, test := range tests {
		avalon := NewAvalon()
		avalon.Players = test.players
		avalon.AvalonConfig.OptionsEnabled = test.options

		ass := NewAssigner(avalon)
		ass.Assign()

		runRequirementTests(t, avalon)
	}
}

func runRequirementTests(t *testing.T, av *Avalon) {
	numPlayers := len(av.Players)
	numGoodsAssigned := len(av.Goods)
	numEvilsAssigned := len(av.Evils)
	numPlayersAssigned := numGoodsAssigned + numEvilsAssigned

	// Number of goods is what we should have
	if numGoodsAssigned != av.NumGoods() {
		t.Errorf("expected %d goods assigned, got %d", av.NumGoods(), numGoodsAssigned)
	}

	// Number of evils is what we should have
	if numEvilsAssigned != av.NumEvils() {
		t.Errorf("expected %d evils assigned, got %d", av.NumEvils(), numEvilsAssigned)
	}

	// goods + evils = number of players
	if numPlayersAssigned != numPlayers {
		t.Errorf("expected %d players assigned to good and evil, got %d", numPlayers, numPlayersAssigned)
	}

	// goods ∪ evils = players
	if !setsEqual(append(av.Goods, av.Evils...), av.Players) {
		t.Errorf("union of goods and evils != players; goods: %v, evils: %v, players: %v", av.Goods, av.Evils, av.Players)
	}

	// Goods and evils should be mutually exclusive
	if intersects(av.Goods, av.Evils) {
		t.Errorf("player appears in goods and evils: %v, %v", av.Goods, av.Evils)
	}

	// None of these sets should have duplicates
	if containsDuplicates(av.Goods) {
		t.Errorf("goods contained duplicates: %v", av.Goods)
	}
	if containsDuplicates(av.Evils) {
		t.Errorf("evils contained duplicates: %v", av.Evils)
	}

	// No player should be assigned to two specials
	if valueAppearsTwice(av.Specials) {
		t.Errorf("specials contained duplicates: %v", av.Specials)
	}

	// Merlin and Assassin must always be assigned
	if !keyExists(av.Specials, "merlin") {
		t.Errorf("merlin not assigned: %v", av.Specials)
	}
	if !keyExists(av.Specials, "assassin") {
		t.Errorf("assassin not assigned: %v", av.Specials)
	}

	// If mordred is enabled, Mordred should be assigned
	if av.IsOptionEnabled("mordred") && !keyExists(av.Specials, "mordred") {
		t.Errorf("mordred enabled, no mordred assigned")
	}

	// If morganapercival is enabled, Morgana and Percival should be assigned
	if av.IsOptionEnabled("morganapercival") {
		if !keyExists(av.Specials, "morgana") {
			t.Errorf("morganapercival enabled, no morgana assigned")
		}
		if !keyExists(av.Specials, "percival") {
			t.Errorf("morganapercival enabled, no percival assigned")
		}
	}

	// If oberon is enabled, Oberon should be assigned
	if av.IsOptionEnabled("oberon") {
		if !keyExists(av.Specials, "oberon") {
			t.Errorf("oberon enabled, no oberon assigned")
		}
	}

	// All specials should be from the list of players
	if !valuesAreSubsetOf(av.Specials, av.Players) {
		t.Errorf("special assigned that isn't a player: %v", av.Specials)
	}

	// First leader should be a player
	if !av.PlayerExists(av.CurrentLeader) {
		t.Errorf("first leader is not a player: %s", av.CurrentLeader)
	}

	// Lake should be a player
	if av.IsOptionEnabled("lake") {
		if !av.PlayerExists(av.CurrentLeader) {
		}
	}

	// If lake is enabled, The Lady of the Lake should be assigned
	if av.IsOptionEnabled("lake") {
		if av.CurrentLake == "" {
			t.Errorf("lake enabled, no lake assigned; players: %v, options: %v", av.Players, av.OptionsEnabled)
		}
	}

	// First leader should never be the same as lake
	if av.IsOptionEnabled("lake") {
		if av.CurrentLeader == av.Specials["lake"] {
			t.Errorf("first leader and lake are same: %s", av.CurrentLeader)
		}
	}
}

func setsEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	copyA := make([]string, len(a))
	copy(copyA, a)
	sort.Strings(copyA)

	copyB := make([]string, len(b))
	copy(copyB, b)
	sort.Strings(copyB)

	for i := range copyA {
		if copyA[i] != copyB[i] {
			return false
		}
	}
	return true
}

func intersects(a []string, b []string) bool {
	for _, e := range a {
		for _, ee := range b {
			if e == ee {
				return true
			}
		}
	}
	return false
}

func containsDuplicates(a []string) bool {
	for i, e := range a {
		for ii, ee := range a {
			if e == ee && i != ii {
				return true
			}
		}
	}
	return false
}

func valueAppearsTwice(a map[string]string) bool {
	for k, v := range a {
		for kk, vv := range a {
			if v == vv && k != kk {
				return true
			}
		}
	}
	return false
}

func keyExists(m map[string]string, key string) bool {
	_, ok := m[key]
	return ok
}

// For every element in a, check that it appears somewhere in b. If it does not,
// immediately return false. If it goes through the entirety of a without
// finding an element not in b, then a is a subset of b.
func isSubset(a []string, b []string) bool {
	for _, e := range a {
		var found bool
		for _, ee := range b {
			if e == ee {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func valuesAreSubsetOf(a map[string]string, b []string) bool {
	values := make([]string, 0, len(a))
	for _, v := range a {
		values = append(values, v)
	}

	return isSubset(values, b)
}
