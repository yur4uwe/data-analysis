package network

import (
	"fmt"
	"labs/charting"
	"math"
)

func RenderWeights(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadData(); err != nil {
		return res.NewErrorf("error loading data: %v", err)
	}
	numClusters := uint32(req.GetVariable(NumClustersID))
	trainingRes := train(data, 100, numClusters, 0.1)

	chartCopy := charting.CopyChart(WeightsChart)

	// Add Unit Circle for reference
	circlePoints := make([]charting.DataPoint, 101)
	for i := 0; i <= 100; i++ {
		angle := float64(i) * 2 * math.Pi / 100
		y := math.Sin(angle)
		circlePoints[i] = charting.DataPoint{
			X: math.Cos(angle),
			Y: &y,
		}
	}
	chartCopy.UpdateDataPointsForDataset(UnitCircleID, circlePoints)

	// Subsample history to around 100-200 frames for smooth but efficient animation
	step := len(trainingRes.WeightsHistory) / 150
	if step == 0 {
		step = 1
	}

	framesCount := (len(trainingRes.WeightsHistory) + step - 1) / step

	// Create a dataset for each cluster
	for clusterIdx := range numClusters {
		clusterID := fmt.Sprintf("cluster-%d", clusterIdx)

		frames := make([][]charting.DataPoint, framesCount)
		for f := range framesCount {
			historyIdx := f * step
			if historyIdx >= len(trainingRes.WeightsHistory) {
				historyIdx = len(trainingRes.WeightsHistory) - 1
			}

			w := trainingRes.WeightsHistory[historyIdx][clusterIdx]
			val := w[1]
			frames[f] = []charting.DataPoint{
				{X: 0, Y: ptr(0)},
				{X: w[0], Y: &val},
			}
		}

		ds := &charting.AnimationDataset{
			BaseDataset: charting.BaseDataset{
				Label:       fmt.Sprintf("Cluster %d", clusterIdx),
				Type:        charting.ChartTypeScatter,
				BorderColor: getClusterColor(int(clusterIdx)),
				BorderWidth: 2,
			},
			Data:   frames[0],
			Frames: frames,
		}

		chartCopy.Datasets[clusterID] = ds
	}

	res = charting.NewRenderResponse()
	res.AddChart(WeightsChartID, &chartCopy)
	return res
}

func ptr(v float64) *float64 {
	return &v
}

func getClusterColor(idx int) charting.Color {
	colors := []charting.Color{
		charting.ColorBlue,
		charting.ColorRed,
		charting.ColorGreen,
		charting.ColorPurple,
		charting.ColorOrange,
	}
	return colors[idx%len(colors)]
}
