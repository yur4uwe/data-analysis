package cluster

import (
	"labs/charting"
	"math/rand"
)

const (
	SimpleChartID = "simple"

	VariableThresholdID = "threshold"
)

var (
	VariableThreshold = charting.MutableField{
		ID:      VariableThresholdID,
		Label:   "Distance Threshold",
		Default: 5,
		Min:     0.0,
		Max:     10,
		Step:    0.1,
		Control: charting.ControlNumber,
	}

	VariableChoiseOption = charting.MutableField{
		ID:      VariableCenroidsChoiseOptionID,
		Label:   "What points should be selected as 1st centroid",
		Control: charting.ControlSelect,
		Default: 0,
		Options: []string{
			"first point",
			"last point",
			"random point",
		},
	}

	SimpleChart = charting.Chart{
		ID:          SimpleChartID,
		Title:       "Simple Clusterization",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "X",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "Y",
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VariableThreshold,
			VariableChoiseOption,
		},
	}
)

func RenderSimple(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadPoints(); err != nil {
		return res.NewError(err.Error())
	}

	threshold, ok := req.GetChartVariable(SimpleChartID, VariableThresholdID)
	if !ok {
		threshold = VariableThreshold.Default
	}
	option, ok := req.GetChartVariable(SimpleChartID, VariableCenroidsChoiseOptionID)
	if !ok {
		option = VariableChoiseOption.Default
	}

	copyChart := charting.CopyChart(SimpleChart)

	var initCentroid charting.DataPoint
	switch int(option) {
	case 0: // first point
		initCentroid = points[0]
	case 1: // last point
		initCentroid = points[len(points)-1]
	case 2: // random point
		initCentroid = points[rand.Intn(len(points))]
	default:
		return res.NewErrorf("unrecognized option index: %.0f", option)
	}

	labels, centroids := simpleClustering(points, initCentroid, threshold)

	clusterData(labels, centroids, &copyChart)

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	if int(option) == 2 { // dont cache random centroids, as they should be different on each render
		res.CachePolicy = charting.CachePolicyDontCache
	}

	return res
}

func simpleClustering(points []charting.DataPoint, initCentroid charting.DataPoint, T float64) (labels []int, centroids []charting.DataPoint) {
	if len(points) == 0 {
		return nil, nil
	}

	labels = make([]int, len(points))
	centroids = []charting.DataPoint{initCentroid}

	for i, p := range points[1:] {
		best, minDist := 0, euclidianDist(p, centroids[0])
		for j := 1; j < len(centroids); j++ {
			if d := euclidianDist(p, centroids[j]); d < minDist {
				minDist, best = d, j
			}
		}

		if minDist > T {
			labels[i+1] = len(centroids)
			centroids = append(centroids, p)
		} else {
			labels[i+1] = best
		}
	}

	return labels, centroids
}
