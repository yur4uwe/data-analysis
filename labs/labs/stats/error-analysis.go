package stats

import (
	"fmt"
	"labs/labs/common"
	"math"
	"math/rand"
)

const (
	ErrorAnalysisChartID = "error-analysis"

	StdDevErrorGraphID = "stddev-error"
	MeanErrorGraphID   = "mean-error"
)

var (
	StdDevErrorGraph = common.ChartDataset{
		Label:           "Std Dev Error",
		BorderColor:     common.Color1,
		BackgroundColor: []string{"rgba(37, 99, 235, 0.1)"},
		PointRadius:     3,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
	}

	MeanErrorGraph = common.ChartDataset{
		Label:           "Mean Error",
		BorderColor:     common.Color2,
		BackgroundColor: []string{"rgba(239, 68, 68, 0.1)"},
		PointRadius:     3,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
	}

	ErrorAnalysisChart = common.Chart{
		ID:          ErrorAnalysisChartID,
		Title:       "Error Analysis: Standard Deviation vs Mean",
		Type:        common.ChartTypeLine,
		XAxisLabel:  "Parameter Value",
		YAxisLabel:  "Error (Absolute)",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			StdDevErrorGraphID: &StdDevErrorGraph,
			MeanErrorGraphID:   &MeanErrorGraph,
		},
		ChartVariables: []common.MutableField{
			{
				ID:      "mean-for-stddev",
				Label:   "Mean (for Std Dev Analysis)",
				Default: 0.0,
				Min:     -10.0,
				Max:     10.0,
				Step:    0.1,
				Control: common.ControlNumber,
			},
			{
				ID:      "stddev-for-mean",
				Label:   "Std Dev (for Mean Analysis)",
				Default: 1.0,
				Min:     0.1,
				Max:     5.0,
				Step:    0.1,
				Control: common.ControlNumber,
			},
			{
				ID:      "min-stddev",
				Label:   "Min Std Dev",
				Default: 0.5,
				Min:     0.1,
				Max:     10.0,
				Step:    0.1,
				Control: common.ControlNumber,
			},
			{
				ID:      "max-stddev",
				Label:   "Max Std Dev",
				Default: 5.0,
				Min:     0.1,
				Max:     20.0,
				Step:    0.1,
				Control: common.ControlNumber,
			},
			{
				ID:      "stddev-step",
				Label:   "Std Dev Step",
				Default: 0.25,
				Min:     0.01,
				Max:     1.0,
				Step:    0.01,
				Control: common.ControlNumber,
			},
			{
				ID:      "min-mean",
				Label:   "Min Mean",
				Default: -5.0,
				Min:     -20.0,
				Max:     0.0,
				Step:    0.1,
				Control: common.ControlNumber,
			},
			{
				ID:      "max-mean",
				Label:   "Max Mean",
				Default: 5.0,
				Min:     0.0,
				Max:     20.0,
				Step:    0.1,
				Control: common.ControlNumber,
			},
			{
				ID:      "mean-step",
				Label:   "Mean Step",
				Default: 0.5,
				Min:     0.01,
				Max:     2.0,
				Step:    0.01,
				Control: common.ControlNumber,
			},
		},
	}

	ErrorAnalysisMeta = ErrorAnalysisChart.Meta()
)

// GenerateNormalSamples generates samples from normal distribution
func GenerateNormalSamples(mean, stddev float64, sampleCount int) []float64 {
	samples := make([]float64, sampleCount)
	for i := range sampleCount {
		samples[i] = rand.NormFloat64()*stddev + mean
	}
	return samples
}

// CalculateMean computes the mean of a slice
func CalculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// CalculateStdDev computes the standard deviation of a slice
func CalculateStdDev(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range data {
		diff := v - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(data))
	return math.Sqrt(variance)
}

