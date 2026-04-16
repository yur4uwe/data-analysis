package optimizations

import (
	"fmt"
	"labs/charting"
	"math"
	"slices"
	"strings"
)

const (
	OneDimChartID = "one-dim"
	TwoDimChartID = "two-dim"

	DataGraphID   = "orig-data"
	MinOptGraphID = "opt-path-min"
	MaxOptGraphID = "opt-path-max"

	VarStartID = "start"
	VarEndID   = "end"
	VarStepID  = "step"

	VarXStartID = "x-start"
	VarXEndID   = "x-end"
	VarXStepID  = "x-step"
	VarYStartID = "y-start"
	VarYEndID   = "y-end"
	VarYStepID  = "y-step"

	VarOptMethodID  = "opt-method"
	VarOptTolID     = "opt-tol"
	VarOptResultID  = "opt-result"
	VarOptX0ID      = "opt-x0"
	VarOptY0ID      = "opt-y0"
	VarOptLRID      = "opt-lr"
	VarOptSamplesID = "opt-samples"
)

var (
	VarOptMethodOneDim = charting.MutableField{
		ID:      VarOptMethodID,
		Label:   "Optimization Method",
		Default: 0, // index 0: None
		Control: charting.ControlSelect,
		Options: []string{
			"None",
			"Dichotomic Search",
			"Random Search",
			"Fast Descent",
		},
	}

	VarOptMethodTwoDim = charting.MutableField{
		ID:      VarOptMethodID,
		Label:   "Optimization Method",
		Default: 0, // index 0: None
		Control: charting.ControlSelect,
		Options: []string{
			"None",
			"Random Search",
			"Fast Descent",
		},
	}

	VarOptTol = charting.MutableField{
		ID:      VarOptTolID,
		Label:   "Tolerance",
		Default: 0.01,
		Min:     0.0001,
		Max:     1.0,
		Step:    0.001,
		Control: charting.ControlNumber,
	}

	VarOptResult = charting.MutableField{
		ID:      VarOptResultID,
		Label:   "Result",
		Control: charting.ControlNoControl,
	}

	VarOptX0 = charting.MutableField{
		ID:      VarOptX0ID,
		Label:   "Initial X",
		Default: 1.5,
		Min:     -10.0,
		Max:     10.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}

	VarOptY0 = charting.MutableField{
		ID:      VarOptY0ID,
		Label:   "Initial Y",
		Default: 0.1,
		Min:     -10.0,
		Max:     10.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}

	VarOptLR = charting.MutableField{
		ID:      VarOptLRID,
		Label:   "Learning Rate",
		Default: 0.1,
		Min:     0.001,
		Max:     1.0,
		Step:    0.01,
		Control: charting.ControlNumber,
	}

	VarOptSamples = charting.MutableField{
		ID:      VarOptSamplesID,
		Label:   "Samples",
		Default: 100,
		Min:     10,
		Max:     1000,
		Step:    10,
		Control: charting.ControlNumber,
	}

	VarStart = charting.MutableField{
		ID:      VarStartID,
		Label:   "Start (x > 0, x != 1)",
		Default: 1.1,
		Min:     0.1,
		Max:     100.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarEnd = charting.MutableField{
		ID:      VarEndID,
		Label:   "End",
		Default: 10.0,
		Min:     0.2,
		Max:     200.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarStep = charting.MutableField{
		ID:      VarStepID,
		Label:   "Step",
		Default: 0.1,
		Min:     0.01,
		Max:     5.0,
		Step:    0.01,
		Control: charting.ControlRange,
	}

	VarXStart = charting.MutableField{
		ID:      VarXStartID,
		Label:   "X Start",
		Default: -5.0,
		Min:     -10.0,
		Max:     10.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarXEnd = charting.MutableField{
		ID:      VarXEndID,
		Label:   "X End",
		Default: 5.0,
		Min:     0.2,
		Max:     20.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarXStep = charting.MutableField{
		ID:      VarXStepID,
		Label:   "X Step",
		Default: 0.2,
		Min:     0.05,
		Max:     2.0,
		Step:    0.05,
		Control: charting.ControlRange,
	}

	VarYStart = charting.MutableField{
		ID:      VarYStartID,
		Label:   "Y Start",
		Default: -5.0,
		Min:     -10.0,
		Max:     10.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarYEnd = charting.MutableField{
		ID:      VarYEndID,
		Label:   "Y End",
		Default: 5.0,
		Min:     -10.0,
		Max:     20.0,
		Step:    0.1,
		Control: charting.ControlNumber,
	}
	VarYStep = charting.MutableField{
		ID:      VarYStepID,
		Label:   "Y Step",
		Default: 0.2,
		Min:     0.05,
		Max:     2.0,
		Step:    0.05,
		Control: charting.ControlRange,
	}

	OneDimGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "f(x) = x + 1/ln(x)",
			BorderColor: charting.ColorBlue,
			BorderWidth: 0,
			GraphVariables: []charting.MutableField{
				VarOptMethodOneDim,
				VarOptX0,
				VarOptTol,
				VarOptLR,
				VarOptSamples,
				VarOptResult,
			},
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     3,
	}

	TwoDimGraph = charting.HeatmapDataset{
		BaseDataset: charting.BaseDataset{
			Label: "f(x, y)",
			GraphVariables: []charting.MutableField{
				VarOptMethodTwoDim,
				VarOptX0, VarOptY0,
				VarOptTol, VarOptLR, VarOptSamples,
				VarOptResult,
			},
		},
	}

	MinOptGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Minimization Path",
			BorderColor: charting.ColorRed,
			BorderWidth: 3,
		},
		BackgroundColor: charting.ColorRed,
		PointRadius:     10,
	}

	MaxOptGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Maximization Path",
			BorderColor: charting.ColorGreen,
			BorderWidth: 3,
		},
		BackgroundColor: charting.ColorGreen,
		PointRadius:     10,
	}

	OneDimChart = charting.Chart{
		ID:          OneDimChartID,
		Title:       "One-Dimensional Function",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "x",
		YAxisLabel:  "f(x)",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VarStart, VarEnd, VarStep,
		},
		Datasets: map[string]charting.Dataset{
			DataGraphID:   &OneDimGraph,
			MinOptGraphID: &MinOptGraph,
			MaxOptGraphID: &MaxOptGraph,
		},
	}

	TwoDimChart = charting.Chart{
		ID:          TwoDimChartID,
		Title:       "Two-Dimensional Function",
		Type:        charting.ChartTypeHeatmap,
		XAxisLabel:  "x",
		YAxisLabel:  "y",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VarXStart, VarXEnd, VarXStep,
			VarYStart, VarYEnd, VarYStep,
		},
		Datasets: map[string]charting.Dataset{
			DataGraphID:   &TwoDimGraph,
			MinOptGraphID: &MinOptGraph,
			MaxOptGraphID: &MaxOptGraph,
		},
	}
)

