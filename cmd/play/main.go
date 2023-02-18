package main

// import "github.com/taebow/evosnake/pkg/game"
import (
	"fmt"
	"github.com/taebow/evosnake/pkg/evosnake"
)

func main() {
	//model := evosnake.NewModel([]int{1, 2}, []float64{1.4, 1.2, 0.8})

	// evosnake.Save("toto", model)
	model := evosnake.Load("toto")
	fmt.Printf("%v\n", model)
	// game.PlayManual(25)
}