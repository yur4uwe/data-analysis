package neuron

import (
	"math"
	"math/rand"
)

type TrainingResult struct {
	WeightsHistory [][]float64
	BiasHistory    []float64
	LossHistory    []float64
}

func train(data *ClassificationData, epochs uint32, lr float64, activation func(float64) float64, actDerivative func(float64) float64) TrainingResult {
	points := make([][]float64, len(data.X))
	for i := range data.X {
		points[i] = []float64{data.X[i], data.Y[i]}
	}

	// Initialization
	w := []float64{rand.Float64() - 0.5, rand.Float64() - 0.5}
	b := 0.0

	res := TrainingResult{
		WeightsHistory: make([][]float64, epochs),
		BiasHistory:    make([]float64, epochs),
		LossHistory:    make([]float64, epochs),
	}

	for epoch := range epochs {
		totalLoss := 0.0
		indices := rand.Perm(len(points))

		for _, pointIdx := range indices {
			z := w[0]*points[pointIdx][0] + w[1]*points[pointIdx][1] + b
			pred := activation(z)

			target := 0.0
			if data.Class[pointIdx] {
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

			loss := 0.5 * (pred - target) * (pred - target)
			if math.IsNaN(loss) || math.IsInf(loss, 0) {
				loss = 1.0 // Penalty for explosion
			}
			totalLoss += loss
		}

		// Check for weight explosion
		if math.IsNaN(w[0]) || math.IsNaN(w[1]) || math.IsNaN(b) || math.IsInf(w[0], 0) || math.IsInf(w[1], 0) || math.IsInf(b, 0) {
			w[0], w[1], b = 0, 0, 0
		}

		// Store snapshot for visualizations - ensure we never save NaN/Inf to history
		safeW0, safeW1, safeB := w[0], w[1], b
		safeLoss := totalLoss / float64(len(points))
		
		if math.IsNaN(safeLoss) || math.IsInf(safeLoss, 0) {
			safeLoss = 1.0
		}

		res.WeightsHistory[epoch] = []float64{safeW0, safeW1}
		res.BiasHistory[epoch] = safeB
		res.LossHistory[epoch] = safeLoss
	}

	return res
}
