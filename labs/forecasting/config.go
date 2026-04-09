package forecasting

import (
	"fmt"
	"labs/analysis"
	"labs/charting"
	"labs/uncsv"
	"math"
	"os"
)

const (
	LabID = "7"

	ChartForecastID          = "forecast"
	ChartOptimalParametersID = "opt-params"

	GraphOriginalDataID    = "original-data"
	GraphTomorrowAsTodayID = "tomorrow-as-today"
	GraphTrendID           = "trend"
	GraphRelativeTrendID   = "relative-trend"
	GraphSimpleAvgID       = "simple-avg"
	GraphSlidingAvgID      = "sliding-avg"
	GraphExponentialAvgID  = "exponential-avg"

	VariableSlidingWindowID = "sliding-window"
	VariableAlphaID         = "alpha"

	DisplayMinErrorID    = "forecast-error-min"
	DisplayMaxErrorID    = "forecast-error-max"
	DisplayAvgErrorID    = "forecast-error-avg"
	DisplayMedianErrorID = "forecast-error-median"
	DisplayModeErrorID   = "forecast-error-mode"
	DisplayStdDevErrorID = "forecast-error-stddev"
)

type ExchageRateHistory struct {
	Date         []string  `csv:"Дата"`
	ExchangeRate []float64 `csv:"Офіційний курс гривні"`
}

var (
	VariableSlidingWindow = charting.MutableField{
		ID:      VariableSlidingWindowID,
		Label:   "Sliding Window Size",
		Default: 5,
		Min:     1,
		Max:     30,
		Step:    1,
		Control: charting.ControlRange,
	}

	VariableAlpha = charting.MutableField{
		ID:      VariableAlphaID,
		Label:   "Alpha (Exponential Smoothing)",
		Default: 0.3,
		Min:     0,
		Max:     1,
		Step:    0.01,
		Control: charting.ControlRange,
	}

	OriginalDataGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Original Rate",
			BorderColor: charting.ColorTeal,
			BorderWidth: 2,
			Togglable:   false,
		},
		PointRadius:     2,
		BackgroundColor: charting.ColorTransparent,
	}

	TomorrowAsTodayGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Tomorrow as Today",
			BorderColor:    charting.ColorAmber,
			BorderWidth:    1,
			Togglable:      true,
			GraphVariables: generateStatFields(GraphTomorrowAsTodayID),
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     2,
	}

	TrendGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Linear Trend",
			BorderColor:    charting.ColorBlue,
			BorderWidth:    1,
			Togglable:      true,
			GraphVariables: generateStatFields(GraphTrendID),
		},
		PointRadius:     2,
		BackgroundColor: charting.ColorTransparent,
	}

	RelativeTrendGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Relative Trend",
			BorderColor:    charting.ColorViolet,
			BorderWidth:    1,
			Togglable:      true,
			GraphVariables: generateStatFields(GraphRelativeTrendID),
		},
		PointRadius:     2,
		BackgroundColor: charting.ColorTransparent,
	}

	AverageGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Simple Average",
			BorderColor:    charting.ColorSlate,
			BorderWidth:    1,
			Togglable:      true,
			GraphVariables: generateStatFields(GraphSimpleAvgID),
		},
		PointRadius:     2,
		BackgroundColor: charting.ColorTransparent,
	}

	SlidingAvgGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Sliding Avg",
			BorderColor:    charting.ColorOrange,
			BorderWidth:    1,
			Togglable:      true,
			GraphVariables: generateStatFields(GraphSlidingAvgID),
		},
		PointRadius:     2,
		BackgroundColor: charting.ColorTransparent,
	}

	ExponentialAvgGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Exp. Smoothing",
			BorderColor:    charting.ColorRed,
			BorderWidth:    1,
			Togglable:      true,
			GraphVariables: generateStatFields(GraphExponentialAvgID),
		},
		PointRadius:     2,
		BackgroundColor: charting.ColorTransparent,
	}

	ForecastChart = charting.Chart{
		ID:          ChartForecastID,
		Title:       "Currency Exchange Rate Forecasting",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphOriginalDataID:    &OriginalDataGraph,
			GraphTomorrowAsTodayID: &TomorrowAsTodayGraph,
			GraphTrendID:           &TrendGraph,
			GraphRelativeTrendID:   &RelativeTrendGraph,
			GraphSimpleAvgID:       &AverageGraph,
			GraphSlidingAvgID:      &SlidingAvgGraph,
			GraphExponentialAvgID:  &ExponentialAvgGraph,
		},
		ChartVariables: []charting.MutableField{
			VariableSlidingWindow,
			VariableAlpha,
		},
	}

	OptimalParametersChart = charting.Chart{
		ID:          ChartOptimalParametersID,
		Title:       "Optimal Parameters Forecast",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphOriginalDataID:   &OriginalDataGraph,
			GraphSlidingAvgID:     &SlidingAvgGraph,
			GraphExponentialAvgID: &ExponentialAvgGraph,
		},
	}

	Config = charting.NewLabConfig(
		LabID,
		"Time Series Forecasting",
		map[string]*charting.Chart{
			ChartForecastID:          &ForecastChart,
			ChartAlphaToErrID:        &AlphaToErrorChart,
			ChartWinSizeToErrID:      &WinSizeToErrChart,
			ChartOptimalParametersID: &OptimalParametersChart,
		},
	)

	Metadata = Config.Lab

	exchangeRateData = (*ExchageRateHistory)(nil)
)

