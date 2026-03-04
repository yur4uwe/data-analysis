package stats

import (
	"labs/labs/common"
	"math"
)

const (
	CorrelationChartID = "errors"

	VariableMaxSampleSizeID  = "max-sample-size"
	VariableSampleSizeStepID = "step-sample-size"

	MeanCorellationGraphID   = "mean-err"
	StdDevCorrelationGraphID = "stddev-err"

	VariableSampleMeanID   = "sample-mean"
	VariableSampleStdDevID = "sample-stddev"
)

var (
	MeanSampleField = common.MutableField{
		ID:      VariableSampleMeanID,
		Label:   "Mean for the sample",
		Default: 0,
		Min:     -100,
		Max:     100,
		Step:    1,
		Control: common.ControlNumber,
	}

	StdDevSampleField = common.MutableField{
		ID:      VariableSampleStdDevID,
		Label:   "Standart Deviation for the sample",
		Default: 10,
		Min:     0,
		Max:     100,
		Step:    1,
		Control: common.ControlRange,
	}

	MeanCorrelationGraph = common.ChartDataset{
		Label:           "Relatioship between error of mean and size of sample",
		BorderColor:     common.Color1,
		BackgroundColor: []string{common.ColorTransparent},
		ShowLine:        true,
		Togglable:       true,
		PointRadius:     0,
		GraphVariables: []common.MutableField{
			MeanSampleField,
		},
	}

	StdDevCorrelationGraph = common.ChartDataset{
		Label:           "Relationship between error of stddev and size of sample",
		BorderColor:     common.Color8,
		BackgroundColor: []string{common.ColorTransparent},
		ShowLine:        true,
		Togglable:       true,
		PointRadius:     0,
		GraphVariables: []common.MutableField{
			StdDevSampleField,
		},
	}

	MaxSampleSizeField = common.MutableField{
		ID:      VariableMaxSampleSizeID,
		Label:   "Max size to calculate error for",
		Default: 50_000,
		Min:     100,
		Max:     100_000,
		Step:    100,
		Control: common.ControlNumber,
	}
	StepSampleSizeField = common.MutableField{
		ID:      VariableSampleSizeStepID,
		Label:   "Step to choose which sizes to calculate error for",
		Default: 10,
		Min:     0,
		Max:     100_000,
		Step:    1,
		Control: common.ControlNumber,
	}

	CorrelationChart = common.Chart{
		ID:          CorrelationChartID,
		Title:       "Relationship between sample size and error of its parameters",
		Type:        common.ChartTypeLine,
		XAxisLabel:  "Sample Size",
		XAxisConfig: common.LogarithmicAxis,
		YAxisLabel:  "Error",
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			MeanCorellationGraphID:   &MeanCorrelationGraph,
			StdDevCorrelationGraphID: &StdDevCorrelationGraph,
		},
		ChartVariables: []common.MutableField{
			MaxSampleSizeField,
		},
		// RenderFunc set in init() to avoid initialization cycle
	}
)

func RenderError(req *common.RenderRequest) (res *common.RenderResponse) {
	max_size, ok := req.GetChartVariable(CorrelationChartID, VariableMaxSampleSizeID)
	if !ok {
		max_size = MaxSampleSizeField.Default
	}

	theoretical_mean, ok := req.GetGraphVariable(CorrelationChartID, MeanCorellationGraphID, VariableSampleMeanID)
	if !ok {
		theoretical_mean = MeanSampleField.Default
	}
	theoretical_stddev, ok := req.GetGraphVariable(CorrelationChartID, StdDevCorrelationGraphID, VariableSampleStdDevID)
	if !ok {
		theoretical_stddev = StdDevSampleField.Default
	}

	n := int(max_size) / 100
	x := make([]float64, 0, n)
	mean_errors := make([]float64, 0, n)
	stddev_errors := make([]float64, 0, n)

	for i := 1; i <= 100; i++ {
		sample := GenerateNormalSamples(theoretical_mean, theoretical_stddev, int(i))

		actual_mean := CalculateMean(sample)
		actual_stddev := CalculateStdDev(sample, actual_mean)

		x = append(x, float64(i*n))
		mean_errors = append(mean_errors, math.Abs(actual_mean-theoretical_mean))
		stddev_errors = append(stddev_errors, math.Abs(actual_stddev-theoretical_stddev))
	}

	copyChart := common.CopyChart(CorrelationChart)

	copyChart.UpdatePointsForDataset(MeanCorellationGraphID, x, mean_errors)
	copyChart.UpdatePointsForDataset(StdDevCorrelationGraphID, x, stddev_errors)

	res = common.NewRenderResponse()
	res.AddChart(CorrelationChartID, &copyChart)
	return res
}
