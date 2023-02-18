package main

import (
	"os"
	"time"
	"math/rand"

	"github.com/taebow/evosnake/pkg/nn"
	"github.com/taebow/evosnake/pkg/game"
	"github.com/taebow/evosnake/pkg/nndriver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	model := nn.LoadModel(os.Args[1])
	nnDriver := nndriver.NewNNDriver(model)
	g := game.NewGame(50, 50, 5, 1, 1)
	g.Run(-1, 100, true, nnDriver)
}