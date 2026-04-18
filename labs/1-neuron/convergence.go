package neuron

import (
	"fmt"
	"labs/charting"
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
			GraphLossSigmoidID:    &LossSigmoidDataset,
			GraphLossTanhID:       &LossTanhDataset,
			GraphLossReLUID:       &LossReLUDataset,
			"val-loss-sigmoid":    &ValLossSigmoidDataset,
			"val-loss-tanh":       &ValLossTanhDataset,
			"val-loss-relu":       &ValLossReLUDataset,
		},
		ChartVariables: append(SharedVariables, DisplayFormula),
	}

	ValLossSigmoidDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Sigmoid (Val)",
			BorderColor: charting.ColorViolet,
			BorderWidth: 1,
		},
		HideLine:    false,
		PointRadius: 0,
	}

	ValLossTanhDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Tanh (Val)",
			BorderColor: charting.ColorAmber,
			BorderWidth: 1,
		},
		HideLine:    false,
		PointRadius: 0,
	}

	ValLossReLUDataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "ReLU (Val)",
			BorderColor: charting.ColorEmerald,
			BorderWidth: 1,
		},
		HideLine:    false,
		PointRadius: 0,
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
	updateLoss := func(graphID string, valGraphID string, actIdx int) {
		trainRes := lastTrainResults[actIdx]

		// Training Loss
		lossPoints := make([]charting.DataPoint, 0, len(trainRes.LossHistory))
		for i, l := range trainRes.LossHistory {
			val := l
			lossPoints = append(lossPoints, charting.DataPoint{X: float64(i), Y: &val})
		}
		chartCopy.UpdateDataPointsForDataset(graphID, lossPoints)

		// Validation Loss
		valLossPoints := make([]charting.DataPoint, 0, len(trainRes.ValidationLossHistory))
		for i, l := range trainRes.ValidationLossHistory {
			val := l
			valLossPoints = append(valLossPoints, charting.DataPoint{X: float64(i), Y: &val})
		}
		chartCopy.UpdateDataPointsForDataset(valGraphID, valLossPoints)

		fmt.Fprintf(&sb, "%s: Epochs: %d | Test Acc: %.2f%% | Formula: %s(%.4fx + %.4fy + %.4f)\n",
			names[actIdx], trainRes.EpochsTrained, trainRes.TestAccuracy*100, names[actIdx],
			trainRes.WeightsHistory[len(trainRes.WeightsHistory)-1][0],
			trainRes.WeightsHistory[len(trainRes.WeightsHistory)-1][1],
			trainRes.BiasHistory[len(trainRes.BiasHistory)-1])
	}

	updateLoss(GraphLossSigmoidID, "val-loss-sigmoid", 0)
	updateLoss(GraphLossTanhID, "val-loss-tanh", 1)
	updateLoss(GraphLossReLUID, "val-loss-relu", 2)

	chartCopy.UpdateVariableLabel(VarFormulaID, sb.String())

	res.AddChart(ConvergenceChartID, &chartCopy)
	return res
}
