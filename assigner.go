package avalon

import (
	"math/rand"
)

// Assigner handles randomly assigning character roles to all players in a game.
// Used by constructing one with an Avalon game and calling Assign(). The
// Assigner will automatically populate the Avalon's fields.
type Assigner struct {
	Avalon *Avalon
}

// NewAssigner creates a new Assigner for the given Avalon game.
func NewAssigner(avalon *Avalon) *Assigner {
	return &Assigner{
		Avalon: avalon,
	}
}

// Assign should be called directly after creating a new Assigner. It populates
// the member Avalon game directly with the assignments.
func (ass *Assigner) Assign() {
	ass.assignGoodEvil()
	ass.assignSpecials()
	ass.assignFirstLeaderAndLake()
}

func (ass *Assigner) assignGoodEvil() {
	randomOrder := rand.Perm(ass.Avalon.NumPlayers())

	for i, n := range randomOrder {
		nick := ass.Avalon.Players[n]
		if i < ass.Avalon.NumGoods() {
			ass.Avalon.Goods = append(ass.Avalon.Goods, nick)
		} else {
			ass.Avalon.Evils = append(ass.Avalon.Evils, nick)
		}
	}
}

func (ass *Assigner) assignSpecials() {
	// ==========
	// == Good ==
	// ==========
	randomOrder := rand.Perm(ass.Avalon.NumGoods())

	// Merlin
	ass.Avalon.Specials["merlin"] = ass.Avalon.Goods[randomOrder[0]]

	// Percival
	if ass.Avalon.IsOptionEnabled("morganapercival") {
		ass.Avalon.Specials["percival"] = ass.Avalon.Goods[randomOrder[1]]
	}

	// ==========
	// == Evil ==
	// ==========
	randIndex := 1
	randomOrder = rand.Perm(ass.Avalon.NumEvils())

	// Assassin
	ass.Avalon.Specials["assassin"] = ass.Avalon.Evils[randomOrder[0]]

	// Mordred
	if ass.Avalon.IsOptionEnabled("mordred") {
		ass.Avalon.Specials["mordred"] = ass.Avalon.Evils[randomOrder[randIndex]]
		randIndex++
	}

	// Morgana
	if ass.Avalon.IsOptionEnabled("morganapercival") {
		ass.Avalon.Specials["morgana"] = ass.Avalon.Evils[randomOrder[randIndex]]
		randIndex++
	}

	// Oberon
	if ass.Avalon.IsOptionEnabled("oberon") {
		ass.Avalon.Specials["oberon"] = ass.Avalon.Evils[randomOrder[randIndex]]
		randIndex++
	}
}

func (ass *Assigner) assignFirstLeaderAndLake() {
	// Assign random first leader
	randLeader := rand.Intn(ass.Avalon.NumPlayers())
	ass.Avalon.CurrentLeader = ass.Avalon.Players[randLeader]

	// Lady of the Lake, cannot be the first quest leader
	if ass.Avalon.IsOptionEnabled("lake") {
		randLake := rand.Intn(ass.Avalon.NumPlayers())
		for randLake != randLeader {
			randLake = rand.Intn(ass.Avalon.NumPlayers())
		}

		ass.Avalon.CurrentLake = ass.Avalon.Players[randLake]
	}
}
