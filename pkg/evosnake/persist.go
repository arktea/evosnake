package evosnake

import (
	"os"
	"encoding/json"
)

type Model struct {
	Config  []int     `json:"config"`
	Weights []float64 `json:"weights"`
}

func NewModel(config []int, weights []float64) *Model {
	return &Model{Config: config, Weights: weights}
}

func Save(name string, model *Model) {
	filename := name + ".json"
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(model); err != nil {
		panic(err)
	}
}

func Load(name string) *Model {
	filename := name + ".json"
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var model Model
	if err := json.NewDecoder(f).Decode(&model); err != nil {
		panic(err)
	}
	return &model
}
