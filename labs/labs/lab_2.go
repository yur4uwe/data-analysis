package labs

import (
	"labs/labs/common"
	"labs/labs/render"
	"math"
	"math/rand/v2"
)

const (
	Lab2ID = "2"

	mainChartID                   = "main-chart"
	originalDataID                = "orig-data"
	noiseDataID                   = "noisy-data"
	recurringAverageID            = "recurrent-avg"
	slidingWindowAverageID        = "sliding-window-avg"
	exponentialSmoothingAverageID = "exponential-average"

	intervalStartID  = "start-interval"
	intervalEndID    = "end-interval"
	intervalStepID   = "step"
	noiseAmplifierID = "noise-amplifier"
	windowSizeID     = "win-size"
	alphaID          = "alpha"
)

var (
	metadata = common.LabMetadata{
		ID:   Lab2ID,
		Name: "Primary data processing",
		Charts: map[string]common.ChartMetadata{
			mainChartID: {
				ID:             mainChartID,
				Title:          "Primary data processing",
				ChartVariables: ChartVariables,
				GraphVariables: map[string][]common.MutableField{
					originalDataID:                originalData.GraphVariables,
					noiseDataID:                   noisyData.GraphVariables,
					recurringAverageID:            recurrentAvg.GraphVariables,
					slidingWindowAverageID:        slidingWindowAvg.GraphVariables,
					exponentialSmoothingAverageID: exponentialAverage.GraphVariables,
				},
			},
		},
	}
	ChartVariables = []common.MutableField{
		{
			ID:      intervalStartID,
			Label:   "Start",
			Default: 0.0,
			Min:     -100.0,
			Max:     100.0,
			Step:    1,
			Control: common.ControlNumber,
		},
		{
			ID:      intervalEndID,
			Label:   "End",
			Default: 10.0,
			Min:     -100.0,
			Max:     100.0,
			Step:    1,
			Control: common.ControlNumber,
		},
		{
			ID:      intervalStepID,
			Label:   "Step",
			Default: 0.1,
			Min:     0.1,
			Max:     1,
			Step:    0.1,
			Control: common.ControlRange,
		},
	}
	originalData = common.ChartDataset{
		Label:           "Original",
		Data:            nil,
		BorderColor:     "#2563eb", // Blue
		BackgroundColor: []string{"rgba(37, 99, 235, 0.1)"},
		Tension:         0,
		Fill:            false,
		Hidden:          false,
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  []common.MutableField{},
	}
	noisyData = common.ChartDataset{
		Label:           "Noisy",
		Data:            nil,
		BorderColor:     "#dc2626", // Red
		BackgroundColor: []string{"rgba(220, 38, 38, 0.1)"},
		Tension:         0,
		Fill:            false,
		Hidden:          false,
		PointRadius:     2,
		BorderWidth:     1,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []common.MutableField{
			{
				ID:      noiseAmplifierID,
				Label:   "Noise Amplifier",
				Default: 1,
				Min:     0.0,
				Max:     100.0,
				Step:    1,
				Control: common.ControlRange,
			},
		},
	}
	recurrentAvg = common.ChartDataset{
		Label:           "Recurrent Average",
		Data:            nil,
		BorderColor:     "#16a34a", // Green
		BackgroundColor: []string{"rgba(22, 163, 74, 0.1)"},
		Tension:         0,
		Fill:            false,
		Hidden:          false,
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  []common.MutableField{},
	}
	slidingWindowAvg = common.ChartDataset{
		Label:           "Sliding Window Average",
		Data:            nil,
		BorderColor:     "#9333ea", // Purple
		BackgroundColor: []string{"rgba(147, 51, 234, 0.1)"},
		Tension:         0,
		Fill:            false,
		Hidden:          false,
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []common.MutableField{
			{
				ID:      windowSizeID,
				Label:   "Window Size",
				Default: 10,
				Min:     1,
				Max:     100,
				Step:    1,
				Control: common.ControlNumber,
			},
		},
	}
	exponentialAverage = common.ChartDataset{
		Label:           "Exponential Average",
		Data:            nil,
		BorderColor:     "#ea580c", // Orange
		BackgroundColor: []string{"rgba(234, 88, 12, 0.1)"},
		Tension:         0,
		Fill:            false,
		Hidden:          false,
		PointRadius:     0,
		BorderWidth:     3,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []common.MutableField{
			{
				ID:      alphaID,
				Label:   "Alpha",
				Default: 0.5,
				Min:     0,
				Max:     1,
				Step:    0.01,
				Control: common.ControlRange,
			},
		},
	}
	main = common.Chart{
		ID:          mainChartID,
		Title:       "Primary Data Processing Methods Visualization",
		Type:        common.ChartTypeLine,
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			originalDataID:                &originalData,
			noiseDataID:                   &noisyData,
			recurringAverageID:            &recurrentAvg,
			slidingWindowAverageID:        &slidingWindowAvg,
			exponentialSmoothingAverageID: &exponentialAverage,
		},
		ChartVariables: ChartVariables,
	}
)

func recurringAvg(x, prevAvg float64, len int) float64 {
	return prevAvg + (x-prevAvg)/float64(len)
}

