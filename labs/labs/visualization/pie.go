package visualization

import (
	"fmt"
	"labs/labs/common"
)

const (
	RadialChartID = "radial"

	RadialGraphID = "orig-data" // due to the default, baseline graph on frontend, one chart should be called orig-data
)

var (
	RadialGraph = common.ChartDataset{
		Label:       "Radial Representation",
		BorderColor: common.Color2,
		BackgroundColor: []string{
			"rgba(220, 38, 38, 0.1)",
			common.Color10,
			common.Color11,
			common.Color12,
			common.Color5,
		},
		BorderWidth: 2,
		PointRadius: 0,
		ShowLine:    true,
		Togglable:   true,
	}

	RadialChart = common.Chart{
		ID:          RadialChartID,
		Title:       "Radial Plot",
		Type:        common.ChartTypePie,
		XAxisLabel:  "Category",
		YAxisLabel:  "Amount",
		XAxisConfig: common.CategoryAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			RadialGraphID: &RadialGraph,
		},
	}

	RadialMeta = RadialChart.Meta()
)

func RenderRadialPlot(req *common.RenderRequest) (res *common.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)
	values, err := ReadCategoricalCSV("../data/lab_4_var_12.csv")
	if err != nil {
		return res.NewErrorf("encountered error while reading csv: %v", err)
	}

	y := make([]float64, 0, len(values))
	labels := make([]string, 0, len(values))
	for k, v := range values {
		y = append(y, v)
		labels = append(labels, k)
	}

	chartCopy := common.CopyChart(RadialChart)

	// Pie charts need simple data array, not point data
	err = chartCopy.UpdateDataForDataset(RadialGraphID, y)
	if err != nil {
		return res.NewErrorf("encountered error while updating data: %v", err)
	}

	chartCopy.Labels = labels

	res = common.NewRenderResponse()
	res.AddChart(RadialChartID, &chartCopy)

	return res
}
