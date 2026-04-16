package optimizations

import (
	"math"
	"math/rand/v2"
	"slices"
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

func gradient(f func(...float64) float64, pos []float64, eps float64) ([]float64, bool) {
	dims := len(pos)
	grad := make([]float64, dims)
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
	}

	return grad, gradIsInvalid
}

func dot(a, b []float64) float64 {
	if len(a) != len(b) {
		panic("dot: slices must be of equal length")
	}
	res := 0.0
	for i := range a {
		res += a[i] * b[i]
	}
	return res
}

// fastDescentNdim returns the path of positions within given bounds
func fastDescentNdim(f func(...float64) float64, initPos []float64, bounds [][]float64, eps, tol, learningRate float64) [][]float64 {
	dims := len(initPos)
	pos := make([]float64, dims)
	copy(pos, initPos)

	// Initial clamp
	for i := range dims {
		if pos[i] < bounds[i][0] {
			pos[i] = bounds[i][0]
		}
		if pos[i] > bounds[i][1] {
			pos[i] = bounds[i][1]
		}
	}

	path := [][]float64{append([]float64(nil), pos...)}
	lastVal := f(pos...)

	if isInvalid(lastVal) {
		return path // Can't start from NaN/Inf
	}

	for range maxIter {
		grad, gradIsInvalid := gradient(f, pos, eps)
		if gradIsInvalid {
			learningRate *= 0.5
			if learningRate < 1e-10 {
				break
			}
			continue
		}

		if dot(grad, grad) < tol*tol {
			break
		}

		// Trial step with clamping
		newPos := make([]float64, dims)
		for i := range dims {
			newPos[i] = pos[i] - grad[i]*learningRate
			if newPos[i] < bounds[i][0] {
				newPos[i] = bounds[i][0]
			}
			if newPos[i] > bounds[i][1] {
				newPos[i] = bounds[i][1]
			}
		}

		// If new clamped positions is the same as the old, exit
		if slices.Equal(newPos, pos) {
			break
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
		newGrad, newGradIsInvalid := gradient(f, newPos, eps)
		if newGradIsInvalid {
			learningRate *= 0.5
		} else {
			dotProduct := dot(grad, newGrad)
			if dotProduct < 0 {
				learningRate *= 0.5
			} else {
				learningRate *= 1.1 // Smaller than 1.6 in original code
				// It was producing NaN/Inf values in the gradient vector
			}
		}

		copy(pos, newPos)
		lastVal = newVal
		path = append(path, append([]float64(nil), pos...))
	}
	return path
}
