package nndriver

import (
	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/nn"
)

func PlaySnakes(rounds int, models ...*nn.Model) *game.Game {
	nnDrivers := make([]game.Driver, len(models))
	for i, model := range models {
		nnDrivers[i] = NewNNDriver(model)
	}
	g := game.NewGame(50, 50, 5, len(models), 1)
	g.Run(rounds, 50, true, nnDrivers...)
	return g
}

func PlayOneSnakeMultiGames(rounds int, nGames int, model *nn.Model) []*game.Game {
	multiDriver := NewMultiNNDriver(model)
	games := make([]*game.Game, nGames)
	for i := range games {
		games[i] = game.NewGame(50, 50, 20, 1, 1)
	}
	game.RunMulti(games, rounds, multiDriver)
	return games
}
