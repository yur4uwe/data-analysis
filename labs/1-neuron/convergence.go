package neuron

import (
	"fmt"
	"labs/charting"
	"math"
	"strings"
)

var (
	ConvergenceChart = charting.Chart{
		ID:          ConvergenceChartID,
		Title:       "Training Convergence",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Epoch",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "MSE",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphLossSigmoidID: &LossSigmoidDataset,
			GraphLossTanhID:    &LossTanhDataset,
			GraphLossReLUID:    &LossReLUDataset,
		},
		ChartVariables: append(SharedVariables, DisplayFormula),
	}

	LossSigmoidDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Sigmoid",
			BorderColor: charting.ColorViolet,
			BorderWidth: 2,
		},
		HideLine:    false,
		PointRadius: 2,
	}

	LossTanhDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Tanh",
			BorderColor: charting.ColorAmber,
			BorderWidth: 2,
		},
		HideLine:    false,
		PointRadius: 2,
	}

	LossReLUDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "ReLU",
			BorderColor: charting.ColorEmerald,
			BorderWidth: 2,
		},
		HideLine:    false,
		PointRadius: 2,
	}
)

func RenderConvergence(req *charting.RenderRequest) (res *charting.RenderResponse) {
	res = charting.NewRenderResponse()
	if err := ensureTrained(req); err != nil {
		return res.NewErrorf("error in training: %v", err)
	}

	chartCopy := charting.CopyChart(ConvergenceChart)

	var sb strings.Builder
	sb.WriteString("Formulas:\n")

	// Update datasets for all three activations
	updateLoss := func(graphID string, actIdx int) {
		trainRes := lastTrainResults[actIdx]
		lossPoints := make([]charting.DataPoint, 0, len(trainRes.LossHistory))
		for i, l := range trainRes.LossHistory {
			if math.IsNaN(l) || math.IsInf(l, 0) {
				l = 1.0 // Penalty for explosion
			}
			val := l
			lossPoints = append(lossPoints, charting.DataPoint{X: float64(i), Y: &val})
		}
		fmt.Fprintf(&sb, "%s(%.4fx + %.4fy + %.4f)\n", names[actIdx], trainRes.WeightsHistory[len(trainRes.WeightsHistory)-1][0], trainRes.WeightsHistory[len(trainRes.WeightsHistory)-1][1], trainRes.BiasHistory[len(trainRes.BiasHistory)-1])
		chartCopy.UpdateDataPointsForDataset(graphID, lossPoints)
	}

	updateLoss(GraphLossSigmoidID, 0)
	updateLoss(GraphLossTanhID, 1)
	updateLoss(GraphLossReLUID, 2)

	chartCopy.UpdateVariableLabel(VarFormulaID, sb.String())

	res.AddChart(ConvergenceChartID, &chartCopy)
	return res
}
