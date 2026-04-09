package holt

import (
	"errors"
	"fmt"
	"labs/analysis"
	"labs/charting"
	"labs/uncsv"
	"math"
	"os"
)

var (
	testL float64
	testT float64
	bestA float64
	bestB float64

	hasTrained bool = false
)

func loadExchangeHistory() error {
	if len(testExchangeRateData.ExchangeRate) > 0 && len(trainExchangeRateData.ExchangeRate) > 0 {
		return nil
	}
	f, err := os.Open("./data/lab_8_var_12.csv")
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.Comma = ','

	exchangeRateData := &ExchangeRateHistory{}
	if err := d.Decode(exchangeRateData); err != nil {
		return err
	}
	n := len(exchangeRateData.ExchangeRate)
	if n < 4 {
		return errors.New("not enough data for training and testing")
	}

	splitIdx := n / 2
	trainExchangeRateData.ExchangeRate = exchangeRateData.ExchangeRate[:splitIdx]
	trainExchangeRateData.Date = exchangeRateData.Date[:splitIdx]
	testExchangeRateData.ExchangeRate = exchangeRateData.ExchangeRate[splitIdx:]
	testExchangeRateData.Date = exchangeRateData.Date[splitIdx:]

	return nil
}

func RenderHoltTest(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if !hasTrained {
		return res.NewError("Run Training on 'Holt's Method - Training Phase' chart first to determine optimal values.")
	}

	if err := loadExchangeHistory(); err != nil {
		return res.NewError(err.Error())
	}

	testData := testExchangeRateData.ExchangeRate
	testDates := testExchangeRateData.Date

	testForecasts := make([]float64, len(testData))
	testLevel := testL
	testTrend := testT

	bestAlpha := bestA
	bestBeta := bestB

	for i := range testData {
		testForecasts[i] = testLevel + testTrend
		prevL := testLevel
		testLevel = bestAlpha*testData[i] + (1-bestAlpha)*(prevL+testTrend)
		testTrend = bestBeta*(testLevel-prevL) + (1-bestBeta)*testTrend
	}

	testMSE := analysis.MSE(testData, testForecasts)

	copyTestChart := charting.CopyChart(TestChart)
	copyTestChart.Labels = testDates
	copyTestChart.UpdateDataPointsForDataset(GraphTestActualID, charting.AnyToPoints(charting.F64ToAny(testData)))
	copyTestChart.UpdateDataPointsForDataset(GraphTestForecastID, charting.AnyToPoints(charting.F64ToAny(testForecasts)))

	copyTestChart.Datasets[GraphTestForecastID].UpdateVariableLabel(0, fmt.Sprintf("Alpha Used: %.4f", bestAlpha))
	copyTestChart.Datasets[GraphTestForecastID].UpdateVariableLabel(1, fmt.Sprintf("Beta Used: %.4f", bestBeta))
	copyTestChart.Datasets[GraphTestForecastID].UpdateVariableLabel(2, fmt.Sprintf("Test MSE: %.4e", testMSE))

	res = charting.NewRenderResponse()
	res.AddChart(copyTestChart.ID, &copyTestChart)
	res.CachePolicy = charting.CachePolicyDontCache
	return res
}

