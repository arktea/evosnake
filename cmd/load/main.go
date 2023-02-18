package main

import (
	"os"

	"github.com/taebow/evosnake/pkg/nn"
	"github.com/taebow/evosnake/pkg/nndriver"
)

func main() {
	models := make([]*nn.Model, len(os.Args[1:]))
	for i := range models {
		models[i] = nn.LoadModel(os.Args[i+1])
	}
	nndriver.PlaySnakes(-1, models...)
}