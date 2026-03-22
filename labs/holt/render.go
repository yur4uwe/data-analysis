package holt

import (
	"errors"
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"os"
	"strings"
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

	testMSE := MSE(testData, testForecasts)

	copyTestChart := charting.CopyChart(TestChart)
	copyTestChart.Labels = testDates
	copyTestChart.UpdateDataForDataset(GraphTestActualID, charting.ToAnySlice(testData))
	copyTestChart.UpdateDataForDataset(GraphTestForecastID, charting.ToAnySlice(testForecasts))

	for i := range copyTestChart.Datasets[GraphTestForecastID].GraphVariables {
		field := &copyTestChart.Datasets[GraphTestForecastID].GraphVariables[i]
		if strings.HasSuffix(field.ID, DisplayTestMSEID) {
			field.Label = fmt.Sprintf("Test MSE: %.4f", testMSE)
		}
	}

	res = charting.NewRenderResponse()
	res.AddChart(copyTestChart.ID, &copyTestChart)
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
	trainMSE := MSE(trainData, trainForecasts)

	// Save for Test phase
	testL = finalL
	testT = finalT
	bestA = bestAlpha
	bestB = bestBeta
	hasTrained = true

	copyTrainChart := charting.CopyChart(TrainChart)
	copyTrainChart.Labels = trainDates
	copyTrainChart.UpdateDataForDataset(GraphTrainActualID, charting.ToAnySlice(trainData))
	copyTrainChart.UpdateDataForDataset(GraphTrainForecastID, charting.ToAnySlice(trainForecasts))

	for i := range copyTrainChart.Datasets[GraphTrainForecastID].GraphVariables {
		field := &copyTrainChart.Datasets[GraphTrainForecastID].GraphVariables[i]
		switch field.ID {
		case DisplayOptimalAlphaID:
			field.Label = fmt.Sprintf("Optimal Alpha: %.4f", bestAlpha)
		case DisplayOptimalBetaID:
			field.Label = fmt.Sprintf("Optimal Beta: %.4f", bestBeta)
		case DisplayTrainMSEID:
			field.Label = fmt.Sprintf("Train MSE: %.4f", trainMSE)
		}
	}

	res = charting.NewRenderResponse()
	res.AddChart(copyTrainChart.ID, &copyTrainChart)
	return res
}

func init() {
	TrainChart.RenderFunc = RenderHolt
	TestChart.RenderFunc = RenderHoltTest
}
