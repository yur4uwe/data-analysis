package holt

import (
	"labs/charting"
)

const (
	LabID = "8"

	ChartHoltTrainID   = "holt-train"
	ChartHoltTestID    = "holt-test"
	ChartHoltOptimalID = "holt-optimal"

	GraphTrainActualID   = "train-actual"
	GraphTrainForecastID = "train-forecast"
	GraphErrHeatmapID    = "error-heatmap"

	GraphTestActualID   = "test-actual"
	GraphTestForecastID = "test-forecast"

	VariableEpochsID       = "epochs"
	VariableLearningRateID = "learning-rate"
	VariableParamStepID    = "param-step"

	DisplayOptimalAlphaID = "optimal-alpha"
	DisplayOptimalBetaID  = "optimal-beta"
	DisplayTrainMSEID     = "train-mse"
	DisplayTestMSEID      = "test-mse"
	DisplayOptimalMSEID   = "optimal-mse"
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
		Default: 10.0,
		Min:     0.01,
		Max:     100.0,
		Step:    0.01,
		Control: charting.ControlRange,
	}

	VariableHeatmapParamStep = charting.MutableField{
		ID:      VariableParamStepID,
		Label:   "Parameter Step Size",
		Default: 0.05,
		Min:     0.001,
		Max:     0.1,
		Step:    0.001,
		Control: charting.ControlRange,
	}

	OptimalAlphaField = charting.MutableField{ID: DisplayOptimalAlphaID, Label: "Optimal Alpha: -", Control: charting.ControlNoControl}
	OptimalBetaField  = charting.MutableField{ID: DisplayOptimalBetaID, Label: "Optimal Beta: -", Control: charting.ControlNoControl}
	TrainMSEField     = charting.MutableField{ID: DisplayTrainMSEID, Label: "Train MSE: -", Control: charting.ControlNoControl}
	TestMSEField      = charting.MutableField{ID: DisplayTestMSEID, Label: "Test MSE: -", Control: charting.ControlNoControl}
	OptimalMSEField   = charting.MutableField{ID: DisplayOptimalMSEID, Label: "Optimal MSE: -", Control: charting.ControlNoControl}

	TrainActualGraph = charting.CategoricalDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Train Data",
			BorderColor: charting.ColorTeal,
			BorderWidth: 2,
			Togglable:   false,
		},
		BackgroundColor: []charting.Color{charting.ColorTransparent},
	}

	TrainForecastGraph = charting.CategoricalDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Holt Forecast (Train)",
			BorderColor: charting.ColorAmber,
			BorderWidth: 2,
			Togglable:   true,
			GraphVariables: []charting.MutableField{
				OptimalAlphaField,
				OptimalBetaField,
				TrainMSEField,
			},
		},
		BackgroundColor: []charting.Color{charting.ColorTransparent},
	}

	TestActualGraph = charting.CategoricalDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Test Data",
			BorderColor: charting.ColorTeal,
			BorderWidth: 2,
			Togglable:   false,
		},
		BackgroundColor: []charting.Color{charting.ColorTransparent},
	}

	TestForecastGraph = charting.CategoricalDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Holt Forecast (Test)",
			BorderColor: charting.ColorRed,
			BorderWidth: 2,
			Togglable:   true,
			GraphVariables: []charting.MutableField{
				TestMSEField,
				OptimalAlphaField,
				OptimalBetaField,
			},
		},
		BackgroundColor: []charting.Color{charting.ColorTransparent},
	}

	HeatmapGraph = charting.HeatmapDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Holt error vs alpha and beta",
			BorderColor: charting.ColorTransparent,
			BorderWidth: 0,
			GraphVariables: []charting.MutableField{
				OptimalMSEField,
				OptimalAlphaField,
				OptimalBetaField,
			},
		},
		BackgroundColor: []charting.Color{charting.ColorBlue, charting.ColorRed},
	}

	TrainChart = charting.Chart{
		ID:          ChartHoltTrainID,
		Title:       "Holt's Method - Training Phase",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Rate (UAH)",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
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
		Datasets: map[string]charting.Dataset{
			GraphTestActualID:   &TestActualGraph,
			GraphTestForecastID: &TestForecastGraph,
		},
	}

	OptimalChart = charting.Chart{
		ID:          ChartHoltOptimalID,
		Title:       "Heatmap of errors vs alpha and beta",
		Type:        charting.ChartTypeHeatmap,
		XAxisLabel:  "Alpha",
		YAxisLabel:  "Beta",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VariableHeatmapParamStep,
		},
		Datasets: map[string]charting.Dataset{
			GraphErrHeatmapID: &HeatmapGraph,
		},
	}

	Config = charting.NewLabConfig(
		LabID,
		"Holt's Linear Trend Forecasting",
		map[string]*charting.Chart{
			ChartHoltTestID:    &TestChart,
			ChartHoltTrainID:   &TrainChart,
			ChartHoltOptimalID: &OptimalChart,
		},
	)

	Metadata = Config.Lab

	testExchangeRateData  = &ExchangeRateHistory{}
	trainExchangeRateData = &ExchangeRateHistory{}
)
