package polyapprox

import (
	"fmt"
	"labs/charting"
	"labs/labs/render"
)

const (
	RandomMSEID     = "random-mse"
	RandomMSEDataID = "mse-data"
)

var (
	mseGraph = charting.ChartDataset{
		Label:       "MSE vs Degree",
		BorderColor: charting.Color11,
		ShowLine:    true,
	}

	RandomMSEChart = charting.Chart{
		ID:          RandomMSEID,
		Title:       "MSE vs Degree (Random)",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Polynomial Degree",
		YAxisLabel:  "Mean Squared Error",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			OriginalDataID: &mseGraph,
		},
		ChartVariables: ChartVariables,
	}

	RandomMSEMetadata = RandomMSEChart.Meta()
)

func RenderRandomPolynomialMSE(req *charting.RenderRequest) *charting.RenderResponse {
	start, hasStart := req.GetChartVariable(RandomMSEID, IntervalStartID)
	end, hasEnd := req.GetChartVariable(RandomMSEID, IntervalEndID)
	step, hasStep := req.GetChartVariable(RandomMSEID, IntervalStepID)
	noiseAmp, hasNoise := req.GetChartVariable(RandomMSEID, NoiseAmpID)

	if !hasStart {
		start = ChartVariables[0].Default
	}
	if !hasEnd {
		end = ChartVariables[1].Default
	}
	if !hasStep {
		step = ChartVariables[2].Default
	}
	if !hasNoise {
		noiseAmp = ChartVariables[3].Default
	}

	if step <= 0 {
		return &charting.RenderResponse{Error: render.NewRenderError("step must be greater than 0")}
	}
	if start > end {
		return &charting.RenderResponse{Error: render.NewRenderError("start interval must be less than or equal to end interval")}
	}

	seed := int64(230420067)
	x, y, _ := GenerateRandomSeries(start, end, step, noiseAmp, seed)
	if len(x) == 0 {
		return &charting.RenderResponse{Error: render.NewRenderError("no data generated with given parameters")}
	}

	maxDegree := min(len(x)-1, 45)
	degrees := make([]float64, 0, maxDegree)
	errs := make([]float64, 0, maxDegree)

	for degree := range maxDegree - 1 {
		degree += 1
		coeffs, err := SolvePolynomialFit(x, y, degree)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		degrees = append(degrees, float64(degree))
		errs = append(errs, CalculateMSE(x, y, coeffs))
	}

	chartCopy := charting.CopyChart(RandomMSEChart)
	chartCopy.UpdatePointsForDataset(OriginalDataID, degrees, errs)

	return &charting.RenderResponse{
		Charts: map[string]charting.Chart{
			RandomMSEID: chartCopy,
		},
	}
}

// Uses Horner's method for evaluating polynomials.
//
// Formula:
//
//	y = (...(((a1x + a2)x + a3)x + a4) ... )x + an)
func EvaluatePolynomial(coeffs []float64, x float64) float64 {
	if len(coeffs) == 0 {
		return 0
	}

	result := coeffs[len(coeffs)-1]
	for i := len(coeffs) - 2; i >= 0; i-- {
		result = result*x + coeffs[i]
	}
	return result
}

func CalculateMSE(xVals, yVals []float64, coeffs []float64) float64 {
	if len(xVals) != len(yVals) || len(xVals) == 0 {
		return 0
	}

	sumSquaredError := 0.0
	for i := range xVals {
		predicted := EvaluatePolynomial(coeffs, xVals[i])
		error := yVals[i] - predicted
		sumSquaredError += error * error
	}

	return sumSquaredError / float64(len(xVals))
}