// AnalyzeErrorByStdDev tests how error in stddev changes as we vary theoretical stddev
func AnalyzeErrorByStdDev(mean float64, minStdDev, maxStdDev, step float64, sampleCount int) ([]float64, []float64) {
	var xValues []float64 // theoretical stddev
	var yValues []float64 // error in stddev

	for theoreticalStdDev := minStdDev; theoreticalStdDev <= maxStdDev; theoreticalStdDev += step {
		samples := GenerateNormalSamples(mean, theoreticalStdDev, sampleCount)
		empiricalMean := CalculateMean(samples)
		empiricalStdDev := CalculateStdDev(samples, empiricalMean)

		error := math.Abs(empiricalStdDev - theoreticalStdDev)

		xValues = append(xValues, theoreticalStdDev)
		yValues = append(yValues, error)
	}

	return xValues, yValues
}

// AnalyzeErrorByMean tests how error in mean changes as we vary theoretical mean
func AnalyzeErrorByMean(stddev float64, minMean, maxMean, step float64, sampleCount int) ([]float64, []float64) {
	var xValues []float64 // theoretical mean
	var yValues []float64 // error in mean

	for theoreticalMean := minMean; theoreticalMean <= maxMean; theoreticalMean += step {
		samples := GenerateNormalSamples(theoreticalMean, stddev, sampleCount)
		empiricalMean := CalculateMean(samples)

		error := math.Abs(empiricalMean - theoreticalMean)

		xValues = append(xValues, theoreticalMean)
		yValues = append(yValues, error)
	}

	return xValues, yValues
}

func RenderErrorAnalysis(req *common.RenderRequest) *common.RenderResponse {
	fmt.Printf("Rendering %s\n", req.ChartID)

	if req == nil {
		return &common.RenderResponse{
			Error: fmt.Errorf("request is nil"),
		}
	}

	// Get parameters from request
	// Default values can be overridden from frontend
	sampleCount := 50000

	// Get mean for stddev analysis
	meanForStdDevAnalysis := 0.0
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "mean-for-stddev"); ok {
		meanForStdDevAnalysis = val
	}

	// Get stddev for mean analysis
	stddevForMeanAnalysis := 1.0
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "stddev-for-mean"); ok {
		stddevForMeanAnalysis = val
	}

	// Get range parameters
	minStdDev := 0.5
	maxStdDev := 5.0
	stdDevStep := 0.25
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "min-stddev"); ok {
		minStdDev = val
	}
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "max-stddev"); ok {
		maxStdDev = val
	}
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "stddev-step"); ok {
		stdDevStep = val
	}

	minMean := -5.0
	maxMean := 5.0
	meanStep := 0.5
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "min-mean"); ok {
		minMean = val
	}
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "max-mean"); ok {
		maxMean = val
	}
	if val, ok := req.GetChartVariable(ErrorAnalysisChartID, "mean-step"); ok {
		meanStep = val
	}

	// Analyze errors
	stdDevXs, stdDevYs := AnalyzeErrorByStdDev(meanForStdDevAnalysis, minStdDev, maxStdDev, stdDevStep, sampleCount)
	meanXs, meanYs := AnalyzeErrorByMean(stddevForMeanAnalysis, minMean, maxMean, meanStep, sampleCount)

	// Create response
	chartCopy := common.CopyChart(ErrorAnalysisChart)

	// Update stddev error data
	err1 := chartCopy.UpdatePointsForDataset(StdDevErrorGraphID, stdDevXs, stdDevYs)
	if err1 != nil {
		return &common.RenderResponse{
			Error: fmt.Errorf("failed to update stddev dataset: %v", err1),
		}
	}

	// Update mean error data
	err2 := chartCopy.UpdatePointsForDataset(MeanErrorGraphID, meanXs, meanYs)
	if err2 != nil {
		return &common.RenderResponse{
			Error: fmt.Errorf("failed to update mean dataset: %v", err2),
		}
	}

	res := common.NewRenderResponse()
	res.AddChart(ErrorAnalysisChartID, &chartCopy)

	return res
}
