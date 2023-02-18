package nn

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

type NeuralNetConfig struct {
	layers []int
}

type NeuralNet struct {
	weights []*mat.Dense
}

func NewNNConfig(layers ...int) *NeuralNetConfig {
	return &NeuralNetConfig{layers: layers}
}

func (c *NeuralNetConfig) RawSize() (size int) {
	for i := range c.layers[:len(c.layers)-1] {
		size += (c.layers[i]+1)*c.layers[i+1]
	}
	return
}

func NewNN(config *NeuralNetConfig, rawWeights []float64) *NeuralNet {
	flowSize := len(config.layers)-1
	weights := make([]*mat.Dense, flowSize)
	var index int
	for i := range config.layers[:flowSize] {
		in, out := config.layers[i], config.layers[i+1]
		weights[i] = mat.NewDense(in+1, out, rawWeights[index:index+(in+1)*out])
		index += (in+1)*out
	}
	return &NeuralNet{weights: weights}
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func ones (n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = 1
	}
	return res
}

func (nn *NeuralNet) Predict(in ...[]float64) (out [][]float64) {
	inputs := mat.NewDense(len(in), len(in[0]), nil)
	for i := range in {
		inputs.SetRow(i, in[i])
	}
	outputs := new(mat.Dense)
	sigm := func(_, _ int, v float64) float64 {return sigmoid(v)}
	for _, dense := range nn.weights {
		inAug, outMul := new(mat.Dense), new(mat.Dense)
		outputs = new(mat.Dense)
		inAug.Augment(inputs,  mat.NewDense(len(in), 1, ones(len(in))))
		outMul.Mul(inAug, dense)
		outputs.Apply(sigm, outMul)
		inputs = outputs
	}
	out = make([][]float64, len(in))
	for i := range out {
		out[i] = outputs.RawRowView(i)
	}
	return
}
