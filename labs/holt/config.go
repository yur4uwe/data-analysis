package holt

import (
	"labs/charting"
)

const (
	LabID = "8"

	ChartHoltTrainID = "holt-train"
	ChartHoltTestID  = "holt-test"

	GraphTrainActualID   = "train-actual"
	GraphTrainForecastID = "train-forecast"

	GraphTestActualID   = "test-actual"
	GraphTestForecastID = "test-forecast"

	VariableEpochsID       = "epochs"
	VariableLearningRateID = "learning-rate"

	DisplayOptimalAlphaID = "optimal-alpha"
	DisplayOptimalBetaID  = "optimal-beta"
	DisplayTrainMSEID     = "train-mse"
	DisplayTestMSEID      = "test-mse"
)

type ExchangeRateHistory struct {
	Date         []string  `csv:"Дата"`
	ExchangeRate []float64 `csv:"Офіційний курс гривні"`
}

var (
	VariableEpochs = charting.MutableField{
		ID:      VariableEpochsID,
		Label:   "Gradient Descent Epochs",
		Default: 1000,
		Min:     10,
		Max:     10000,
		Step:    10,
		Control: charting.ControlRange,
	}

	VariableLearningRate = charting.MutableField{
		ID:      VariableLearningRateID,
		Label:   "Learning Rate",
		Default: 0.01,
		Min:     0.001,
		Max:     1.0,
		Step:    0.001,
		Control: charting.ControlRange,
	}

	OptimalAlphaField = charting.MutableField{ID: DisplayOptimalAlphaID, Label: "Optimal Alpha: -", Control: charting.ControlNoControl}
	OptimalBetaField  = charting.MutableField{ID: DisplayOptimalBetaID, Label: "Optimal Beta: -", Control: charting.ControlNoControl}
	TrainMSEField     = charting.MutableField{ID: DisplayTrainMSEID, Label: "Train MSE: -", Control: charting.ControlNoControl}
	TestMSEField      = charting.MutableField{ID: DisplayTestMSEID, Label: "Test MSE: -", Control: charting.ControlNoControl}

	TrainActualGraph = charting.ChartDataset{
		Label:           "Train Data",
		BorderColor:     charting.ColorTeal,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       false,
	}

	TrainForecastGraph = charting.ChartDataset{
		Label:           "Holt Forecast (Train)",
		BorderColor:     charting.ColorAmber,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []charting.MutableField{
			OptimalAlphaField,
			OptimalBetaField,
			TrainMSEField,
		},
	}

	TestActualGraph = charting.ChartDataset{
		Label:           "Test Data",
		BorderColor:     charting.ColorTeal,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       false,
	}

	TestForecastGraph = charting.ChartDataset{
		Label:           "Holt Forecast (Test)",
		BorderColor:     charting.ColorRed,
		BackgroundColor: []string{charting.ColorTransparent},
		BorderWidth:     2,
		PointRadius:     0,
		ShowLine:        true,
		Togglable:       true,
		GraphVariables: []charting.MutableField{
			TestMSEField,
		},
	}

	TrainChart = charting.Chart{
		ID:          ChartHoltTrainID,
		Title:       "Holt's Method - Training Phase",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			GraphTrainActualID:   &TrainActualGraph,
			GraphTrainForecastID: &TrainForecastGraph,
		},
		ChartVariables: []charting.MutableField{
			VariableEpochs,
			VariableLearningRate,
		},
	}

	TestChart = charting.Chart{
		ID:          ChartHoltTestID,
		Title:       "Holt's Method - Testing Phase",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			GraphTestActualID:   &TestActualGraph,
			GraphTestForecastID: &TestForecastGraph,
		},
	}

	Config = charting.NewLabConfig(
		LabID,
		"Holt's Linear Trend Forecasting",
		map[string]*charting.Chart{
			ChartHoltTestID:  &TestChart,
			ChartHoltTrainID: &TrainChart,
		},
	)

	Metadata = Config.Lab

	testExchangeRateData  = &ExchangeRateHistory{}
	trainExchangeRateData = &ExchangeRateHistory{}
)
