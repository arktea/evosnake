package persist

import (
	"os"
	"time"
	"strings"
	"fmt"
	"encoding/json"
)

const DT_LAYOUT = "20060102-150405"

type Model struct {
	Config  []int     `json:"config"`
	Weights []float64 `json:"weights"`
}

func NewModel(config []int, weights []float64) *Model {
	return &Model{Config: config, Weights: weights}
}

func Save(name string, model *Model) {
	f, err := os.Create(getFileName(name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(model); err != nil {
		panic(err)
	}
}

func Load(filename string) *Model {
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

func getFileName(name string) (filename string) {
	return fmt.Sprintf("%v-%v.model.json", strings.TrimSpace(name), time.Now().Format(DT_LAYOUT))
}