func onedimf(x float64) float64 {
	if x <= 0 || x == 1 {
		return math.NaN()
	}
	return x + (1 / math.Log(x))
}

func twodimf(x, y float64) float64 {
	c1 := 1.9
	c2 := 0.0
	term1 := (x - c1) * (x - c1)
	term2 := (y - c2) * (y - c2)
	inner := term1 - term2
	res := (inner * inner) + (c1 * c1) - (c2 * c2)
	if res <= 0 {
		return math.NaN()
	}
	return -math.Log(res)
}

// processOptimizationPath is a generalized pseudo n-dim path processor.
func processOptimizationPath(
	path [][]float64,
	f func(...float64) float64,
) ([]float64, []float64, string, bool) {
	if len(path) == 0 {
		return nil, nil, "Failed", false
	}

	dims := len(path[0])
	px := make([]float64, 0, len(path))
	py := make([]float64, 0, len(path))

	for _, p := range path {
		if slices.ContainsFunc(p, isInvalid) {
			continue
		}
		val := f(p...)
		if isInvalid(val) {
			continue
		}

		px = append(px, p[0])
		if dims == 1 {
			py = append(py, val)
		} else {
			py = append(py, p[1])
		}
	}

	if len(px) > 0 {
		last := path[len(path)-1]
		lastVal := f(last...)

		var sb strings.Builder
		sb.WriteString("f(")
		for i, coord := range last {
			if i > 0 {
				sb.WriteString(", ")
			}
			fmt.Fprintf(&sb, "%.4f", coord)
		}
		fmt.Fprintf(&sb, ") = %.4f, path length = %d", lastVal, len(px))
		return px, py, sb.String(), true
	}

	return nil, nil, "Failed", false
}

