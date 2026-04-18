package neuron

import (
	"fmt"
	"labs/analysis"
	"labs/charting"
	"math"
	"strings"
)

var (
	BoundaryChart = charting.Chart{
		ID:          BoundaryChartID,
		Title:       "Decision Boundary & Probability Heatmap",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "x",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "y",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphHeatmapID:  &ProbabilityHeatmapDataset,
			GraphClass0ID:   &BoundaryClass0Dataset,
			GraphClass1ID:   &BoundaryClass1Dataset,
			GraphBoundaryID: &DecisionBoundaryLineDataset,
		},
		ChartVariables: append(
			SharedVariables,
			ActivationFuncField,
			DisplayFormula,
		),
	}

	ProbabilityHeatmapDataset = charting.HeatmapDataset{
		BaseDataset: charting.BaseDataset{
			Label: "Probability Area",
			GraphVariables: []charting.MutableField{
				VarHeatmapPrecision,
			},
		},
	}

	BoundaryClass0Dataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Class 0",
			BorderColor: charting.ColorRed,
		},
		HideLine:    true,
		PointRadius: 5,
	}

	BoundaryClass1Dataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Class 1",
			BorderColor: charting.ColorBlue,
		},
		HideLine:    true,
		PointRadius: 5,
	}

	DecisionBoundaryLineDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Decision Boundary",
			BorderColor: charting.ColorAmber,
			BorderWidth: 4,
		},
		HideLine:    false,
		PointRadius: 0,
	}
)

var names = []string{"Sigmoid", "Tanh", "ReLU"}

func RenderBoundary(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()
	if err := ensureTrained(req); err != nil {
		return res.NewErrorf("error in training: %v", err)
	}

	chartCopy := charting.CopyChart(BoundaryChart)

	alpha := req.GetVariable(VarAlphaID)
	actIdx := int(req.GetVariable(VarActivationID))
	varPrecision, _ := req.GetGraphVariable(BoundaryChartID, GraphHeatmapID, VarHeatmapPrecisionID)
	precision := int(varPrecision)

	act, _ := getActivation(actIdx, alpha)

	// Get specific training result
	trainRes := lastTrainResults[actIdx]

	finalW := trainRes.WeightsHistory[len(trainRes.WeightsHistory)-1]
	finalB := trainRes.BiasHistory[len(trainRes.BiasHistory)-1]

	var sb strings.Builder
	fmt.Fprintf(&sb, "Formula: %s(%.4fx + %.4fy + %.4f)", names[actIdx], finalW[0], finalW[1], finalB)
	chartCopy.UpdateVariableLabel(VarFormulaID, sb.String())
	fmt.Println(sb.String())

	minX, maxX := analysis.MinMax(trainData.X)
	minY, maxY := analysis.MinMax(trainData.Y)

	heatPoints := make([]any, 0)
	for i := 0; i <= precision; i++ {
		gx := minX + (maxX-minX)*float64(i)/float64(precision)
		for j := 0; j <= precision; j++ {
			gy := minY + (maxY-minY)*float64(j)/float64(precision)
			z := finalW[0]*gx + finalW[1]*gy + finalB
			conf := act(z)

			// Safety: Plotly/JSON cannot handle NaN or Inf
			if math.IsNaN(conf) || math.IsInf(conf, 0) {
				conf = 0.5 // Default to uncertain
			}

			heatPoints = append(heatPoints, &charting.HeatmapPoint{
				DataPoint: charting.DataPoint{X: gx, Y: &gy},
				Value:     &conf,
			})
		}
	}
	chartCopy.UpdateDataForDataset(GraphHeatmapID, heatPoints)

	c0, c1 := splitData(trainData)
	chartCopy.UpdateDataPointsForDataset(GraphClass0ID, c0)
	chartCopy.UpdateDataPointsForDataset(GraphClass1ID, c1)

	// Boundary Line: y = -(w1*x + b) / w2
	// Handle division by zero (vertical line)
	var yStart, yEnd float64
	lineMinX, lineMaxX := minX, maxX
	if math.Abs(finalW[1]) < 1e-9 {
		// Almost vertical line
		yStart, yEnd = minY, maxY
		if math.Abs(finalW[0]) > 1e-9 {
			lineMinX = -finalB / finalW[0]
		} else {
			lineMinX = 0
		}
		lineMaxX = lineMinX
	} else {
		yStart = -(finalW[0]*lineMinX + finalB) / finalW[1]
		yEnd = -(finalW[0]*lineMaxX + finalB) / finalW[1]
	}

	// Final JSON safety check for all coordinates
	safe := func(f float64) float64 {
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return 0
		}
		return f
	}

	chartCopy.UpdateDataPointsForDataset(GraphBoundaryID, []charting.DataPoint{
		{X: safe(lineMinX), Y: Pointer(safe(yStart))},
		{X: safe(lineMaxX), Y: Pointer(safe(yEnd))},
	})

	res.AddChart(BoundaryChartID, &chartCopy)
	return res
}

func Pointer[T any](v T) *T {
	return &v
}