func generateStatFields(graphID string) []charting.MutableField {
	return []charting.MutableField{
		{ID: graphID + "-" + DisplayMinErrorID, Label: "Min Error", Control: charting.ControlNoControl},
		{ID: graphID + "-" + DisplayMaxErrorID, Label: "Max Error", Control: charting.ControlNoControl},
		{ID: graphID + "-" + DisplayAvgErrorID, Label: "Avg Error", Control: charting.ControlNoControl},
		{ID: graphID + "-" + DisplayMedianErrorID, Label: "Median Error", Control: charting.ControlNoControl},
		{ID: graphID + "-" + DisplayModeErrorID, Label: "Mode Error", Control: charting.ControlNoControl},
		{ID: graphID + "-" + DisplayStdDevErrorID, Label: "Std Dev Error", Control: charting.ControlNoControl},
	}
}

func loadExchageHistory() error {
	if exchangeRateData != nil && len(exchangeRateData.ExchangeRate) > 0 {
		return nil
	}
	f, err := os.Open("./data/lab_7_var_12.csv")
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.Comma = ','

	exchangeRateData = &ExchageRateHistory{}
	return d.Decode(exchangeRateData)
}

func updateGraphStats(dataset charting.Dataset, rates []float64, forecast []*float64) {
	errors := make([]float64, 0)
	for i := 0; i < len(rates) && i < len(forecast); i++ {
		if forecast[i] == nil {
			continue
		}
		errors = append(errors, math.Abs(rates[i]-*forecast[i]))
	}

	if len(errors) == 0 {
		return
	}

	minErr, maxErr := analysis.MinMax(errors)
	avgErr := analysis.Mean(errors)
	medErr := analysis.Median(errors)
	modErr := analysis.Mode(errors, 4)
	stdErr := analysis.StdDev(errors)

	dataset.UpdateVariableLabel(0, fmt.Sprintf("Min Error: %.4f", minErr))
	dataset.UpdateVariableLabel(1, fmt.Sprintf("Max Error: %.4f", maxErr))
	dataset.UpdateVariableLabel(2, fmt.Sprintf("Avg Error: %.4f", avgErr))
	dataset.UpdateVariableLabel(3, fmt.Sprintf("Median Error: %.4f", medErr))
	dataset.UpdateVariableLabel(4, fmt.Sprintf("Mode Error: %.4f", modErr))
	dataset.UpdateVariableLabel(5, fmt.Sprintf("Std Dev: %.4f", stdErr))
}

func RenderForecasting(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchageHistory(); err != nil {
		return res.NewError(err.Error())
	}

	slidingWindow, ok := req.GetChartVariable(ChartForecastID, VariableSlidingWindowID)
	if !ok {
		slidingWindow = VariableSlidingWindow.Default
	}

	alpha, ok := req.GetChartVariable(ChartForecastID, VariableAlphaID)
	if !ok {
		alpha = VariableAlpha.Default
	}

	rates := exchangeRateData.ExchangeRate
	n := len(rates)
	if n < 2 {
		return res.NewError("not enough data for forecasting")
	}

	copyChart := charting.CopyChart(ForecastChart)
	copyChart.Labels = exchangeRateData.Date

	// 1. Original Data
	copyChart.UpdateDataPointsForDataset(GraphOriginalDataID, charting.F64ToPoints(rates))

	// 2. Tomorrow as Today
	tatForecast := make([]*float64, n)
	tatForecast[0] = nil
	for i := 1; i < n; i++ {
		*tatForecast[i] = tomorrowAsToday(rates[i-1])
	}
	copyChart.UpdateDataPointsForDataset(GraphTomorrowAsTodayID, charting.F64PtrToPoints(tatForecast))
	updateGraphStats(copyChart.Datasets[GraphTomorrowAsTodayID], rates, tatForecast)

	// 3. Trend
	trendForecast := make([]*float64, n)
	trendForecast[0] = nil
	trendForecast[1] = nil
	for i := 2; i < n; i++ {
		*trendForecast[i] = trend(rates[i-1], rates[i-2])
	}
	copyChart.UpdateDataPointsForDataset(GraphTrendID, charting.F64PtrToPoints(trendForecast))
	updateGraphStats(copyChart.Datasets[GraphTrendID], rates, trendForecast)

	// 4. Relative Trend
	relTrendForecast := make([]*float64, n)
	relTrendForecast[0] = nil
	relTrendForecast[1] = nil
	for i := 2; i < n; i++ {
		*relTrendForecast[i] = relativeTrend(rates[i-1], rates[i-2])
	}
	copyChart.UpdateDataPointsForDataset(GraphRelativeTrendID, charting.F64PtrToPoints(relTrendForecast))
	updateGraphStats(copyChart.Datasets[GraphRelativeTrendID], rates, relTrendForecast)

	simpleAvgForecast := make([]*float64, n)
	for i := range n {
		*simpleAvgForecast[i] = simpleAvg(rates[:min(i+1, n)])
	}
	copyChart.UpdateDataPointsForDataset(GraphSimpleAvgID, charting.F64PtrToPoints(simpleAvgForecast))
	updateGraphStats(copyChart.Datasets[GraphSimpleAvgID], rates, simpleAvgForecast)

	// 5. Sliding Average
	slidingForecast := make([]*float64, n)
	win := int(slidingWindow)
	slidingForecast[0] = &rates[0]
	for i := 1; i < n; i++ {
		limit := min(i, win)
		*slidingForecast[i] = slidingAvg(rates[:i], limit)
	}
	copyChart.UpdateDataPointsForDataset(GraphSlidingAvgID, charting.F64PtrToPoints(slidingForecast))
	slidingDs := copyChart.Datasets[GraphSlidingAvgID]
	slidingDs.UpdateLabel(fmt.Sprintf("Sliding Avg (n=%d)", win))
	updateGraphStats(slidingDs, rates, slidingForecast)

	// 6. Exponential Smoothing
	expForecast := make([]*float64, n)
	expForecast[0] = &rates[0] // Initial seed
	for i := 1; i < n; i++ {
		*expForecast[i] = exponentialAvg(rates[i-1], *expForecast[i-1], alpha)
	}
	copyChart.UpdateDataPointsForDataset(GraphExponentialAvgID, charting.F64PtrToPoints(expForecast))
	expDs := copyChart.Datasets[GraphExponentialAvgID]
	expDs.UpdateLabel(fmt.Sprintf("Exp. Smoothing (α=%.2f)", alpha))
	updateGraphStats(expDs, rates, expForecast)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	return res
}

