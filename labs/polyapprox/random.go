package polyapprox

import (
	"fmt"
	"labs/charting"
	"math/rand"
)

const (
	LabID = "3"

	RandomFitsID = "random-fits"

	OriginalDataID = "orig-data"
	NoisyDataID    = "noisy-data"
	LinearApproxID = "linear-approx"
	QuadApproxID   = "quad-approx"

	IntervalStartID = "start-interval"
	IntervalEndID   = "end-interval"
	IntervalStepID  = "step"
	NoiseAmpID      = "noise-amplifier"

	LinearFitCoefficientsID    = "linear-fit-coefficients"
	QuadraticFitCoefficientsID = "quadratic-fit-coefficients"
)

var (
	ChartVariables = []charting.MutableField{
		{
			ID:      IntervalStartID,
			Label:   "Start",
			Default: -5.0,
			Min:     -100.0,
			Max:     100.0,
			Step:    0.5,
			Control: charting.ControlNumber,
		},
		{
			ID:      IntervalEndID,
			Label:   "End",
			Default: 5.0,
			Min:     -100.0,
			Max:     100.0,
			Step:    0.5,
			Control: charting.ControlNumber,
		},
		{
			ID:      IntervalStepID,
			Label:   "Step",
			Default: 0.1,
			Min:     0.01,
			Max:     1,
			Step:    0.01,
			Control: charting.ControlRange,
		},
		{
			ID:      NoiseAmpID,
			Label:   "Noise Amplifier",
			Default: 10,
			Min:     1,
			Max:     100,
			Step:    1,
			Control: charting.ControlRange,
		},
	}

	LinearFitCoefficients = charting.MutableField{
		ID:      LinearFitCoefficientsID,
		Label:   "Linear Fit Coefficients",
		Control: charting.ControlNoControl,
	}

	QuadraticFitCoefficients = charting.MutableField{
		ID:      QuadraticFitCoefficientsID,
		Label:   "Quadratic Fit Coefficients",
		Control: charting.ControlNoControl,
	}

	OriginalData = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Original",
			BorderColor: charting.ColorAmber,
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorAmber,
	}

	NoisyData = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Noisy",
			BorderColor: charting.ColorCyan,
			BorderWidth: 1,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorCyan,
		PointRadius:     2,
	}

	LinearApprox = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Linear Approximation",
			BorderColor:    "#16a34a",
			BorderWidth:    2,
			Togglable:      true,
			GraphVariables: []charting.MutableField{LinearFitCoefficients},
		},
		BackgroundColor: "rgba(22, 163, 74, 0.1)",
	}

	QuadApprox = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Quadratic Approximation",
			BorderColor:    "#9333ea",
			BorderWidth:    2,
			Togglable:      true,
			GraphVariables: []charting.MutableField{QuadraticFitCoefficients},
		},
		BackgroundColor: "rgba(147, 51, 234, 0.1)",
		PointRadius:     0,
	}

	RandomFitsChart = charting.Chart{
		ID:          RandomFitsID,
		Title:       "Random Data Fits",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			OriginalDataID: &OriginalData,
			NoisyDataID:    &NoisyData,
			LinearApproxID: &LinearApprox,
			QuadApproxID:   &QuadApprox,
		},
		ChartVariables: ChartVariables,
	}

	RandomFitsMetadata = RandomFitsChart.Meta()

	Metadata = charting.LabMetadata{
		ID:   LabID,
		Name: "Least Squares Approximation",
		Charts: map[string]charting.ChartMetadata{
			RandomFitsID: RandomFitsMetadata,
			SampleDataID: SampleDataMetadata,
			RandomMSEID:  RandomMSEMetadata,
			SampleMSEID:  SampleMSEMetadata,
		},
	}
)

func RenderRandomFits(req *charting.RenderRequest) (res *charting.RenderResponse) {
	start, hasStart := req.GetChartVariable(RandomFitsID, IntervalStartID)
	end, hasEnd := req.GetChartVariable(RandomFitsID, IntervalEndID)
	step, hasStep := req.GetChartVariable(RandomFitsID, IntervalStepID)
	noiseAmp, hasNoise := req.GetChartVariable(RandomFitsID, NoiseAmpID)

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
		return res.NewError("step must be greater than 0")
	}
	if start > end {
		return res.NewError("start interval must be less than or equal to end interval")
	}

	seed := int64(230420067)
	x, y, origY := GenerateRandomSeries(start, end, step, noiseAmp, seed)

	if len(x) == 0 {
		return res.NewError("no data generated with given parameters")
	}

	chartCopy := charting.CopyChart(RandomFitsChart)

	chartCopy.UpdatePointsForDataset(OriginalDataID, x, origY)
	chartCopy.UpdatePointsForDataset(NoisyDataID, x, y)

	if coefs, err := SolvePolynomialFit(x, y, 1); err == nil {
		approx := make([]float64, 0, len(x))
		for _, xi := range x {
			approx = append(approx, EvaluatePolynomial(coefs, xi))
		}
		mse := CalculateMSE(x, y, coefs)
		fmt.Printf("Linear Fit MSE: %.4e\n", mse)
		chartCopy.UpdatePointsForDataset(LinearApproxID, x, approx)
		chartCopy.Datasets[LinearApproxID].UpdateVariableLabel(0, fmt.Sprintf("Linear Fit (MSE: %.4e, a=%.4f, b=%.4f) for y=bx+a", mse, coefs[0], coefs[1]))
	} else {
		fmt.Println("linear fit failed:", err)
	}

	if coefs, err := SolvePolynomialFit(x, y, 2); err == nil {
		approx := make([]float64, 0, len(x))
		for _, xi := range x {
			approx = append(approx, EvaluatePolynomial(coefs, xi))
		}
		mse := CalculateMSE(x, y, coefs)
		fmt.Printf("Quadratic Fit MSE: %.4e\n", mse)
		chartCopy.UpdatePointsForDataset(QuadApproxID, x, approx)
		chartCopy.Datasets[QuadApproxID].UpdateVariableLabel(0, fmt.Sprintf("Quadratic Fit (MSE: %.4e, a=%.4f, b=%.4f, c=%.4f) for y=cx^2+bx+a", mse, coefs[0], coefs[1], coefs[2]))
	} else {
		fmt.Println("quadratic fit failed:", err)
	}

	res = charting.NewRenderResponse()
	res.AddChart(RandomFitsID, &chartCopy)

	return res
}

func GenerateRandomSeries(start, end, step, noiseAmp float64, seed int64) ([]float64, []float64, []float64) {
	r := rand.New(rand.NewSource(seed))
	n := int((end-start)/step) + 1
	x := make([]float64, 0, n)
	y := make([]float64, 0, n)
	origY := make([]float64, 0, n)

	for i := start; i <= end; i += step {
		noise := r.NormFloat64() * 0.2 * noiseAmp
		curr := 0.8 - 4*i
		x = append(x, i)
		y = append(y, curr+noise)
		origY = append(origY, curr)
	}

	return x, y, origY
}
