package cluster

import (
	"cmp"
	"fmt"
	"labs/charting"
	"math"
	"math/rand"
	"slices"
)

const (
	SilhouetteChartID    = "silhouette"
	SilhouetteGraphFmtID = "silhouette-%d"

	VariableAlgChoiseID = "algorithm-choise"
)

var SilhouetteChart = charting.Chart{
	ID:         SilhouetteChartID,
	Title:      "Silhouette Plot",
	Type:       charting.ChartTypeMultiBar,
	XAxisLabel: "Points (sorted by score)",
	YAxisLabel: "Silhouette coefficient",
	Datasets:   map[string]charting.Dataset{},
	ChartVariables: []charting.MutableField{
		{
			ID:      VariableAlgChoiseID,
			Label:   "Algorithm",
			Control: charting.ControlSelect,
			Options: []string{"K-Means", "Simple"},
			Default: 0,
		},
		VariableNumCentroids, // reuse same field def as in kmeans lab
		VariableThreshold,
		VariableCentroidsChoiseOption,
	},
}

// SilhouetteScores computes per-point silhouette coefficients.
// points: all data points
// labels: cluster assignment per point (same length as points)
// returns: scores in the same order as points
func SilhouetteScores(points []charting.DataPoint, labels []int) []float64 {
	n := len(points)
	scores := make([]float64, n)

	for i := range n {
		ci := labels[i]

		var a, b float64

		// a(i): mean distance to points in same cluster
		var sumA float64
		var countA int
		for j := range n {
			if j == i || labels[j] != ci {
				continue
			}
			sumA += dist(points[i], points[j])
			countA++
		}
		if countA > 0 {
			a = sumA / float64(countA)
		}

		// b(i): mean distance to points in nearest other cluster
		b = math.MaxFloat64
		// find unique clusters
		clusterSums := map[int]float64{}
		clusterCounts := map[int]int{}
		for j := range n {
			if labels[j] == ci {
				continue
			}
			clusterSums[labels[j]] += dist(points[i], points[j])
			clusterCounts[labels[j]]++
		}
		for cj, sum := range clusterSums {
			mean := sum / float64(clusterCounts[cj])
			if mean < b {
				b = mean
			}
		}

		if b == math.MaxFloat64 {
			// only one cluster exists
			scores[i] = 0
			continue
		}

		scores[i] = (b - a) / math.Max(a, b)
	}

	return scores
}

func dist(a, b charting.DataPoint) float64 {
	ay, by := 0.0, 0.0
	if a.Y != nil {
		ay = *a.Y
	}
	if b.Y != nil {
		by = *b.Y
	}
	dx := a.X - b.X
	dy := ay - by
	return math.Sqrt(dx*dx + dy*dy)
}

func RenderSilhouette(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadPoints(); err != nil {
		return res.NewError(err.Error())
	}

	algChoice, ok := req.GetChartVariable(SilhouetteChartID, VariableAlgChoiseID)
	if !ok {
		algChoice = 0
	}
	num_centroids, ok := req.GetChartVariable(SilhouetteChartID, VariableNumCentriodsID)
	if !ok {
		num_centroids = VariableNumCentroids.Default
	}
	if algChoice == 1 {
		fmt.Print("\n\nSIMPLE ALGORITHM DETECTED, setting amount of initial centroids to 1\n\n")
		num_centroids = 1
	}
	threshold, ok := req.GetChartVariable(SilhouetteChartID, VariableThresholdID)
	if !ok {
		threshold = VariableThreshold.Default
	}
	option, ok := req.GetChartVariable(SilhouetteChartID, VariableCenroidsChoiseOptionID)
	if !ok {
		option = VariableCentroidsChoiseOption.Default
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

	// run the chosen algorithm
	var labels []int
	var err error
	switch int(algChoice) {
	case 0:
		labels, _, err = kmeans(points, centroids, 1000)
	case 1:
		labels, _ = simpleClustering(points, centroids[0], threshold)
	default:
		return res.NewErrorf("unrecognized algorithm: %.0f", algChoice)
	}
	if err != nil {
		return res.NewErrorf("clustering failed: %s", err.Error())
	}

	scores := SilhouetteScores(points, labels)

	// group scores by cluster, sorted descending
	maxCluster := 0
	for _, label := range labels {
		if label > maxCluster {
			maxCluster = label
		}
	}
	k := maxCluster + 1
	clusterScores := make([][]float64, k)
	for i, label := range labels {
		clusterScores[label] = append(clusterScores[label], scores[i])
	}

	strlabels := make([]string, len(points))
	for i := range points {
		strlabels[i] = fmt.Sprintf("P%d", i+1)
	}

	colors := [...]charting.Color{
		charting.ColorBlue, charting.ColorRed, charting.ColorGreen,
		charting.ColorViolet, charting.ColorAmber, charting.ColorCyan,
		charting.ColorPink, charting.ColorOrange, charting.ColorTeal,
		charting.ColorIndigo,
	}

	copyChart := charting.CopyChart(SilhouetteChart)
	copyChart.Labels = strlabels

	for c := range k {
		datasetID := fmt.Sprintf(SilhouetteGraphFmtID, c)
		color := colors[c%len(colors)]

		slices.SortFunc(clusterScores[c], func(a, b float64) int {
			return cmp.Compare(b, a) // descending
		})

		copyChart.Datasets[datasetID] = &charting.CategoricalDataset{
			BaseDataset: charting.BaseDataset{
				Label:       fmt.Sprintf("Cluster %d", c),
				BorderColor: color,
				BorderWidth: 1,
			},
			Data:            charting.ToFloat64PtrSlice(clusterScores[c]),
			BackgroundColor: []charting.Color{color},
		}
	}

	res = charting.NewRenderResponse()
	res.AddChart(copyChart.ID, &copyChart)
	return res
}
