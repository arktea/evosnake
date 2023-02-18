package nn

type ModelConfig []int

type Model struct {
	Config  ModelConfig `json:"config"`
	Weights []float64   `json:"weights"`
}

func NewModel(config []int, weights []float64) *Model {
	return &Model{Config: config, Weights: weights}
}

func (m ModelConfig) Size() (size int) {
	for i := range m[:len(m)-1] {
		size += (m[i]+1)*m[i+1]
	}
	return
}