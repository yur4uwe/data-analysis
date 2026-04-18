package visualization

import (
	"labs/charting"
	"math"
)

const (
	FunctionChartID = "function"

	FunctionGraphID = "orig-data" // due to the default, baseline graph on frontend, one chart should be called orig-data

	VariableStartID = "start"
	VariableEndID   = "end"
	VariableStepID  = "step"
)

var (
	VariableStart = charting.MutableField{
		ID:      VariableStartID,
		Label:   "Interval Start",
		Default: -5.0,
		Min:     -100.0,
		Max:     100.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VariableEnd = charting.MutableField{
		ID:      VariableEndID,
		Label:   "Interval End",
		Default: 5.0,
		Min:     -100.0,
		Max:     100.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VariableStep = charting.MutableField{
		ID:      VariableStepID,
		Label:   "Step Size",
		Default: 0.1,
		Min:     0.1,
		Max:     10.0,
		Step:    0.1,
		Control: charting.ControlRange,
	}

	FunctionChart = charting.Chart{
		ID:          FunctionChartID,
		Title:       "Function Representation",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VariableStart,
			VariableEnd,
			VariableStep,
		},
		Datasets: map[string]charting.Dataset{
			FunctionGraphID: &FunctionGraph,
		},
	}

	FunctionGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Function cos(x)*e^(-(|x|)) Graph",
			BorderColor: charting.ColorIndigo,
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     0,
	}
)

func f(x float64) float64 {
	return math.Cos(x) * math.Exp(-math.Abs(x))
}

func RenderFunction(req *charting.RenderRequest) (res *charting.RenderResponse) {
	start, ok := req.GetChartVariable(FunctionChartID, VariableStartID)
	if !ok {
		start = VariableStart.Default
	}
	end, ok := req.GetChartVariable(FunctionChartID, VariableEndID)
	if !ok {
		end = VariableEnd.Default
	}
	step, ok := req.GetChartVariable(FunctionChartID, VariableStepID)
	if !ok {
		step = VariableStep.Default
	}

	if start >= end {
		return res.NewErrorf(
			"invalid interval: start (%v) must be less than end (%v)",
			start, end,
		)
	}
	if step <= 0 {
		return res.NewErrorf(
			"invalid step size: step (%v) must be positive",
			step,
		)
	}
	if step > (end - start) {
		return res.NewErrorf(
			"invalid step size: step (%v) must be less than or equal to the interval length (%v)",
			step, end-start,
		)
	}

	n := int((end-start)/step) + 1

	x := make([]float64, 0, n)
	y := make([]float64, 0, n)

	for i := range n {
		xVal := start + float64(i)*step
		x = append(x, xVal)
		y = append(y, f(xVal))
	}

	chartCopy := charting.CopyChart(FunctionChart)
	chartCopy.UpdatePointsForDataset(FunctionGraphID, x, y)
	chartCopy.GenerateLabels(2)

	res = charting.NewRenderResponse()
	res.AddChart(FunctionChartID, &chartCopy)

	return res
}
