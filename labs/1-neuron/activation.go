package neuron

import "math"

func newSigmoid(alpha float64) (func(float64) float64, func(float64) float64) {
	f := func(x float64) float64 {
		return 1 / (1 + math.Exp(-alpha*x))
	}
	df := func(y float64) float64 {
		return alpha * y * (1 - y)
	}
	return f, df
}

func newTahn(alpha float64) (func(float64) float64, func(float64) float64) {
	f := func(x float64) float64 {
		numerator := math.Exp(alpha*x) - math.Exp(-alpha*x)
		denominator := math.Exp(alpha*x) + math.Exp(-alpha*x)
		return numerator / denominator
	}
	df := func(y float64) float64 {
		return alpha * (1 - y*y)
	}
	return f, df
}

func newReLU(alpha float64) (func(float64) float64, func(float64) float64) {
	f := func(x float64) float64 {
		if x > 0 {
			return x
		}
		return alpha * x
	}
	df := func(y float64) float64 {
		if y > 0 {
			return 1
		}
		return alpha
	}
	return f, df
}
