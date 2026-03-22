package holt

import "math"

// HoltForecast computes the Holt's Linear Trend forecast for a given series.
// Returns the fitted/forecasted values, and the final Level and Trend.
func HoltForecast(data []float64, alpha, beta float64) ([]float64, float64, float64) {
	if len(data) == 0 {
		return nil, 0, 0
	}

	n := len(data)
	forecasts := make([]float64, n)

	L := data[0]
	T := 0.0
	if n > 1 {
		T = data[1] - data[0]
	}

	forecasts[0] = L

	for i := 1; i < n; i++ {
		forecasts[i] = L + T
		prevL := L
		L = alpha*data[i] + (1-alpha)*(prevL+T)
		T = beta*(L-prevL) + (1-beta)*T
	}

	return forecasts, L, T
}

func MSE(data []float64, forecasts []float64) float64 {
	if len(data) == 0 || len(data) != len(forecasts) {
		return math.MaxFloat64
	}
	sum := 0.0
	for i := range data {
		diff := data[i] - forecasts[i]
		sum += diff * diff
	}
	return sum / float64(len(data))
}

// OptimizeHolt finds optimal alpha and beta via numerical gradient descent
func OptimizeHolt(data []float64, epochs int, lr float64) (bestAlpha, bestBeta float64) {
	alpha := 0.5
	beta := 0.5
	h := 1e-4

	for range epochs {
		fAlphaPlus, _, _ := HoltForecast(data, alpha+h, beta)
		fAlphaMinus, _, _ := HoltForecast(data, alpha-h, beta)
		gradAlpha := (MSE(data, fAlphaPlus) - MSE(data, fAlphaMinus)) / (2 * h)

		fBetaPlus, _, _ := HoltForecast(data, alpha, beta+h)
		fBetaMinus, _, _ := HoltForecast(data, alpha, beta-h)
		gradBeta := (MSE(data, fBetaPlus) - MSE(data, fBetaMinus)) / (2 * h)

		alpha -= lr * gradAlpha
		beta -= lr * gradBeta

		alpha = math.Max(0.001, math.Min(0.999, alpha))
		beta = math.Max(0.001, math.Min(0.999, beta))
	}

	return alpha, beta
}
