package optimizations

import (
	"fmt"
	"labs/charting"
	"math"
)

const (
	OneDimChartID = "one-dim"
	TwoDimChartID = "two-dim"

	OneDimGraphID = "orig-data"
	OneDimOptID   = "opt-path"
	TwoDimGraphID = "orig-data"
	TwoDimOptID   = "opt-path"

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

	OneDimChart = charting.Chart{
		ID:          OneDimChartID,
		Title:       "One-Dimensional Function",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "x",
		YAxisLabel:  "f(x)",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VarStart,
			VarEnd,
			VarStep,
		},
		Datasets: map[string]charting.Dataset{
			OneDimGraphID: &OneDimGraph,
			OneDimOptID:   &OneDimOptGraph,
		},
	}

	OneDimGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "f(x) = x + 1/ln(x)",
			BorderColor: charting.ColorBlue,
			BorderWidth: 2,
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
		PointRadius:     0,
	}

	OneDimOptGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Optimization Path",
			BorderColor: charting.ColorRed,
			BorderWidth: 3,
		},
		BackgroundColor: charting.ColorRed,
		PointRadius:     5,
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
			TwoDimGraphID: &TwoDimGraph,
			TwoDimOptID:   &TwoDimOptGraph,
		},
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

	TwoDimOptGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Optimization Path",
			BorderColor: charting.ColorYellow,
			BorderWidth: 4,
		},
		BackgroundColor: charting.ColorYellow,
		PointRadius:     6,
	}
)

func updateVariableLabel(c *charting.Chart, id string, label string) {
	// Check chart variables
	for i := range c.ChartVariables {
		if c.ChartVariables[i].ID == id {
			c.ChartVariables[i].Label = label
			return
		}
	}
	// Check graph variables in datasets
	for _, ds := range c.Datasets {
		fields := ds.GetFields()
		for i := range fields {
			if fields[i].ID == id {
				ds.UpdateVariableLabel(i, label)
				return
			}
		}
	}
}

func onedimf(x float64) float64 {
	if x <= 0 || x == 1 {
		return math.NaN()
	}
	return x + (1 / math.Log(x))
}

// z=-\ln\left(\left(\left(x-c_{1}\right)^{2}-\left(y-c_{2}\right)^{2}\right)^{2}+c_{1}^{2}-c_{2}^{2}\right)
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

func RenderOneDim(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()

	start, _ := req.GetChartVariable(OneDimChartID, VarStartID)
	end, _ := req.GetChartVariable(OneDimChartID, VarStartID) // oops, should be VarEndID
	step, _ := req.GetChartVariable(OneDimChartID, VarStepID)

	// wait, end, _ := req.GetChartVariable(OneDimChartID, VarStartID) was a typo in previous edit
	end, _ = req.GetChartVariable(OneDimChartID, VarEndID)

	methodIdx, _ := req.GetGraphVariable(OneDimChartID, OneDimGraphID, VarOptMethodID)
	tol, _ := req.GetGraphVariable(OneDimChartID, OneDimGraphID, VarOptTolID)
	lr, _ := req.GetGraphVariable(OneDimChartID, OneDimGraphID, VarOptLRID)
	samples, _ := req.GetGraphVariable(OneDimChartID, OneDimGraphID, VarOptSamplesID)
	x0, _ := req.GetGraphVariable(OneDimChartID, OneDimGraphID, VarOptX0ID)

	if start <= 0 || (start <= 1 && end >= 1) {
		return res.NewError("invalid start value: x must be > 0 and != 1")
	}

	n := int((end-start)/step) + 1
	x := make([]float64, 0, n)
	y := make([]float64, 0, n)

	for i := range n {
		xVal := start + float64(i)*step
		val := onedimf(xVal)
		if !math.IsNaN(val) && !math.IsInf(val, 0) {
			x = append(x, xVal)
			y = append(y, val)
		}
	}

	chartCopy := charting.CopyChart(OneDimChart)
	chartCopy.UpdatePointsForDataset(OneDimGraphID, x, y)

	f1d := func(args ...float64) float64 {
		return onedimf(args[0])
	}

	if methodIdx != 0 {
		var path [][]float64
		hideLine := false

		switch methodIdx {
		case 1: // Dichotomic Search
			p := dichotomicSearch(onedimf, start, end, tol)
			path = make([][]float64, len(p))
			for i, v := range p {
				path[i] = []float64{v}
			}
		case 2: // Random Search
			bounds := [][]float64{{start, end}}
			minPath, _ := randomSearchNdim(f1d, int(samples), bounds)
			path = minPath
			hideLine = true
		case 3: // Fast Descent
			path = fastDescentNdim(f1d, []float64{x0}, 0.01, tol, lr)
		}

		px := make([]float64, 0, len(path))
		py := make([]float64, 0, len(path))
		for _, p := range path {
			val := onedimf(p[0])
			if !math.IsNaN(val) && !math.IsInf(val, 0) && !math.IsNaN(p[0]) && !math.IsInf(p[0], 0) {
				px = append(px, p[0])
				py = append(py, val)
			}
		}

		if len(px) > 0 {
			chartCopy.UpdatePointsForDataset(OneDimOptID, px, py)
			chartCopy.Datasets[OneDimOptID].(*charting.GridDataset).HideLine = hideLine

			lastX := px[len(px)-1]
			lastY := py[len(py)-1]
			updateVariableLabel(&chartCopy, VarOptResultID, fmt.Sprintf("Result: min f(%.4f) = %.4f, path length = %d", lastX, lastY, len(path)))
		} else {
			chartCopy.UpdatePointsForDataset(OneDimOptID, nil, nil)
			updateVariableLabel(&chartCopy, VarOptResultID, "Result: Failed (NaN)")
		}
	} else {
		chartCopy.UpdatePointsForDataset(OneDimOptID, nil, nil)
		updateVariableLabel(&chartCopy, VarOptResultID, "Result: None")
	}

	chartCopy.GenerateLabels(2)
	res.AddChart(OneDimChartID, &chartCopy)

	return res
}

