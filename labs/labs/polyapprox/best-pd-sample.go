package polyapprox

import (
	"errors"
	"fmt"
	"labs/charting"
	"labs/labs/render"
	"labs/uncsv"
	"math"
	"os"
)

const (
	SampleMSEID     = "sample-mse"
	SampleMSEDataID = "mse-data"
)

var (
	SampleMSEChart = charting.Chart{
		ID:          SampleMSEID,
		Title:       "MSE vs Degree (CSV)",
		XAxisLabel:  "Polynomial Degree",
		YAxisLabel:  "Mean Squared Error",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			OriginalDataID: &mseGraph,
		},
		ChartVariables: []charting.MutableField{},
	}

	SampleMSEMetadata = SampleMSEChart.Meta()
)

func RenderSamplePolynomialMSE(req *charting.RenderRequest) *charting.RenderResponse {
	if points == nil {
		f, err := os.Open("../data/lab_3_var_12.csv")
		if err != nil {
			fmt.Println("failed to open file:", err)
			return &charting.RenderResponse{
				Error: render.NewRenderError("failed to read sample data file"),
			}
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ','
		points = &Points{}
		if err := d.Decode(points); err != nil {
			fmt.Println("failed to decode csv:", err)
			return &charting.RenderResponse{
				Error: render.NewRenderError("failed to decode sample data file"),
			}
		}
	}

	maxDegree := min(len(points.X)-1, 45)
	degrees := make([]float64, 0, maxDegree)
	errs := make([]float64, 0, maxDegree)

	for degree := range maxDegree - 1 {
		degree += 1
		coeffs, err := SolvePolynomialFit(points.X, points.Y, degree)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		degrees = append(degrees, float64(degree))
		errs = append(errs, CalculateMSE(points.X, points.Y, coeffs))
	}

	chartCopy := charting.CopyChart(RandomMSEChart)
	chartCopy.UpdatePointsForDataset(OriginalDataID, degrees, errs)

	return &charting.RenderResponse{
		Charts: map[string]charting.Chart{
			SampleMSEID: chartCopy,
		},
	}
}

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
