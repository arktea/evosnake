package nn

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

type NeuralNetConfig struct {
	inputNeurons, hiddenNeurons, outputNeurons int
}

type NeuralNet struct {
	wHidden, bHidden, wOut, bOut *mat.Dense
	config *NeuralNetConfig
}

func NewNNConfig(inputNeurons, outputNeurons, hiddenNeurons int) *NeuralNetConfig {
	return &NeuralNetConfig{
		inputNeurons: inputNeurons, 
		outputNeurons: outputNeurons, 
		hiddenNeurons: hiddenNeurons,
	}
}

func NewNN(config *NeuralNetConfig) *NeuralNet {
	wHidden := mat.NewDense(config.inputNeurons, config.hiddenNeurons, nil)
	bHidden := mat.NewDense(1, config.hiddenNeurons, nil)
	wOut := mat.NewDense(config.hiddenNeurons, config.outputNeurons, nil)
	bOut := mat.NewDense(1, config.outputNeurons, nil)
	return &NeuralNet{wHidden: wHidden, bHidden: bHidden, wOut: wOut, bOut: bOut, config: config}
}

func (nn *NeuralNet) totalNeurons() (res int) {
	for _, dense := range []*mat.Dense{nn.wHidden, nn.bHidden, nn.wOut, nn.bOut} {
		r, c := dense.Dims()
		res += r*c
	}
	return
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func (nn *NeuralNet) Predict(x *mat.Dense) (*mat.Dense, error) {

	// Check to make sure that our neuralNet value
	// represents a trained model.
	if nn.wHidden == nil || nn.wOut == nil {
		return nil, errors.New("the supplied weights are empty")
	}
	if nn.bHidden == nil || nn.bOut == nil {
		return nil, errors.New("the supplied biases are empty")
	}

	// Define the output of the neural network.
	output := new(mat.Dense)

	// Complete the feed forward process.
	hiddenLayerInput := new(mat.Dense)
	hiddenLayerInput.Mul(x, nn.wHidden)
	addBHidden := func(_, col int, v float64) float64 { return v + nn.bHidden.At(0, col) }
	hiddenLayerInput.Apply(addBHidden, hiddenLayerInput)

	hiddenLayerActivations := new(mat.Dense)
	applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
	hiddenLayerActivations.Apply(applySigmoid, hiddenLayerInput)

	outputLayerInput := new(mat.Dense)
	outputLayerInput.Mul(hiddenLayerActivations, nn.wOut)
	addBOut := func(_, col int, v float64) float64 { return v + nn.bOut.At(0, col) }
	outputLayerInput.Apply(addBOut, outputLayerInput)
	output.Apply(applySigmoid, outputLayerInput)

	return output, nil
}

func (nn *NeuralNet) ToSlice() []float64 {
	wHiddenRaw := nn.wHidden.RawMatrix().Data
	bHiddenRaw := nn.bHidden.RawMatrix().Data
	wOutRaw := nn.wOut.RawMatrix().Data
	bOutRaw := nn.bOut.RawMatrix().Data
	res := make([]float64, nn.totalNeurons())
	i := 0
	for _, param := range [][]float64{wHiddenRaw, bHiddenRaw, wOutRaw, bOutRaw} {
		for _, p := range param {
			res[i] = p
			i++
		}
	}
	return res
}

func (nn *NeuralNet) InitFromSlice(s []float64) {
	wHiddenRaw := nn.wHidden.RawMatrix().Data
	bHiddenRaw := nn.bHidden.RawMatrix().Data
	wOutRaw := nn.wOut.RawMatrix().Data
	bOutRaw := nn.bOut.RawMatrix().Data
	index := 0
	for _, param := range [][]float64{
		wHiddenRaw,
		bHiddenRaw,
		wOutRaw,
		bOutRaw,
	} {
		for i := range param {
			param[i] = s[index]
			index++
		}
	}
}

func (nn *NeuralNet) InitRandom() {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	s := make([]float64, nn.totalNeurons())
	for i := range s {
		s[i] = (randGen.Float64()*2)-1
	}
	nn.InitFromSlice(s)
}