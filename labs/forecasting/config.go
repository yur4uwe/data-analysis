package forecasting

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"math"
	"os"
	"strings"
)

const (
	LabID = "7"

	ChartForecastID = "forecast"

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

	OriginalDataGraph = charting.ChartDataset{
		Label:           "Original Rate",
		BorderColor:     charting.ColorTeal,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       false,
	}

	TomorrowAsTodayGraph = charting.ChartDataset{
		Label:           "Tomorrow as Today",
		BorderColor:     charting.ColorAmber,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     1,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  generateStatFields(GraphTomorrowAsTodayID),
	}

	TrendGraph = charting.ChartDataset{
		Label:           "Linear Trend",
		BorderColor:     charting.ColorBlue,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     1,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  generateStatFields(GraphTrendID),
	}

	RelativeTrendGraph = charting.ChartDataset{
		Label:           "Relative Trend",
		BorderColor:     charting.ColorViolet,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     1,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  generateStatFields(GraphRelativeTrendID),
	}

	AverageGraph = charting.ChartDataset{
		Label:           "Simple Average",
		BorderColor:     charting.ColorSlate,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     1,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  generateStatFields(GraphSimpleAvgID),
	}

	SlidingAvgGraph = charting.ChartDataset{
		Label:           "Sliding Avg",
		BorderColor:     charting.ColorOrange,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     1,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  generateStatFields(GraphSlidingAvgID),
	}

	ExponentialAvgGraph = charting.ChartDataset{
		Label:           "Exp. Smoothing",
		BorderColor:     charting.ColorRed,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     1,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables:  generateStatFields(GraphExponentialAvgID),
	}

	ForecastChart = charting.Chart{
		ID:          ChartForecastID,
		Title:       "Currency Exchange Rate Forecasting",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
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

	Config = charting.NewLabConfig(
		LabID,
		"Time Series Forecasting",
		map[string]*charting.Chart{
			ChartForecastID:     &ForecastChart,
			ChartAlphaToErrID:   &AlphaToErrorChart,
			ChartWinSizeToErrID: &WinSizeToErrChart,
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

func updateGraphStats(dataset *charting.ChartDataset, rates []float64, forecast []any) {
	errors := make([]float64, 0)
	for i := 0; i < len(rates) && i < len(forecast); i++ {
		if forecast[i] == nil {
			continue
		}
		val, ok := forecast[i].(float64)
		if !ok {
			continue
		}
		errors = append(errors, math.Abs(rates[i]-val))
	}

	if len(errors) == 0 {
		return
	}

	minErr, maxErr := CalculateMinMax(errors)
	avgErr := CalculateMean(errors)
	medErr := CalculateMedian(errors)
	modErr := CalculateMode(errors, 4)
	stdErr := CalculateStdDev(errors)

	for i := range dataset.GraphVariables {
		field := &dataset.GraphVariables[i]
		switch {
		case strings.HasSuffix(field.ID, DisplayMinErrorID):
			field.Label = fmt.Sprintf("Min Error: %.4f", minErr)
		case strings.HasSuffix(field.ID, DisplayMaxErrorID):
			field.Label = fmt.Sprintf("Max Error: %.4f", maxErr)
		case strings.HasSuffix(field.ID, DisplayAvgErrorID):
			field.Label = fmt.Sprintf("Avg Error: %.4f", avgErr)
		case strings.HasSuffix(field.ID, DisplayMedianErrorID):
			field.Label = fmt.Sprintf("Median Error: %.4f", medErr)
		case strings.HasSuffix(field.ID, DisplayModeErrorID):
			field.Label = fmt.Sprintf("Mode Error: %.4f", modErr)
		case strings.HasSuffix(field.ID, DisplayStdDevErrorID):
			field.Label = fmt.Sprintf("Std Dev: %.4f", stdErr)
		}
	}
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
	copyChart.UpdateDataForDataset(GraphOriginalDataID, charting.ToAnySlice(rates))

	// 2. Tomorrow as Today
	tatForecast := make([]any, n)
	for i := 1; i < n; i++ {
		tatForecast[i] = tomorrowAsToday(rates[i-1])
	}
	copyChart.UpdateDataForDataset(GraphTomorrowAsTodayID, tatForecast)
	updateGraphStats(copyChart.Datasets[GraphTomorrowAsTodayID], rates, tatForecast)

	// 3. Trend
	trendForecast := make([]any, n)
	for i := 2; i < n; i++ {
		trendForecast[i] = trend(rates[i-1], rates[i-2])
	}
	copyChart.UpdateDataForDataset(GraphTrendID, trendForecast)
	updateGraphStats(copyChart.Datasets[GraphTrendID], rates, trendForecast)

	// 4. Relative Trend
	relTrendForecast := make([]any, n)
	for i := 2; i < n; i++ {
		relTrendForecast[i] = relativeTrend(rates[i-1], rates[i-2])
	}
	copyChart.UpdateDataForDataset(GraphRelativeTrendID, relTrendForecast)
	updateGraphStats(copyChart.Datasets[GraphRelativeTrendID], rates, relTrendForecast)

	simpleAvgForecast := make([]any, n)
	for i := range n {
		simpleAvgForecast[i] = simpleAvg(rates[:min(i+1, n)])
	}
	copyChart.UpdateDataForDataset(GraphSimpleAvgID, simpleAvgForecast)
	updateGraphStats(copyChart.Datasets[GraphSimpleAvgID], rates, simpleAvgForecast)

	fmt.Println("Simple avg forecast:", simpleAvgForecast)

	// 5. Sliding Average
	slidingForecast := make([]any, n)
	win := int(slidingWindow)
	slidingForecast[0] = rates[0]
	for i := 1; i < n; i++ {
		limit := min(i, win)
		slidingForecast[i] = slidingAvg(rates[:i], limit)
	}
	copyChart.UpdateDataForDataset(GraphSlidingAvgID, slidingForecast)
	slidingDs := copyChart.Datasets[GraphSlidingAvgID]
	slidingDs.Label = fmt.Sprintf("Sliding Avg (n=%d)", win)
	updateGraphStats(slidingDs, rates, slidingForecast)

	// 6. Exponential Smoothing
	expForecast := make([]any, n)
	expForecast[0] = rates[0] // Initial seed
	for i := 1; i < n; i++ {
		expForecast[i] = exponentialAvg(rates[i-1], expForecast[i-1].(float64), alpha)
	}
	copyChart.UpdateDataForDataset(GraphExponentialAvgID, expForecast)
	expDs := copyChart.Datasets[GraphExponentialAvgID]
	expDs.Label = fmt.Sprintf("Exp. Smoothing (α=%.2f)", alpha)
	updateGraphStats(expDs, rates, expForecast)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	return res
}

func init() {
	ForecastChart.RenderFunc = RenderForecasting
	WinSizeToErrChart.RenderFunc = RenderWinSizeErrChart
	AlphaToErrorChart.RenderFunc = RenderAlphaErrChart
}
