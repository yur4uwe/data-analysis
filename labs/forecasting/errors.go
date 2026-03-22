package forecasting

import (
	"fmt"
	"labs/charting"
	"math"
	"strconv"
)

const (
	ChartAlphaToErrID   = "alpha-err-rel"
	ChartWinSizeToErrID = "winsize-err-rel"

	GraphMinErrID = "min-err"
	GraphMaxErrID = "max-err"
	GraphAvgErrID = "avg-err"
	GraphModErrID = "mode-err"
	GraphMedErrID = "median-err"
	GraphStdErrID = "stddev-err"
)

var (
	AlphaToErrorChart = charting.Chart{
		ID:          ChartAlphaToErrID,
		Title:       "Error vs Alpha (Exponential Smoothing)",
		Type:        charting.ChartTypeMultiBar,
		XAxisLabel:  "Alpha",
		YAxisLabel:  "Error",
		XAxisConfig: charting.CategoryAxis,
		YAxisConfig: charting.LinearAxis,
	}

	WinSizeToErrChart = charting.Chart{
		ID:          ChartWinSizeToErrID,
		Title:       "Error vs Window Size (Sliding Average)",
		Type:        charting.ChartTypeMultiBar,
		XAxisLabel:  "Window size",
		YAxisLabel:  "Error",
		XAxisConfig: charting.CategoryAxis,
		YAxisConfig: charting.LinearAxis,
	}

	ErrorGraphBase = charting.ChartDataset{
		BackgroundColor: []string{charting.ColorAmber},
		BorderWidth:     0,
		PointRadius:     0,
	}
)

func toPointLabels(data []any) []string {
	labels := make([]string, len(data))
	for i, val := range data {
		labels[i] = strconv.FormatFloat(val.(float64), 'f', 4, 64)
	}
	return labels
}

func RenderAlphaErrChart(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchageHistory(); err != nil {
		return res.NewError(err.Error())
	}

	rates := exchangeRateData.ExchangeRate
	n := len(rates)
	if n < 2 {
		return res.NewError("not enough data for forecasting")
	}

	alphas := []float64{0.01, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.99}
	labels := make([]string, len(alphas))
	for i, a := range alphas {
		labels[i] = fmt.Sprintf("%.2f", a)
	}

	metrics := map[string][]any{
		GraphMinErrID: make([]any, 0, len(alphas)),
		GraphMaxErrID: make([]any, 0, len(alphas)),
		GraphAvgErrID: make([]any, 0, len(alphas)),
		GraphMedErrID: make([]any, 0, len(alphas)),
		GraphModErrID: make([]any, 0, len(alphas)),
		GraphStdErrID: make([]any, 0, len(alphas)),
	}

	for _, alpha := range alphas {
		expForecast := make([]any, n)
		expForecast[0] = rates[0]
		for i := 1; i < n; i++ {
			expForecast[i] = exponentialAvg(rates[i-1], expForecast[i-1].(float64), alpha)
		}

		errors := make([]float64, 0)
		for i := range n {
			val := expForecast[i].(float64)
			errors = append(errors, math.Abs(rates[i]-val))
		}

		minErr, maxErr := CalculateMinMax(errors)
		metrics[GraphMinErrID] = append(metrics[GraphMinErrID], minErr)
		metrics[GraphMaxErrID] = append(metrics[GraphMaxErrID], maxErr)
		metrics[GraphAvgErrID] = append(metrics[GraphAvgErrID], CalculateMean(errors))
		metrics[GraphMedErrID] = append(metrics[GraphMedErrID], CalculateMedian(errors))
		metrics[GraphModErrID] = append(metrics[GraphModErrID], CalculateMode(errors, 4))
		metrics[GraphStdErrID] = append(metrics[GraphStdErrID], CalculateStdDev(errors))
	}

	copyChart := charting.CopyChart(AlphaToErrorChart)
	copyChart.Labels = labels
	copyChart.Datasets = createErrorDatasets(metrics)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	return res
}

func RenderWinSizeErrChart(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchageHistory(); err != nil {
		return res.NewError(err.Error())
	}

	rates := exchangeRateData.ExchangeRate
	n := len(rates)
	if n < 2 {
		return res.NewError("not enough data for forecasting")
	}

	windows := []int{1, 2, 3, 5, 7, 10, 14, 21, 30}
	labels := make([]string, len(windows))
	for i, w := range windows {
		labels[i] = strconv.Itoa(w)
	}

	metrics := map[string][]any{
		GraphMinErrID: make([]any, 0, len(windows)),
		GraphMaxErrID: make([]any, 0, len(windows)),
		GraphAvgErrID: make([]any, 0, len(windows)),
		GraphMedErrID: make([]any, 0, len(windows)),
		GraphModErrID: make([]any, 0, len(windows)),
		GraphStdErrID: make([]any, 0, len(windows)),
	}

	for _, win := range windows {
		slidingForecast := make([]any, n)
		slidingForecast[0] = rates[0]
		for i := 1; i < n; i++ {
			limit := min(i, win)
			slidingForecast[i] = slidingAvg(rates[:i], limit)
		}

		errors := make([]float64, 0)
		for i := range n {
			val := slidingForecast[i].(float64)
			errors = append(errors, math.Abs(rates[i]-val))
		}

		minErr, maxErr := CalculateMinMax(errors)
		metrics[GraphMinErrID] = append(metrics[GraphMinErrID], minErr)
		metrics[GraphMaxErrID] = append(metrics[GraphMaxErrID], maxErr)
		metrics[GraphAvgErrID] = append(metrics[GraphAvgErrID], CalculateMean(errors))
		metrics[GraphMedErrID] = append(metrics[GraphMedErrID], CalculateMedian(errors))
		metrics[GraphModErrID] = append(metrics[GraphModErrID], CalculateMode(errors, 4))
		metrics[GraphStdErrID] = append(metrics[GraphStdErrID], CalculateStdDev(errors))
	}

	copyChart := charting.CopyChart(WinSizeToErrChart)
	copyChart.Labels = labels
	copyChart.Datasets = createErrorDatasets(metrics)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	return res
}

func createErrorDatasets(metrics map[string][]any) map[string]*charting.ChartDataset {
	datasets := make(map[string]*charting.ChartDataset)

	names := map[string]string{
		GraphMinErrID: "Min Error",
		GraphMaxErrID: "Max Error",
		GraphAvgErrID: "Avg Error",
		GraphMedErrID: "Median Error",
		GraphModErrID: "Mode Error",
		GraphStdErrID: "Std Dev Error",
	}

	colors := map[string]string{
		GraphMinErrID: charting.ColorTeal,
		GraphMaxErrID: charting.ColorRed,
		GraphAvgErrID: charting.ColorAmber,
		GraphMedErrID: charting.ColorBlue,
		GraphModErrID: charting.ColorViolet,
		GraphStdErrID: charting.ColorOrange,
	}

	for id, data := range metrics {
		ds := ErrorGraphBase
		ds.Label = names[id]
		ds.Data = data
		ds.BackgroundColor = []string{colors[id]}
		ds.PointLabels = toPointLabels(data)
		datasets[id] = &ds
	}

	return datasets
}
