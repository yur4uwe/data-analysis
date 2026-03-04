package visualization

import (
	"fmt"
	"labs/labs/common"
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
	VariableStart = common.MutableField{
		ID:      VariableStartID,
		Label:   "Interval Start",
		Default: -5.0,
		Min:     -100.0,
		Max:     100.0,
		Step:    0.1,
		Control: common.ControlNumber,
	}
	VariableEnd = common.MutableField{
		ID:      VariableEndID,
		Label:   "Interval End",
		Default: 5.0,
		Min:     -100.0,
		Max:     100.0,
		Step:    0.1,
		Control: common.ControlNumber,
	}
	VariableStep = common.MutableField{
		ID:      VariableStepID,
		Label:   "Step Size",
		Default: 0.1,
		Min:     0.1,
		Max:     10.0,
		Step:    0.1,
		Control: common.ControlRange,
	}

	FunctionChart = common.Chart{
		ID:          FunctionChartID,
		Title:       "Function Representation",
		Type:        common.ChartTypeLine,
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		ChartVariables: []common.MutableField{
			VariableStart,
			VariableEnd,
			VariableStep,
		},
		Datasets: map[string]*common.ChartDataset{
			FunctionGraphID: &FunctionGraph,
		},
	}

	FunctionGraph = common.ChartDataset{
		Label:           "Function cos(x)*e^(-(|x|)) Graph",
		BorderColor:     common.Color1,
		BackgroundColor: []string{common.ColorTransparent},
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
	}

	FunctionMeta = FunctionChart.Meta()
)

func f(x float64) float64 {
	return math.Cos(x) * math.Exp(-math.Abs(x))
}

func RenderFunction(req *common.RenderRequest) (res *common.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)
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

	chartCopy := common.CopyChart(FunctionChart)
	chartCopy.UpdatePointsForDataset(FunctionGraphID, x, y)

	res = common.NewRenderResponse()
	res.AddChart(FunctionChartID, &chartCopy)

	return res
}