func RenderOptimal(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchageHistory(); err != nil {
		return res.NewError(err.Error())
	}

	rates := exchangeRateData.ExchangeRate
	n := len(rates)
	if n < 2 {
		return res.NewError("not enough data for forecasting")
	}

	// 1. Find optimal window size
	bestWin := 1
	minWinMSE := math.MaxFloat64
	var bestSlidingForecast []*float64

	for win := 1; win <= 30; win++ { // Max window 30
		slidingForecast := make([]float64, n)
		slidingForecast[0] = rates[0]
		for i := 1; i < n; i++ {
			limit := min(i, win)
			slidingForecast[i] = slidingAvg(rates[:i], limit)
		}
		mse := analysis.MSE(rates, slidingForecast)
		if mse < minWinMSE {
			minWinMSE = mse
			bestWin = win
			bestSlidingForecast = charting.F64ToPtr(slidingForecast)
		}
		fmt.Printf("Sliding window: win size = %d, MSE = %.5f\n", win, mse)
	}

	// 2. Find optimal alpha
	bestAlpha := 0.01
	minAlphaMSE := math.MaxFloat64
	var bestExpForecast []*float64

	for a := 1; a <= 99; a++ {
		alpha := float64(a) / 100.0
		expForecast := make([]float64, n)
		expForecast[0] = rates[0] // Initial seed
		for i := 1; i < n; i++ {
			expForecast[i] = exponentialAvg(rates[i-1], expForecast[i-1], alpha)
		}
		mse := analysis.MSE(rates, expForecast)
		if mse < minAlphaMSE {
			minAlphaMSE = mse
			bestAlpha = alpha
			bestExpForecast = charting.F64ToPtr(expForecast)
		}
		fmt.Printf("Exponential smoothing: alpha = %.2f, MSE = %.5f\n", alpha, mse)
	}

	copyChart := charting.CopyChart(OptimalParametersChart)
	copyChart.Labels = exchangeRateData.Date

	copyChart.UpdateDataPointsForDataset(GraphOriginalDataID, charting.F64ToPoints(rates))

	copyChart.UpdateDataPointsForDataset(GraphSlidingAvgID, charting.F64PtrToPoints(bestSlidingForecast))
	slidingDs := copyChart.Datasets[GraphSlidingAvgID]
	slidingDs.UpdateLabel(fmt.Sprintf("Opt. Sliding Avg (n=%d, MSE=%.4f)", bestWin, minWinMSE))
	updateGraphStats(slidingDs, rates, bestSlidingForecast)

	copyChart.UpdateDataPointsForDataset(GraphExponentialAvgID, charting.F64PtrToPoints(bestExpForecast))
	expDs := copyChart.Datasets[GraphExponentialAvgID]
	expDs.UpdateLabel(fmt.Sprintf("Opt. Exp. Smoothing (α=%.2f, MSE=%.4f)", bestAlpha, minAlphaMSE))
	updateGraphStats(expDs, rates, bestExpForecast)

	res = charting.NewRenderResponse()
	res.AddChart(ChartOptimalParametersID, &copyChart)
	return res
}

func init() {
	ForecastChart.RenderFunc = RenderForecasting
	WinSizeToErrChart.RenderFunc = RenderWinSizeErrChart
	AlphaToErrorChart.RenderFunc = RenderAlphaErrChart
	OptimalParametersChart.RenderFunc = RenderOptimal
}
