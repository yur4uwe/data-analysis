package visualization

import (
	"encoding/csv"
	"fmt"
	"labs/labs/common"
	"labs/uncsv"
	"os"
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

type RevenueSources struct {
	Sources []string  `csv:"Джерело доходу"`
	Sum     []float64 `csv:"Сума (грн)"`
}

func RenderRadialPlot(req *common.RenderRequest) (res *common.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)
	f, err := os.Open("../data/lab_4_var_12_revenue_sources.csv")
	if err != nil {
		return res.NewErrorf("encountered error while opening file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	d := uncsv.NewDecoder(r)
	rec := &RevenueSources{}
	if err := d.Decode(rec); err != nil {
		return res.NewErrorf("encountered error while decoding csv: %v", err)
	}

	chartCopy := common.CopyChart(RadialChart)

	// Pie charts need simple data array, not point data
	err = chartCopy.UpdateDataForDataset(RadialGraphID, rec.Sum)
	if err != nil {
		return res.NewErrorf("encountered error while updating data: %v", err)
	}

	chartCopy.Labels = rec.Sources

	res = common.NewRenderResponse()
	res.AddChart(RadialChartID, &chartCopy)

	return res
}