func slidingAvg(data []float64, i int, winSize int) float64 {
	sum := 0.0
	for k := range winSize {
		sum += data[i+k]
	}
	return sum / float64(winSize)
}

func exponentialAvg(x, prevAvg float64, alpha float64) float64 {
	return prevAvg + alpha*(x-prevAvg)
}

type Lab2Provider struct {
}

var _ common.LabProvider = Lab2Provider{}

func NewLab2() *Lab2Provider {
	return &Lab2Provider{}
}

func (lp Lab2Provider) GetMetadata() common.LabMetadata {
	return metadata
}

func (lp Lab2Provider) Render(req *common.RenderRequest) *common.RenderResponse {
	interval_start, has_start := req.GetChartVariable(mainChartID, intervalStartID)
	interval_end, has_end := req.GetChartVariable(mainChartID, intervalEndID)
	interval_step, has_step := req.GetChartVariable(mainChartID, intervalStepID)
	noise_amplifier, has_noise_amplifier := req.GetGraphVariable(mainChartID, noiseDataID, noiseAmplifierID)

	if !has_start {
		interval_start = main.ChartVariables[0].Default
	}
	if !has_end {
		interval_end = main.ChartVariables[1].Default
	}
	if !has_step {
		interval_step = main.ChartVariables[2].Default
	}
	if !has_noise_amplifier {
		noise_amplifier = noisyData.GraphVariables[0].Default
	}

	// Validate input parameters
	if interval_step <= 0 {
		return &common.RenderResponse{
			Error: render.NewRenderError("step must be greater than 0"),
		}
	}
	if interval_start > interval_end {
		return &common.RenderResponse{
			Error: render.NewRenderError("start interval must be less than or equal to end interval"),
		}
	}

	// Generate data
	n := int((interval_end - interval_start + 1) / interval_step)
	x := make([]float64, 0, n)
	y := make([]float64, 0, n)
	origY := make([]float64, 0, n)

	for i := interval_start; i <= interval_end; i += interval_step {
		// gaussian distribution noise from standard normal distribution
		noise := rand.NormFloat64() * 0.2 * noise_amplifier
		f := math.Sin(0.1*i) + math.Pow(math.Sin(i), 2)
		x = append(x, i)
		y = append(y, f+noise)
		origY = append(origY, f)
	}

	// Check if we have data to work with
	if len(origY) == 0 {
		return &common.RenderResponse{
			Error: render.NewRenderError("no data generated with given parameters"),
		}
	}

	// Deep copy the chart template and update only the data
	chartCopy := common.CopyChart(main)

	// Helper function to update dataset with point data (x, y pairs)
	updateDatasetWithPoints := func(key string, xData, yData []float64) {
		dataset := chartCopy.Datasets[key]
		pointData := make([]common.DataPoint, len(yData))
		for i := range yData {
			pointData[i] = common.DataPoint{X: xData[i], Y: yData[i]}
		}
		dataset.PointData = pointData
		chartCopy.Datasets[key] = dataset
	}

	updateDatasetWithPoints(originalDataID, x, origY)
	updateDatasetWithPoints(noiseDataID, x, y)

	// Add sliding window average (with default window size if not provided)
	winSize := slidingWindowAvg.GraphVariables[0].Default
	if size, ok := req.GetGraphVariable(mainChartID, slidingWindowAverageID, windowSizeID); ok {
		winSize = size
	}
	if len(y) > int(winSize) {
		slidingWinAvg := make([]float64, 0, len(y))
		slidingWinX := make([]float64, 0, len(y))
		for i := 0; i < len(y)-int(winSize)+1; i++ {
			slidingWinAvg = append(slidingWinAvg, slidingAvg(origY, i, int(winSize)))
			slidingWinX = append(slidingWinX, x[i+int(winSize)/2]) // Center point of window
		}
		updateDatasetWithPoints(slidingWindowAverageID, slidingWinX, slidingWinAvg)
	}

	// Add recurrent average
	recAvg := make([]float64, 0, len(y))
	prevAvg := y[0]
	recAvg = append(recAvg, prevAvg)
	for i := 1; i < len(y); i++ {
		prevAvg = recurringAvg(y[i], prevAvg, i)
		recAvg = append(recAvg, prevAvg)
	}
	updateDatasetWithPoints(recurringAverageID, x, recAvg)

	// Add exponential average (with default alpha if not provided)
	alpha := exponentialAverage.GraphVariables[0].Default
	if a, ok := req.GetGraphVariable(mainChartID, exponentialSmoothingAverageID, alphaID); ok {
		alpha = a
	}
	expAvg := make([]float64, 0, len(y))
	prevAvg = y[0]
	expAvg = append(expAvg, prevAvg)
	for i := 1; i < len(y); i++ {
		prevAvg = exponentialAvg(y[i], prevAvg, alpha)
		expAvg = append(expAvg, prevAvg)
	}
	updateDatasetWithPoints(exponentialSmoothingAverageID, x, expAvg)

	return &common.RenderResponse{
		Charts: map[string]common.Chart{
			mainChartID: chartCopy,
		},
	}
}

func (lp Lab2Provider) GetConfig() common.LabConfig {
	return common.LabConfig{
		Lab: metadata,
		Charts: map[string]*common.Chart{
			mainChartID: &main,
		},
	}
}
