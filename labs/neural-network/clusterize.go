package network

import (
	"fmt"
	"labs/charting"
	"math"
)

func norm(x *ClusterizationPoint) {
	euc := math.Sqrt(x.X*x.X + x.Y*x.Y)
	x.X /= euc
	x.Y /= euc
}

func label(w [][]float64, x ClusterizationPoint) int {
	y := forward(w, []float64{x.X, x.Y})
	maxIdx := 0
	maxVal := y[0]
	for i := range y {
		if y[i] > maxVal {
			maxIdx = i
			maxVal = y[i]
		}
	}
	if maxVal > 0.5 {
		return maxIdx
	}
	return -1
}

type ClusterResult struct {
	Clusters []int
	Centers  [][]float64
}

func RenderClusters(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadData(); err != nil {
		return res.NewErrorf("error loading data: %v", err)
	}

	numClusters := uint32(req.GetVariable(NumClustersID))
	maxEpochs := uint32(req.GetVariable(NumEpochsID))
	lr := req.GetVariable(LearningRateID)
	if numClusters == 0 {
		numClusters = 3
	}

	// Work on a copy of data because clusterize normalizes it in-place
	dataCopy := make([]ClusterizationPoint, len(data))
	copy(dataCopy, data)

	res_cluster := clusterize(dataCopy, numClusters, maxEpochs, lr)

	chartCopy := charting.CopyChart(ClustersChart)

	// First, calculate cluster points and update centers to original space
	for cluster := range int(numClusters) {
		cluster_points := make([]charting.DataPoint, 0)
		centerX, centerY := 0.0, 0.0
		count := 0
		for i := range data {
			if res_cluster.Clusters[i] == cluster {
				cluster_points = append(cluster_points, charting.DataPoint{
					X: data[i].X,
					Y: &data[i].Y,
				})
				centerX += data[i].X
				centerY += data[i].Y
				count++
			}
		}

		if count > 0 {
			res_cluster.Centers[cluster][0] = centerX / float64(count)
			res_cluster.Centers[cluster][1] = centerY / float64(count)
		}

		key := fmt.Sprintf("cluster-%d", cluster)
		chartCopy.Datasets[key] = &charting.GridDataset{
			BaseDataset: charting.BaseDataset{
				Label:       fmt.Sprintf("Cluster %d", cluster),
				BorderColor: GetClusterColor(cluster),
				BorderWidth: 2,
			},
			BackgroundColor: charting.ColorTransparent,
			PointRadius:     4,
			Data:            cluster_points,
			HideLine:        true,
		}
	}

	// Now create the centroids dataset using recalculated centers
	centroidLabels := make([]string, len(res_cluster.Centers))
	centroids := make([]charting.DataPoint, len(res_cluster.Centers))
	for i := range res_cluster.Centers {
		centroidLabels[i] = fmt.Sprintf("Cluster %d", i)
		centroids[i] = charting.DataPoint{
			X: res_cluster.Centers[i][0],
			Y: &res_cluster.Centers[i][1],
		}
	}

	chartCopy.Datasets["centroids"] = &charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Centroids",
			BorderColor: "#000000",
			BorderWidth: 3,
			DataLabels:  centroidLabels,
			ZIndex:      10,
		},
		BackgroundColor: "#ffffff",
		PointRadius:     12,
		PointStyle:      "star",
		Data:            centroids,
		HideLine:        true,
	}

	// Add unclustered points if any
	unclustered := make([]charting.DataPoint, 0)
	for i := range data {
		if res_cluster.Clusters[i] == -1 {
			unclustered = append(unclustered, charting.DataPoint{
				X: data[i].X,
				Y: &data[i].Y,
			})
		}
	}
	if len(unclustered) > 0 {
		chartCopy.Datasets["unclustered"] = &charting.GridDataset{
			BaseDataset: charting.BaseDataset{
				Label:       "Unclustered",
				BorderColor: charting.ColorSlate,
				BorderWidth: 1,
			},
			BackgroundColor: charting.ColorTransparent,
			PointRadius:     3,
			Data:            unclustered,
			HideLine:        true,
		}
	}

	res = charting.NewRenderResponse()
	res.AddChart(ClustersChartID, &chartCopy)
	return res
}

func clusterize(data []ClusterizationPoint, numClusters uint32, maxEpochs uint32, lr float64) *ClusterResult {
	for i := range data {
		norm(&data[i])
	}

	res := &ClusterResult{}

	trainingRes := train(data, maxEpochs, numClusters, lr)

	weights := trainingRes.WeightsHistory[len(trainingRes.WeightsHistory)-1]

	clusters := make([]int, len(data))
	for i := range data {
		clusters[i] = label(weights, data[i])
	}

	res.Clusters = clusters

	clusterCenters := make([][]float64, numClusters)
	clusterCounts := make([]int, numClusters)
	for i := range clusterCenters {
		clusterCenters[i] = make([]float64, 2)
	}
	for i := range data {
		if clusters[i] != -1 {
			clusterCenters[clusters[i]][0] += data[i].X
			clusterCenters[clusters[i]][1] += data[i].Y
			clusterCounts[clusters[i]]++
		}
	}
	for i := range clusterCenters {
		if clusterCounts[i] > 0 {
			clusterCenters[i][0] /= float64(clusterCounts[i])
			clusterCenters[i][1] /= float64(clusterCounts[i])
		}
	}

	res.Centers = clusterCenters

	return res
}