func RenderOneDim(req *charting.RenderRequest) (res *charting.RenderResponse) {

	start, _ := req.GetChartVariable(OneDimChartID, VarStartID)
	end, _ := req.GetChartVariable(OneDimChartID, VarEndID)
	step, _ := req.GetChartVariable(OneDimChartID, VarStepID)

	methodIdx, _ := req.GetGraphVariable(OneDimChartID, DataGraphID, VarOptMethodID)
	tol, _ := req.GetGraphVariable(OneDimChartID, DataGraphID, VarOptTolID)
	lr, _ := req.GetGraphVariable(OneDimChartID, DataGraphID, VarOptLRID)
	samples, _ := req.GetGraphVariable(OneDimChartID, DataGraphID, VarOptSamplesID)
	x0, _ := req.GetGraphVariable(OneDimChartID, DataGraphID, VarOptX0ID)

	if start <= 0 || (start <= 1 && end >= 1) {
		return res.NewError("invalid start value: x must be > 0 and != 1")
	}

	chartCopy := charting.CopyChart(OneDimChart)

	n := int((end-start)/step) + 1
	x := make([]float64, 0, n)
	y := make([]float64, 0, n)
	for i := range n {
		xVal := start + float64(i)*step
		val := onedimf(xVal)
		if !isInvalid(val) {
			x = append(x, xVal)
			y = append(y, val)
		}
	}
	chartCopy.UpdatePointsForDataset(DataGraphID, x, y)

	f1d := func(args ...float64) float64 { return onedimf(args[0]) }
	f1dNeg := func(args ...float64) float64 { return -onedimf(args[0]) }
	fNeg := func(arg float64) float64 { return -onedimf(arg) }

	if methodIdx != 0 {
		var pathMin, pathMax [][]float64

		switch methodIdx {
		case 1: // Dichotomic Search
			pMin := dichotomicSearch(onedimf, start, end, tol)
			pathMin = make([][]float64, len(pMin))
			for i, v := range pMin {
				pathMin[i] = []float64{v}
			}
			pMax := dichotomicSearch(fNeg, start, end, tol)
			pathMax = make([][]float64, len(pMax))
			for i, v := range pMax {
				pathMax[i] = []float64{v}
			}
		case 2: // Random Search
			chartCopy.Datasets[MaxOptGraphID].(*charting.GridDataset).HideLine = true
			chartCopy.Datasets[MinOptGraphID].(*charting.GridDataset).HideLine = true
			bounds := [][]float64{{start, end}}
			pathMin, pathMax = randomSearchNdim(f1d, int(samples), bounds)
		case 3: // Fast Descent
			bounds := [][]float64{{start, end}}
			pathMin = fastDescentNdim(f1d, []float64{x0}, bounds, 0.01, tol, lr)
			pathMax = fastDescentNdim(f1dNeg, []float64{x0}, bounds, 0.01, tol, lr)
		}

		px, py, resMin, okMin := processOptimizationPath(pathMin, f1d)
		chartCopy.UpdatePointsForDataset(MinOptGraphID, px, py)
		px, py, resMax, okMax := processOptimizationPath(pathMax, f1d)
		chartCopy.UpdatePointsForDataset(MaxOptGraphID, px, py)

		var resMsg string
		if okMin || okMax {
			resMsg = fmt.Sprintf("Min: %s | Max: %s", resMin, resMax)
		} else {
			resMsg = "Result: Failed (NaN)"
		}
		chartCopy.Datasets[DataGraphID].UpdateVariableLabel(VarOptResultID, resMsg)
	} else {
		chartCopy.UpdatePointsForDataset(MinOptGraphID, nil, nil)
		chartCopy.UpdatePointsForDataset(MaxOptGraphID, nil, nil)
		chartCopy.Datasets[DataGraphID].UpdateVariableLabel(VarOptResultID, "Result: None")
	}

	res = charting.NewRenderResponse()
	chartCopy.GenerateLabels(2)
	res.AddChart(OneDimChartID, &chartCopy)
	return res
}

