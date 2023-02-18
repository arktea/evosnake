package main

import (
	"math/rand"
	"time"

	"github.com/taebow/evosnake/pkg/game"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	g := game.NewGame(50, 50, 5, 1, 1)
	g.Run(-1, 25, true, game.NewKeyboardDriver(3))
}