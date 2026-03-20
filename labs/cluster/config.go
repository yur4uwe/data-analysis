package cluster

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	LabID = "6"
)

var (
	Config = charting.NewLabConfig(
		LabID,
		"Data Clustering",
		map[string]*charting.Chart{
			SimpleChartID: &SimpleChart,
			KmeansChartID: &KmeansChart,
		},
	)

	Metadata = Config.Lab

	points = []charting.DataPoint{}
)

func loadPoints() error {
	if len(points) != 0 {
		return nil
	}

	f, err := os.Open("./data/lab_6_var_12.csv")
	if err != nil {
		return fmt.Errorf("clustering points chart: error while reading file: %s", err.Error())
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.Comma = ','
	if err := d.Decode(&points); err != nil {
		return fmt.Errorf("clustering points chart: error while decoding csv: %s", err.Error())
	}

	return nil
}

func clusterData(labels []int, centroids []charting.DataPoint, chart *charting.Chart) {
	colors := [...]string{
		charting.ColorAmber,
		charting.ColorBlue,
		charting.ColorCyan,
		charting.ColorEmerald,
		charting.ColorLightPurple,
		charting.ColorIndigo,
		charting.ColorCrimson,
		charting.ColorYellow,
		charting.ColorLime,
		charting.ColorPink,
		charting.ColorFuchsia,
	}
	for cluster := range len(centroids) {
		cluster_points := make([]charting.DataPoint, 0)
		for i := range points {
			if labels[i] == cluster {
				cluster_points = append(cluster_points, points[i])
			}
		}

		key := fmt.Sprintf("cluster-%d", cluster)
		chart.Datasets[key] = &charting.ChartDataset{
			Label:           fmt.Sprintf("Cluster %d", cluster),
			BorderColor:     colors[cluster%len(colors)],
			BackgroundColor: []string{charting.ColorTransparent},
			PointRadius:     4,
			BorderWidth:     2,
			PointData:       cluster_points,
		}
	}

	centroidLabels := make([]string, len(centroids))
	for i := range centroids {
		centroidLabels[i] = fmt.Sprintf("Cluster %d", i)
	}
	chart.Datasets["centroids"] = &charting.ChartDataset{
		Label:           "Centroids",
		BorderColor:     "#000000",
		BackgroundColor: []string{"#ffffff"},
		PointRadius:     12,
		BorderWidth:     3,
		PointStyle:      "star",
		PointData:       centroids,
		PointLabels:     centroidLabels,
	}
}

func init() {
	SimpleChart.RenderFunc = RenderSimple
	KmeansChart.RenderFunc = RenderKmeans
}
