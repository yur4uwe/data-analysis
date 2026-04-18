package neuron

import (
	"math"
	"math/rand"
)

type TrainingResult struct {
	WeightsHistory        [][]float64
	BiasHistory           []float64
	LossHistory           []float64
	ValidationLossHistory []float64
	EpochsTrained         uint32
	TestAccuracy          float64
}

func newForward(activation func(float64) float64, w []float64, b float64) func([]float64) float64 {
	return func(z []float64) float64 {
		return activation(w[0]*z[0] + w[1]*z[1] + b)
	}
}

func calculateLoss(data *ClassificationData, w []float64, b float64, activation func(float64) float64) float64 {
	totalLoss := 0.0
	forward := newForward(activation, w, b)
	for i := range data.X {
		pred := forward([]float64{data.X[i], data.Y[i]})
		target := 0.0
		if data.Class[i] {
			target = 1.0
		}
		loss := 0.5 * (pred - target) * (pred - target)
		totalLoss += loss
	}
	return totalLoss / float64(len(data.X))
}

func train(split *ClassificationSplit, maxEpochs uint32, lr float64, targetAccuracy float64, activation func(float64) float64, actDerivative func(float64) float64) TrainingResult {
	trainSet := split.Train
	valSet := split.Validation

	points := make([][]float64, len(trainSet.X))
	for i := range trainSet.X {
		points[i] = []float64{trainSet.X[i], trainSet.Y[i]}
	}

	// Initialization
	w := []float64{rand.Float64() - 0.5, rand.Float64() - 0.5}
	b := 0.0

	res := TrainingResult{
		WeightsHistory:        make([][]float64, 0, maxEpochs),
		BiasHistory:           make([]float64, 0, maxEpochs),
		LossHistory:           make([]float64, 0, maxEpochs),
		ValidationLossHistory: make([]float64, 0, maxEpochs),
	}

	for epoch := range maxEpochs {
		indices := rand.Perm(len(points))

		for _, pointIdx := range indices {
			forward := newForward(activation, w, b)
			pred := forward(points[pointIdx])

			target := 0.0
			if trainSet.Class[pointIdx] {
				target = 1.0
			}

			errorGrad := (pred - target)
			delta := errorGrad * actDerivative(pred)

			// Safety: Clip delta to prevent explosion and check for NaN
			if math.IsNaN(delta) || math.IsInf(delta, 0) {
				delta = 0
			}

			w[0] -= lr * delta * points[pointIdx][0]
			w[1] -= lr * delta * points[pointIdx][1]
			b -= lr * delta
		}

		// Check for weight explosion
		if math.IsNaN(w[0]) || math.IsNaN(w[1]) || math.IsNaN(b) || math.IsInf(w[0], 0) || math.IsInf(w[1], 0) || math.IsInf(b, 0) {
			w[0], w[1], b = 0, 0, 0
		}

		// Calculate metrics
		avgLoss := calculateLoss(trainSet, w, b, activation)
		valLoss := calculateLoss(valSet, w, b, activation)

		if math.IsNaN(avgLoss) || math.IsInf(avgLoss, 0) {
			avgLoss = 1.0
		}
		if math.IsNaN(valLoss) || math.IsInf(valLoss, 0) {
			valLoss = 1.0
		}

		res.WeightsHistory = append(res.WeightsHistory, []float64{w[0], w[1]})
		res.BiasHistory = append(res.BiasHistory, b)
		res.LossHistory = append(res.LossHistory, avgLoss)
		res.ValidationLossHistory = append(res.ValidationLossHistory, valLoss)
		res.EpochsTrained = epoch + 1

		// Target Accuracy Stopping Condition
		if avgLoss <= targetAccuracy {
			break
		}
	}

	// Calculate Final Test Accuracy
	testSet := split.Test
	correct := 0
	forwardTest := newForward(activation, w, b)
	for i := range testSet.X {
		pred := forwardTest([]float64{testSet.X[i], testSet.Y[i]})
		predClass := pred >= 0.5
		if predClass == testSet.Class[i] {
			correct++
		}
	}
	res.TestAccuracy = float64(correct) / float64(len(testSet.X))

	return res
}
