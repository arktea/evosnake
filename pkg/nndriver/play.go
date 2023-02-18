package nndriver

import (
	"github.com/taebow/evosnake/pkg/nn"
	"github.com/taebow/evosnake/pkg/game"
)

func NNPlayGame(rounds int, model *nn.Model) *game.Game {
	nnDriver := NewNNDriver(model)
	g := game.NewGame(50, 50, 45, 1, 1)
	g.Run(rounds, -1, false, nnDriver)
	return g
}

func PlaySnakes(models []*nn.Model) {
	nnDrivers := make([]game.Driver, len(models))
	for i, model := range models {
		nnDrivers[i] = NewNNDriver(model)
	}
	g := game.NewGame(50, 50, 5, len(models), 1)
	g.Run(-1, 25, true, nnDrivers...)
}

func MultiPlayGames(rounds int, model *nn.Model, nGames int) []*game.Game {
	multiDriver := NewMultiDriver([]*nn.Model{model})
	games := make([]*game.Game, nGames)
	for i := range games {
		games[i] = game.NewGame(50, 50, 20, 1, 1)
	}
	game.RunMulti(games, rounds, multiDriver)
	return games
}