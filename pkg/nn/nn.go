package nn

import (
	// "errors"
	"math"
	// "math/rand"
	// "time"

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

func NewNN(config *NeuralNetConfig) *NeuralNet {
	var weights []*mat.Dense
	for i := range config.layers[:len(config.layers)-1] {
		in, out := config.layers[i], config.layers[i+1]
		weights = append(weights, mat.NewDense(in+1, out, nil))
	}
	return &NeuralNet{weights: weights}
}

func (nn *NeuralNet) GetRawWeights() []float64 {
	var rawWeights []float64
	for _, dense := range nn.weights {
		rawWeights = append(rawWeights, dense.RawMatrix().Data...)
	}
	return rawWeights
}

func (nn *NeuralNet) InitFromRawWeights(weights []float64) {
	var index int
	for _, dense := range nn.weights {
		data := dense.RawMatrix().Data
		for i := range data {
			data[i] = weights[index]
			index++
		}
	}
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func (nn *NeuralNet) Predict(x *mat.Dense) *mat.Dense {

	output := x
	sigm := func(_, _ int, v float64) float64 {return sigmoid(v)}
	for _, dense := range nn.weights {
		r, c := dense.Dims()
		temp := new(mat.Dense)
		temp.Mul(output, dense.Slice(0, r-1, 0, c))
		output = temp
		temp = new(mat.Dense)
		temp.Add(output, dense.Slice(r-1, r, 0, c))
		output = temp
		temp = new(mat.Dense)
		temp.Apply(sigm, output)
		output = temp
	}
	return output


	// // Complete the feed forward process.
	// hiddenLayerInput := new(mat.Dense)
	// hiddenLayerInput.Mul(x, nn.wHidden)
	// addBHidden := func(_, col int, v float64) float64 { return v + nn.bHidden.At(0, col) }
	// hiddenLayerInput.Apply(addBHidden, hiddenLayerInput)

	// hiddenLayerActivations := new(mat.Dense)
	// applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
	// hiddenLayerActivations.Apply(applySigmoid, hiddenLayerInput)

	// outputLayerInput := new(mat.Dense)
	// outputLayerInput.Mul(hiddenLayerActivations, nn.wOut)
	// addBOut := func(_, col int, v float64) float64 { return v + nn.bOut.At(0, col) }
	// outputLayerInput.Apply(addBOut, outputLayerInput)
	// output.Apply(applySigmoid, outputLayerInput)

	// return output, nil
}

// func (nn *NeuralNet) InitFromSlice(s []float64) {
// 	wHiddenRaw := nn.wHidden.RawMatrix().Data
// 	bHiddenRaw := nn.bHidden.RawMatrix().Data
// 	wOutRaw := nn.wOut.RawMatrix().Data
// 	bOutRaw := nn.bOut.RawMatrix().Data
// 	index := 0
// 	for _, param := range [][]float64{
// 		wHiddenRaw,
// 		bHiddenRaw,
// 		wOutRaw,
// 		bOutRaw,
// 	} {
// 		for i := range param {
// 			param[i] = s[index]
// 			index++
// 		}
// 	}
// }

// func (nn *NeuralNet) InitRandom() {
// 	randSource := rand.NewSource(time.Now().UnixNano())
// 	randGen := rand.New(randSource)
// 	s := make([]float64, nn.totalNeurons())
// 	for i := range s {
// 		s[i] = (randGen.Float64()*2)-1
// 	}
// 	nn.InitFromSlice(s)
// }