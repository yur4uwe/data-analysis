package forecastinglinparab

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	LabID = "9"

	ChartTrainDataID              = "train-data"
	ChartTestDataID               = "test-data"
	ChartOptimalParabolicParamsID = "optimal-parabolic-params"
	ChartOptimalLinearParamsID    = "optimal-linear-params"

	GraphOriginalDataID    = "original-data"
	GraphLinearApproxID    = "linear-approx"
	GraphParabolicApproxID = "parabolic-approx"

	VariableParabolicFitCoefficientsID = "parabolic-fit-coefficients"
	VariableLinearFitCoefficientsID    = "linear-fit-coefficients"
)

var (
	bestA = 0.0
	bestB = 0.0
	bestC = 0.0

	hasTrained = false

	exchangeRateData = &ExchageRateHistory{}
)

type ExchageRateHistory struct {
	Date         []string  `csv:"Дата"`
	ExchangeRate []float64 `csv:"Офіційний курс гривні"`
}

var (
	LinParabConfig = charting.NewLabConfig(
		LabID,
		"Linear and Parabolic Approximation",
		map[string]*charting.Chart{
			ChartTrainDataID:           &TrainDataChart,
			ChartTestDataID:            &TestDataChart,
			ChartOptimalLinearParamsID: &OptimalLinearParamsChart,
			// ChartOptimalParabolicParamsID: &OptimalParabolicParamsChart,
		},
	)

	TrainDataChart = charting.Chart{
		ID:          ChartTrainDataID,
		Type:        charting.ChartTypeLine,
		Title:       "Train Data",
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Amount",
		YAxisConfig: charting.LinearAxis,
	}

	TestDataChart = charting.Chart{
		ID:          ChartTestDataID,
		Type:        charting.ChartTypeLine,
		Title:       "Test Data",
		XAxisLabel:  "Date",
		XAxisConfig: charting.CategoryAxis,
		YAxisLabel:  "Amount",
		YAxisConfig: charting.LinearAxis,
	}

	// We'll leave it for now, i don't know how
	// tp make a chart of the optimal parabolic params
	OptimalParabolicParamsChart = charting.Chart{}

	OptimalLinearParamsChart = charting.Chart{
		ID:          ChartOptimalLinearParamsID,
		Type:        charting.ChartTypeHeatmap,
		Title:       "Optimal Linear Params (A + Bt)",
		XAxisLabel:  "A",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "B",
		YAxisConfig: charting.LinearAxis,
	}

	OriginalDataGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Original Data",
			BorderColor: charting.ColorEmerald,
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     2,
	}

	LinearFitCoefficients = charting.MutableField{
		ID:      VariableLinearFitCoefficientsID,
		Label:   "Linear Fit Coefficients",
		Control: charting.ControlNoControl,
	}

	ParabolicFitCoefficients = charting.MutableField{
		ID:      VariableParabolicFitCoefficientsID,
		Label:   "Parabolic Fit Coefficients",
		Control: charting.ControlNoControl,
	}

	LinearApproxGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Linear Approximation",
			BorderColor:    charting.ToColor("#16a34a"),
			BorderWidth:    2,
			Togglable:      true,
			GraphVariables: []charting.MutableField{LinearFitCoefficients},
		},
		BackgroundColor: charting.ToColor("rgba(22, 163, 74, 0.1)"),
		PointRadius:     0,
	}

	ParabolicApproxGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:          "Parabolic Approximation",
			BorderColor:    charting.ToColor("#9333ea"),
			BorderWidth:    2,
			Togglable:      true,
			GraphVariables: []charting.MutableField{ParabolicFitCoefficients},
		},
		BackgroundColor: charting.ToColor("rgba(147, 51, 234, 0.1)"),
		PointRadius:     0,
	}
)

func loadExchageHistory() error {
	if len(exchangeRateData.ExchangeRate) > 0 {
		return nil
	}
	f, err := os.Open("./data/lab_9_var_12.csv")
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.Comma = ';'
	exchangeRateData = &ExchageRateHistory{}
	if err := d.Decode(exchangeRateData); err != nil {
		return fmt.Errorf("error decoding csv: %w", err)
	}

	return nil
}

func init() {
	TrainDataChart.RenderFunc = RenderTrain
	TestDataChart.RenderFunc = RenderTest
}
