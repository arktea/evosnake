package main

import (
	"github.com/taebow/evosnake/pkg/game"
)

func main() {
	g := game.NewGame(50, 50, 5, 1, 1)
	g.Run(-1, 25, true, game.NewKeyboardDriver(3))
}