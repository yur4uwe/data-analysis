package neuron

import (
	"fmt"
	"labs/charting"
	"math"
	"strings"
)

var (
	TrajectoryChart = charting.Chart{
		ID:          TrajectoryChartID,
		Title:       "Weight Trajectory & Loss Landscape",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "w1",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "w2",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphLandscapeID: &LandscapeHeatmapDataset,
			GraphPathID:      &LearningPathDataset,
		},
		ChartVariables: append(SharedVariables, ActivationFuncField, DisplayFormula),
	}

	LandscapeHeatmapDataset = charting.HeatmapDataset{
		BaseDataset: charting.BaseDataset{
			Label: "Loss Landscape",
			GraphVariables: []charting.MutableField{
				VarHeatmapPrecision,
			},
		},
	}

	LearningPathDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Learning Path",
			BorderColor: charting.ColorAmber,
			BorderWidth: 3,
		},
		HideLine:    false,
		PointRadius: 2,
	}
)

func roundTo(val float64, precision float64) float64 {
	if precision == 0 {
		return val
	}
	return math.Round(val/precision) * precision
}

func RenderTrajectory(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()
	if err := ensureTrained(req); err != nil {
		return res.NewErrorf("error in training: %v", err)
	}

	chartCopy := charting.CopyChart(TrajectoryChart)

	safe := func(f float64) float64 {
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return 0
		}
		return f
	}

	alpha := req.GetVariable(VarAlphaID)
	actIdx := int(req.GetVariable(VarActivationID))

	// Ensure precision is sane
	varPrecision, _ := req.GetGraphVariable(TrajectoryChartID, GraphLandscapeID, VarHeatmapPrecisionID)
	precision := int(varPrecision)
	margin := 0.5

	act, _ := getActivation(actIdx, alpha)

	// Get specific training result
	trainRes := lastTrainResults[actIdx]

	finalW := trainRes.WeightsHistory[len(trainRes.WeightsHistory)-1]
	finalB := trainRes.BiasHistory[len(trainRes.BiasHistory)-1]
	var sb strings.Builder
	fmt.Fprintf(&sb, "Epochs: %d | Test Acc: %.2f%% | Formula: %s(%.4fx + %.4fy + %.4f)", trainRes.EpochsTrained, trainRes.TestAccuracy*100, names[actIdx], finalW[0], finalW[1], finalB)
	chartCopy.UpdateVariableLabel(VarFormulaID, sb.String())

	// 2. Determine Weight Range from History
	var minW1, maxW1, minW2, maxW2 float64
	initialized := false

	for _, w := range trainRes.WeightsHistory {
		w1, w2 := w[0], w[1]
		if math.IsNaN(w1) || math.IsInf(w1, 0) || math.IsNaN(w2) || math.IsInf(w2, 0) {
			continue
		}
		if !initialized {
			minW1, maxW1 = w1, w1
			minW2, maxW2 = w2, w2
			initialized = true
		} else {
			if w1 < minW1 {
				minW1 = w1
			}
			if w1 > maxW1 {
				maxW1 = w1
			}
			if w2 < minW2 {
				minW2 = w2
			}
			if w2 > maxW2 {
				maxW2 = w2
			}
		}
	}

	if !initialized {
		minW1, maxW1 = -1, 1
		minW2, maxW2 = -1, 1
	}

	// Add margin and ensure range is not zero
	minW1 -= margin
	maxW1 += margin
	minW2 -= margin
	maxW2 += margin
	if math.Abs(maxW1-minW1) < 0.001 {
		minW1 -= 0.1
		maxW1 += 0.1
	}
	if math.Abs(maxW2-minW2) < 0.001 {
		minW2 -= 0.1
		maxW2 += 0.1
	}

	// Calculate step size for rounding
	step1 := (maxW1 - minW1) / float64(precision)
	step2 := (maxW2 - minW2) / float64(precision)

	if math.IsNaN(finalB) || math.IsInf(finalB, 0) {
		finalB = 0
	}

	heatPoints := make([]any, 0)
	for i := 0; i <= precision; i++ {
		gw1 := safe(roundTo(minW1+float64(i)*step1, step1/1000))
		for j := 0; j <= precision; j++ {
			gw2 := safe(roundTo(minW2+float64(j)*step2, step2/1000))

			forward := newForward(act, []float64{gw1, gw2}, finalB)
			var totalLoss float64
			for k := range trainData.X {
				pred := forward([]float64{trainData.X[k], trainData.Y[k]})
				target := 0.0
				if trainData.Class[k] {
					target = 1.0
				}
				diff := pred - target
				totalLoss += 0.5 * diff * diff
			}
			avgLoss := totalLoss / float64(len(trainData.X))
			if math.IsNaN(avgLoss) || math.IsInf(avgLoss, 0) {
				avgLoss = 1.0 // Penalty value
			}

			valY := gw2 // local copy for pointer
			heatPoints = append(heatPoints, &charting.HeatmapPoint{
				DataPoint: charting.DataPoint{X: gw1, Y: &valY},
				Value:     &avgLoss,
			})
		}
	}
	chartCopy.UpdateDataForDataset(GraphLandscapeID, heatPoints)

	trajPoints := make([]charting.DataPoint, 0, len(trainRes.WeightsHistory))
	for _, w := range trainRes.WeightsHistory {
		w1, w2 := safe(w[0]), safe(w[1])
		val := w2 // local copy for pointer
		trajPoints = append(trajPoints, charting.DataPoint{X: w1, Y: &val})
	}
	chartCopy.UpdateDataPointsForDataset(GraphPathID, trajPoints)

	res.AddChart(TrajectoryChartID, &chartCopy)
	return res
}