func RenderTwoDim(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()

	xStart, _ := req.GetChartVariable(TwoDimChartID, VarXStartID)
	xEnd, _ := req.GetChartVariable(TwoDimChartID, VarXEndID)
	xStep, _ := req.GetChartVariable(TwoDimChartID, VarXStepID)

	yStart, _ := req.GetChartVariable(TwoDimChartID, VarYStartID)
	yEnd, _ := req.GetChartVariable(TwoDimChartID, VarYEndID)
	yStep, _ := req.GetChartVariable(TwoDimChartID, VarYStepID)

	methodIdx, _ := req.GetGraphVariable(TwoDimChartID, TwoDimGraphID, VarOptMethodID)
	tol, _ := req.GetGraphVariable(TwoDimChartID, TwoDimGraphID, VarOptTolID)
	lr, _ := req.GetGraphVariable(TwoDimChartID, TwoDimGraphID, VarOptLRID)
	samples, _ := req.GetGraphVariable(TwoDimChartID, TwoDimGraphID, VarOptSamplesID)
	x0, _ := req.GetGraphVariable(TwoDimChartID, TwoDimGraphID, VarOptX0ID)
	y0, _ := req.GetGraphVariable(TwoDimChartID, TwoDimGraphID, VarOptY0ID)

	nx := int((xEnd-xStart)/xStep) + 1
	ny := int((yEnd-yStart)/yStep) + 1

	points := make([]any, 0, nx*ny)

	for i := range nx {
		for j := range ny {
			xVal := xStart + float64(i)*xStep
			yVal := yStart + float64(j)*yStep
			zVal := twodimf(xVal, yVal)

			if !math.IsNaN(zVal) && !math.IsInf(zVal, 0) {
				yCopy := yVal
				zCopy := zVal

				points = append(points, charting.HeatmapPoint{
					DataPoint: charting.DataPoint{X: xVal, Y: &yCopy},
					Value:     &zCopy,
				})
			}
		}
	}

	chartCopy := charting.CopyChart(TwoDimChart)
	chartCopy.UpdateDataForDataset(TwoDimGraphID, points)

	f2d := func(args ...float64) float64 {
		return twodimf(args[0], args[1])
	}

	if methodIdx != 0 {
		var path [][]float64
		hideLine := false

		switch methodIdx {
		case 1: // Random Search
			bounds := [][]float64{{xStart, xEnd}, {yStart, yEnd}}
			minPath, _ := randomSearchNdim(f2d, int(samples), bounds)
			path = minPath
			hideLine = true
		case 2: // Fast Descent
			path = fastDescentNdim(f2d, []float64{x0, y0}, 0.01, tol, lr)
		}

		px := make([]float64, 0, len(path))
		py := make([]float64, 0, len(path))
		for _, p := range path {
			val := twodimf(p[0], p[1])
			if !math.IsNaN(val) && !math.IsInf(val, 0) && !math.IsNaN(p[0]) && !math.IsInf(p[0], 0) && !math.IsNaN(p[1]) && !math.IsInf(p[1], 0) {
				px = append(px, p[0])
				py = append(py, p[1])
			}
		}

		if len(px) > 0 {
			chartCopy.UpdatePointsForDataset(TwoDimOptID, px, py)
			chartCopy.Datasets[TwoDimOptID].(*charting.GridDataset).HideLine = hideLine

			lastIdx := len(path) - 1
			last := path[lastIdx]
			updateVariableLabel(&chartCopy, VarOptResultID, fmt.Sprintf("Result: min f(%.4f, %.4f) = %.4f, path length = %d", last[0], last[1], twodimf(last[0], last[1]), len(path)-1))
		} else {
			chartCopy.UpdatePointsForDataset(TwoDimOptID, nil, nil)
			updateVariableLabel(&chartCopy, VarOptResultID, "Result: Failed (NaN)")
		}
	} else {
		chartCopy.UpdatePointsForDataset(TwoDimOptID, nil, nil)
		updateVariableLabel(&chartCopy, VarOptResultID, "Result: None")
	}

	chartCopy.GenerateLabels(2)
	res.AddChart(TwoDimChartID, &chartCopy)

	return res
}
