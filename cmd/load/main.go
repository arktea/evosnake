package main

import (
	"os"

	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/nn"
	"github.com/taebow/evosnake/pkg/nndriver"
	"github.com/taebow/evosnake/pkg/persist"
)

func main() {
	model := persist.Load(os.Args[1])
	nnDriver := nndriver.NewNNDriver(nn.NewNNConfig(model.Config...), model.Weights)
	g := game.NewGame(50, 50, 5, 1, 1)
	g.Run(-1, 50, true, nnDriver)
}