package cluster

import (
	"errors"
	"fmt"
	"labs/charting"
	"math"
	"math/rand"
)

const (
	KmeansGraphID = "kmeans"
	KmeansChartID = "kmeans"

	VariableCenroidsChoiseOptionID = "centroids-choise-method"
	VariableNumCentriodsID         = "centroids-num"
)

var (
	VariableNumCentroids = charting.MutableField{
		ID:      VariableNumCentriodsID,
		Label:   "Number of cetriods",
		Default: 5,
		Min:     1,
		Max:     100,
		Step:    1,
		Control: charting.ControlNumber,
	}

	VariableCentroidsChoiseOption = charting.MutableField{
		ID:      VariableCenroidsChoiseOptionID,
		Label:   "What points should be selected",
		Control: charting.ControlSelect,
		Default: 0,
		Options: []string{
			"first n points",
			"last n points",
			"random n points",
		},
	}

	KmeansChart = charting.Chart{
		ID:          KmeansChartID,
		Title:       "Kmeans Clusterization",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "X",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "Y",
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VariableNumCentroids,
			VariableCentroidsChoiseOption,
		},
	}
)

func RenderKmeans(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadPoints(); err != nil {
		return res.NewError(err.Error())
	}

	option, ok := req.GetChartVariable(KmeansChartID, VariableCenroidsChoiseOptionID)
	if !ok {
		option = VariableCentroidsChoiseOption.Default
	}
	num_centroids, ok := req.GetChartVariable(KmeansChartID, VariableNumCentriodsID)
	if !ok {
		num_centroids = VariableNumCentroids.Default
	}

	if num_centroids < 1 {
		return res.NewErrorf("number of centroids too small")
	}

	centroids := make([]charting.DataPoint, int(num_centroids))
	for i := range int(num_centroids) {
		switch int(option) {
		case 0: // first points
			centroids[i] = points[i]
		case 1: // last points
			centroids[i] = points[len(points)-1-i]
		case 2: // random points
			centroids[i] = points[rand.Intn(len(points))]
		default:
			return res.NewErrorf("unrecognized option index: %.0f", option)
		}
	}

	labels, centroids, err := kmeans(points, centroids, 1000)
	if err != nil {
		return res.NewErrorf("failed to cluster data: %s", err.Error())
	}

	copyChart := charting.CopyChart(KmeansChart)

	clusterData(labels, len(centroids), &copyChart)

	res = charting.NewRenderResponse()
	if int(option) == 2 { // dont cache random centroids, as they should be different on each render
		res.CachePolicy = charting.CachePolicyDontCache
	}
	res.AddChart(copyChart.ID, &copyChart)

	return res
}

func kmeans(points []charting.DataPoint, centroids []charting.DataPoint, maxIter uint32) ([]int, []charting.DataPoint, error) {
	var labels []int
	for i := range maxIter {
		labels = assignPointsToCentroids(points, centroids)
		newCentroids := updateCentroids(points, labels, uint32(len(centroids)))

		if fmt.Sprint(newCentroids) == fmt.Sprint(centroids) {
			break
		}
		centroids = newCentroids
		if i == maxIter-1 {
			return nil, nil, errors.New("k-means: centroids haven't coverged")
		}
	}

	return labels, centroids, nil
}

func assignPointsToCentroids(points []charting.DataPoint, centroids []charting.DataPoint) []int {
	labels := make([]int, len(points))
	for i, p := range points {
		best, minDist := 0, math.MaxFloat64
		for j, c := range centroids {
			if d := euclidianDist(p, c); d < minDist {
				minDist, best = d, j
			}
		}
		labels[i] = best
	}
	return labels
}

func updateCentroids(points []charting.DataPoint, labels []int, k uint32) []charting.DataPoint {
	sums := make([]charting.DataPoint, k)
	counts := make([]int, k)
	for i, p := range points {
		sums[labels[i]].X += p.X
		sums[labels[i]].Y += p.Y
		counts[labels[i]]++
	}
	centroids := make([]charting.DataPoint, k)
	for i := range centroids {
		if counts[i] > 0 {
			centroids[i] = charting.DataPoint{X: sums[i].X / float64(counts[i]), Y: sums[i].Y / float64(counts[i])}
		}
	}
	return centroids
}

func euclidianDist(a, b charting.DataPoint) float64 {
	return math.Sqrt(math.Pow(a.X-b.X, 2) + math.Pow(a.Y-b.Y, 2))
}
