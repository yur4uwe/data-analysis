package analysis

import (
	"errors"
	"math"
)

// SolvePolynomialFit finds the least squares polynomial fit of given degree.
// Returns coefficients [a0, a1, a2, ..., an] where y = a0 + a1*x + a2*x^2 + ... + an*x^n.
// Uses the standard normal equations method (sum of powers).
func SolvePolynomialFit(xVals []float64, yVals []float64, degree int) ([]float64, error) {
	if len(xVals) != len(yVals) || len(xVals) == 0 {
		return nil, errors.New("invalid input")
	}
	if degree < 0 {
		return nil, errors.New("invalid degree")
	}
	if degree >= len(xVals) {
		return nil, errors.New("not enough data points")
	}

	n := degree + 1

	// Precompute sums of x powers: powers_of_x[power] = sum(x_k^p) for power = 0 to 2*degree
	powers_of_x := make([]float64, 2*degree+1)
	for _, x := range xVals {
		x_power_sum := 1.0
		for power := 0; power <= 2*degree; power++ {
			powers_of_x[power] += x_power_sum
			x_power_sum *= x
		}
	}

	// Precompute sums of y * x powers: B[i] = sum(y_k * x_k^i) for i = 0 to degree
	B := make([]float64, n)
	for k, x := range xVals {
		y := yVals[k]
		x_power_sum := 1.0
		for i := range n {
			B[i] += y * x_power_sum
			x_power_sum *= x
		}
	}

	// Build matrix M: M[i][j] = powers_of_x[i+j]
	M := make([][]float64, n)
	for i := range n {
		M[i] = make([]float64, n)
		for j := range n {
			M[i][j] = powers_of_x[i+j]
		}
	}

	return gaussianElimination(M, B)
}

// gaussianElimination solves the system Ax = b using Gaussian elimination with partial pivoting.
func gaussianElimination(A [][]float64, b []float64) ([]float64, error) {
	n := len(b)
	if len(A) != n {
		return nil, errors.New("matrix dimensions mismatch")
	}

	// Create augmented matrix [A|b].
	aug := make([][]float64, n)
	for i := range aug {
		aug[i] = make([]float64, n+1)
		copy(aug[i], A[i])
		aug[i][n] = b[i]
	}

	// Forward elimination with partial pivoting.
	for col := range n {
		// Find pivot.
		maxRow := col
		maxVal := math.Abs(aug[col][col])
		for row := col + 1; row < n; row++ {
			if absVal := math.Abs(aug[row][col]); absVal > maxVal {
				maxVal = absVal
				maxRow = row
			}
		}

		// Check for singular matrix.
		if maxVal < 1e-12 {
			return nil, errors.New("matrix is singular")
		}

		// Swap rows if needed.
		if maxRow != col {
			aug[col], aug[maxRow] = aug[maxRow], aug[col]
		}

		// Eliminate column entries below pivot.
		for row := col + 1; row < n; row++ {
			factor := aug[row][col] / aug[col][col]
			for j := col; j <= n; j++ {
				aug[row][j] -= factor * aug[col][j]
			}
		}
	}

	// Back substitution.
	x := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = aug[i][n]
		for j := i + 1; j < n; j++ {
			x[i] -= aug[i][j] * x[j]
		}
		x[i] /= aug[i][i]
	}

	return x, nil
}