func RenderTwoDim(req *charting.RenderRequest) (res *charting.RenderResponse) {
	xStart, _ := req.GetChartVariable(TwoDimChartID, VarXStartID)
	xEnd, _ := req.GetChartVariable(TwoDimChartID, VarXEndID)
	xStep, _ := req.GetChartVariable(TwoDimChartID, VarXStepID)
	yStart, _ := req.GetChartVariable(TwoDimChartID, VarYStartID)
	yEnd, _ := req.GetChartVariable(TwoDimChartID, VarYEndID)
	yStep, _ := req.GetChartVariable(TwoDimChartID, VarYStepID)

	methodIdx, _ := req.GetGraphVariable(TwoDimChartID, DataGraphID, VarOptMethodID)
	tol, _ := req.GetGraphVariable(TwoDimChartID, DataGraphID, VarOptTolID)
	lr, _ := req.GetGraphVariable(TwoDimChartID, DataGraphID, VarOptLRID)
	samples, _ := req.GetGraphVariable(TwoDimChartID, DataGraphID, VarOptSamplesID)
	x0, _ := req.GetGraphVariable(TwoDimChartID, DataGraphID, VarOptX0ID)
	y0, _ := req.GetGraphVariable(TwoDimChartID, DataGraphID, VarOptY0ID)

	chartCopy := charting.CopyChart(TwoDimChart)

	nx := int((xEnd-xStart)/xStep) + 1
	ny := int((yEnd-yStart)/yStep) + 1
	points := make([]any, 0, nx*ny)
	for i := range nx {
		for j := range ny {
			xVal := xStart + float64(i)*xStep
			yVal := yStart + float64(j)*yStep
			zVal := twodimf(xVal, yVal)
			if !isInvalid(zVal) {
				yC, zC := yVal, zVal
				points = append(points, charting.HeatmapPoint{
					DataPoint: charting.DataPoint{X: xVal, Y: &yC},
					Value:     &zC,
				})
			}
		}
	}
	chartCopy.UpdateDataForDataset(DataGraphID, points)

	f2d := func(args ...float64) float64 { return twodimf(args[0], args[1]) }
	f2dNeg := func(args ...float64) float64 { return -twodimf(args[0], args[1]) }

	if methodIdx != 0 {
		var pathMin, pathMax [][]float64

		switch methodIdx {
		case 1: // Random Search
			chartCopy.Datasets[MaxOptGraphID].(*charting.GridDataset).HideLine = true
			chartCopy.Datasets[MinOptGraphID].(*charting.GridDataset).HideLine = true
			bounds := [][]float64{{xStart, xEnd}, {yStart, yEnd}}
			pathMin, pathMax = randomSearchNdim(f2d, int(samples), bounds)
		case 2: // Fast Descent
			bounds := [][]float64{{xStart, xEnd}, {yStart, yEnd}}
			pathMin = fastDescentNdim(f2d, []float64{x0, y0}, bounds, 0.01, tol, lr)
			pathMax = fastDescentNdim(f2dNeg, []float64{x0, y0}, bounds, 0.01, tol, lr)
		}

		px, py, resMin, okMin := processOptimizationPath(pathMin, f2d)
		chartCopy.UpdatePointsForDataset(MinOptGraphID, px, py)
		px, py, resMax, okMax := processOptimizationPath(pathMax, f2d)
		chartCopy.UpdatePointsForDataset(MaxOptGraphID, px, py)

		var resMsg string
		if okMin || okMax {
			resMsg = fmt.Sprintf("Min: %s | Max: %s", resMin, resMax)
		} else {
			resMsg = "Result: Failed (NaN)"
		}
		chartCopy.Datasets[DataGraphID].UpdateVariableLabel(VarOptResultID, resMsg)
	} else {
		chartCopy.UpdatePointsForDataset(MinOptGraphID, nil, nil)
		chartCopy.UpdatePointsForDataset(MaxOptGraphID, nil, nil)
		chartCopy.Datasets[DataGraphID].UpdateVariableLabel(VarOptResultID, "Result: None")
	}

	res = charting.NewRenderResponse()
	chartCopy.GenerateLabels(2)
	res.AddChart(TwoDimChartID, &chartCopy)
	return res
}