func RenderHolt(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchangeHistory(); err != nil {
		return res.NewError(err.Error())
	}

	epochs, ok := req.GetChartVariable(ChartHoltTrainID, VariableEpochsID)
	if !ok {
		epochs = VariableEpochs.Default
	}

	lr, ok := req.GetChartVariable(ChartHoltTrainID, VariableLearningRateID)
	if !ok {
		lr = VariableLearningRate.Default
	}

	// Train Phase
	trainData := trainExchangeRateData.ExchangeRate
	trainDates := trainExchangeRateData.Date

	bestAlpha, bestBeta := OptimizeHolt(trainData, int(epochs), lr)
	trainForecasts, finalL, finalT := HoltForecast(trainData, bestAlpha, bestBeta)
	trainMSE := analysis.MSE(trainData, trainForecasts)

	// Save for Test phase
	testL = finalL
	testT = finalT
	bestA = bestAlpha
	bestB = bestBeta
	hasTrained = true

	copyTrainChart := charting.CopyChart(TrainChart)
	copyTrainChart.Labels = trainDates
	copyTrainChart.UpdateDataPointsForDataset(GraphTrainActualID, charting.AnyToPoints(charting.F64ToAny(trainData)))
	copyTrainChart.UpdateDataPointsForDataset(GraphTrainForecastID, charting.AnyToPoints(charting.F64ToAny(trainForecasts)))
	copyTrainChart.Datasets[GraphTrainForecastID].UpdateVariableLabel(0, fmt.Sprintf("Optimal Alpha: %.4f", bestAlpha))
	copyTrainChart.Datasets[GraphTrainForecastID].UpdateVariableLabel(1, fmt.Sprintf("Optimal Beta: %.4f", bestBeta))
	copyTrainChart.Datasets[GraphTrainForecastID].UpdateVariableLabel(2, fmt.Sprintf("Train MSE: %.4e", trainMSE))

	res = charting.NewRenderResponse()
	res.AddChart(copyTrainChart.ID, &copyTrainChart)
	res.CachePolicy = charting.CachePolicyDontCache
	return res
}

func RenderError(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadExchangeHistory(); err != nil {
		return res.NewError(err.Error())
	}

	// Holt parameters MUST be in [0, 1] to be stable
	minAlpha, maxAlpha := 0.0, 1.0
	minBeta, maxBeta := 0.0, 1.0

	step, ok := req.GetChartVariable(ChartHoltOptimalID, VariableParamStepID)
	if !ok {
		step = VariableHeatmapParamStep.Default
	}

	nAlpha := int((maxAlpha-minAlpha)/step) + 1
	nBeta := int((maxBeta-minBeta)/step) + 1

	heatmapData := make([]any, 0, nAlpha*nBeta)
	rawValues := make([]float64, nAlpha*nBeta)
	coords := make([]charting.DataPoint, nAlpha*nBeta)

	bestAlpha, bestBeta, bestMSE := math.MaxFloat64, math.MaxFloat64, math.MaxFloat64

	fmt.Printf("Calculating MSE grid for %s\n", OptimalChart.ID)
	for i := range nAlpha {
		alpha := minAlpha + step*float64(i)
		for j := range nBeta {
			beta := minBeta + step*float64(j)

			forecasts, _, _ := HoltForecast(trainExchangeRateData.ExchangeRate, alpha, beta)
			forecastMSE := analysis.MSE(trainExchangeRateData.ExchangeRate, forecasts)
			fmt.Printf("Alpha: %.2f, Beta: %.2f, MSE: %.4e\n", alpha, beta, forecastMSE)

			if forecastMSE < bestMSE {
				bestMSE = forecastMSE
				bestAlpha = alpha
				bestBeta = beta
			}

			index := i*nBeta + j
			rawValues[index] = forecastMSE
			coords[index] = charting.DataPoint{X: alpha, Y: &beta}
		}
	}

	// Outlier Suppression: Cap values at 10x bestMSE to keep the heatmap gradient useful
	capValue := bestMSE * 10
	for i, v := range rawValues {
		val := v
		if v > capValue {
			val = capValue
		}
		heatmapData = append(heatmapData, charting.HeatmapPoint{
			DataPoint: coords[i],
			Value:     &val,
		})
	}

	copyChart := charting.CopyChart(OptimalChart)
	copyChart.UpdateDataForDataset(GraphErrHeatmapID, heatmapData)

	copyChart.Datasets[GraphErrHeatmapID].UpdateVariableLabel(0, fmt.Sprintf("Optimal Alpha: %.4f", bestAlpha))
	copyChart.Datasets[GraphErrHeatmapID].UpdateVariableLabel(1, fmt.Sprintf("Optimal Beta: %.4f", bestBeta))
	copyChart.Datasets[GraphErrHeatmapID].UpdateVariableLabel(2, fmt.Sprintf("Optimal MSE: %.4e", bestMSE))

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	return res
}

func init() {
	TrainChart.RenderFunc = RenderHolt
	TestChart.RenderFunc = RenderHoltTest
	OptimalChart.RenderFunc = RenderError
}
