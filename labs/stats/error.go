package stats

import (
	"labs/charting"
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
	MeanSampleField = charting.MutableField{
		ID:      VariableSampleMeanID,
		Label:   "Mean for the sample",
		Default: 0,
		Min:     -100,
		Max:     100,
		Step:    1,
		Control: charting.ControlNumber,
	}

	StdDevSampleField = charting.MutableField{
		ID:      VariableSampleStdDevID,
		Label:   "Standart Deviation for the sample",
		Default: 10,
		Min:     0,
		Max:     100,
		Step:    1,
		Control: charting.ControlRange,
	}

	MeanCorrelationGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Relatioship between error of mean and size of sample",
			BorderColor: charting.ColorLightPurple,
			BorderWidth: 2,
			Togglable:   true,
			GraphVariables: []charting.MutableField{
				MeanSampleField,
			},
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}

	StdDevCorrelationGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Relationship between error of stddev and size of sample",
			BorderColor: charting.ColorOrange,
			BorderWidth: 2,
			Togglable:   true,
			GraphVariables: []charting.MutableField{
				StdDevSampleField,
			},
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}

	MaxSampleSizeField = charting.MutableField{
		ID:      VariableMaxSampleSizeID,
		Label:   "Max size to calculate error for",
		Default: 50_000,
		Min:     100,
		Max:     100_000,
		Step:    100,
		Control: charting.ControlNumber,
	}
	StepSampleSizeField = charting.MutableField{
		ID:      VariableSampleSizeStepID,
		Label:   "Step to choose which sizes to calculate error for",
		Default: 10,
		Min:     0,
		Max:     100_000,
		Step:    1,
		Control: charting.ControlNumber,
	}

	CorrelationChart = charting.Chart{
		ID:          CorrelationChartID,
		Title:       "Relationship between sample size and error of its parameters",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Sample Size",
		XAxisConfig: charting.LogarithmicAxis,
		YAxisLabel:  "Error",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			MeanCorellationGraphID:   &MeanCorrelationGraph,
			StdDevCorrelationGraphID: &StdDevCorrelationGraph,
		},
		ChartVariables: []charting.MutableField{
			MaxSampleSizeField,
		},
	}
)

func RenderError(req *charting.RenderRequest) (res *charting.RenderResponse) {
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

	sample := GenerateNormalSamples(theoretical_mean, theoretical_stddev, int(max_size))
	for i := 1; i <= 100; i++ {
		sliced_sample := sample[:i*n]

		actual_mean := CalculateMean(sliced_sample)
		actual_stddev := CalculateStdDev(sliced_sample, actual_mean)

		x = append(x, float64(i*n))
		mean_errors = append(mean_errors, (actual_mean - theoretical_mean))
		stddev_errors = append(stddev_errors, (actual_stddev - theoretical_stddev))
	}

	copyChart := charting.CopyChart(CorrelationChart)

	copyChart.UpdatePointsForDataset(MeanCorellationGraphID, x, mean_errors)
	copyChart.UpdatePointsForDataset(StdDevCorrelationGraphID, x, stddev_errors)

	res = charting.NewRenderResponse()
	res.AddChart(CorrelationChartID, &copyChart)
	return res
}
