package optimizations

import (
	"math"
	"math/rand/v2"
)

const (
	maxIter = 100
)

func dichotomicSearch(f func(float64) float64, a, b, tol float64) []float64 {
	path := []float64{(a + b) / 2}
	// epsilon must be significantly smaller than tol
	// typically eps < tol/2. If tol=0.01, eps=0.001 is better.
	eps := tol / 4
	for range maxIter {
		if math.Abs(a-b) < tol {
			break
		}
		p := (a+b)/2 + eps
		q := (a+b)/2 - eps

		if f(p) < f(q) {
			a = q
		} else {
			b = p
		}
		path = append(path, (a+b)/2)
	}
	return path
}

// randomSearchNdim finds the extremum for a function with N variables.
// bounds is an Nx2 slice where each element is [min, max] for that dimension.
func randomSearchNdim(f func(...float64) float64, nSamples int, bounds [][]float64) (minPath, maxPath [][]float64) {
	dims := len(bounds)
	minVal := math.MaxFloat64
	maxVal := -math.MaxFloat64
	
	currMin := make([]float64, dims)
	currMax := make([]float64, dims)

	for range nSamples {
		point := make([]float64, dims)
		for d := range dims {
			point[d] = bounds[d][0] + rand.Float64()*(bounds[d][1]-bounds[d][0])
		}

		val := f(point...)
		if val < minVal {
			minVal = val
			copy(currMin, point)
			p := make([]float64, dims)
			copy(p, currMin)
			minPath = append(minPath, p)
		}
		if val > maxVal {
			maxVal = val
			copy(currMax, point)
			p := make([]float64, dims)
			copy(p, currMax)
			maxPath = append(maxPath, p)
		}
	}
	return minPath, maxPath
}

func isInvalid(f float64) bool {
	return math.IsNaN(f) || math.IsInf(f, 0)
}

// fastDescentNdim returns the path of positions
func fastDescentNdim(f func(...float64) float64, initPos []float64, eps, tol, learningRate float64) [][]float64 {
	dims := len(initPos)
	pos := make([]float64, dims)
	copy(pos, initPos)

	path := [][]float64{append([]float64(nil), pos...)}
	lastVal := f(pos...)

	if isInvalid(lastVal) {
		return path // Can't start from NaN/Inf
	}

	for range maxIter {
		grad := make([]float64, dims)
		var gradNormSq float64
		gradIsInvalid := false

		// Calculate gradient vector with NaN/Inf protection
		for i := range dims {
			original := pos[i]
			pos[i] = original + eps
			fPlus := f(pos...)
			pos[i] = original - eps
			fMinus := f(pos...)
			pos[i] = original // restore

			if isInvalid(fPlus) || isInvalid(fMinus) {
				gradIsInvalid = true
				break
			}

			grad[i] = (fPlus - fMinus) / (2 * eps)
			gradNormSq += grad[i] * grad[i]
		}

		if gradIsInvalid || math.IsNaN(gradNormSq) {
			learningRate *= 0.5
			if learningRate < 1e-10 {
				break
			}
			continue
		}

		gradNorm := math.Sqrt(gradNormSq)
		if gradNorm < tol {
			break
		}

		// Trial step
		newPos := make([]float64, dims)
		for i := range dims {
			newPos[i] = pos[i] - grad[i]*learningRate
		}

		newVal := f(newPos...)

		// If we jumped into NaN/Inf or value increased, reduce LR and retry from same pos
		if isInvalid(newVal) || newVal > lastVal {
			learningRate *= 0.5
			if learningRate < 1e-10 {
				break
			}
			continue
		}

		// Calculate gradient at new position for adjustment
		newGrad := make([]float64, dims)
		var dotProduct float64
		newGradIsInvalid := false
		for i := range dims {
			original := newPos[i]
			newPos[i] = original + eps
			fPlus := f(newPos...)
			newPos[i] = original - eps
			fMinus := f(newPos...)
			newPos[i] = original

			if isInvalid(fPlus) || isInvalid(fMinus) {
				newGradIsInvalid = true
				break
			}

			newGrad[i] = (fPlus - fMinus) / (2 * eps)
			dotProduct += grad[i] * newGrad[i]
		}

		if newGradIsInvalid || math.IsNaN(dotProduct) || dotProduct < 0 {
			learningRate *= 0.5
		} else {
			learningRate *= 1.1
		}

		copy(pos, newPos)
		lastVal = newVal
		path = append(path, append([]float64(nil), pos...))
	}
	return path
}
