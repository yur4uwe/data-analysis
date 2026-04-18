package network

import (
	"math"
	"math/rand"
)

func forward(w [][]float64, x []float64) []float64 {
	y := make([]float64, len(w))
	for i := range w {
		for j := range x {
			y[i] += x[j] * w[i][j]
		}
	}
	return y
}

type TrainingResult struct {
	WeightsHistory [][][]float64
	EpochsTrained  uint32
	TestAccuracy   float64
	NumClusters    uint32
}

func normalize(x []float64) []float64 {
	sumSq := 0.0
	for _, v := range x {
		sumSq += v * v
	}
	norm := math.Sqrt(sumSq)
	if norm == 0.0 {
		return x
	}
	res := make([]float64, len(x))
	for i, v := range x {
		res[i] = v / norm
	}
	return res
}

func train(split []ClusterizationPoint, maxEpochs, numClusters uint32, lr float64) *TrainingResult {
	// Initialization
	w := make([][]float64, numClusters)
	for i := range w {
		w[i] = make([]float64, 2)
		for j := range w[i] {
			w[i][j] = -1.0 + rand.Float64()*2.0 // Random in [-1, 1]
		}

		// Normalize weights to unit circle
		w[i] = normalize(w[i])
	}

	res := TrainingResult{
		WeightsHistory: make([][][]float64, 0, maxEpochs),
		EpochsTrained:  0,
		NumClusters:    numClusters,
		TestAccuracy:   0.0,
	}

	for range maxEpochs {
		for _, point := range split {
			x := []float64{point.X, point.Y}
			y := forward(w, x)

			maxIdx := 0
			for j := range y {
				if y[j] > y[maxIdx] {
					maxIdx = j
				}
			}

			newW := []float64{
				w[maxIdx][0] + lr*(x[0]-w[maxIdx][0]),
				w[maxIdx][1] + lr*(x[1]-w[maxIdx][1]),
			}

			// Re-normalize to stay on unit circle
			w[maxIdx] = normalize(newW)

			// Copy current weights to history
			wCopy := make([][]float64, len(w))
			for i := range w {
				wCopy[i] = make([]float64, len(w[i]))
				copy(wCopy[i], w[i])
			}
			res.WeightsHistory = append(res.WeightsHistory, wCopy)
		}
	}

	return &res
}
