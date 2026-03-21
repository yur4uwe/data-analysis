package visualization

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	RadialChartID = "radial"

	RadialGraphID = "orig-data" // due to the default, baseline graph on frontend, one chart should be called orig-data
)

var (
	RadialGraph = charting.ChartDataset{
		Label: "Revenue Sources", BorderColor: charting.ColorTransparent,
		BackgroundColor: []string{charting.ColorAmber, charting.ColorBlue,
			charting.ColorCyan, charting.ColorEmerald, charting.ColorFuchsia},
		BorderWidth: 2,
		PointRadius: 0,
		ShowLine:    true,
		Togglable:   true,
	}

	RadialChart = charting.Chart{
		ID:          RadialChartID,
		Title:       "Revenue Sources",
		Type:        charting.ChartTypePie,
		XAxisLabel:  "Category",
		YAxisLabel:  "Amount",
		XAxisConfig: charting.CategoryAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			RadialGraphID: &RadialGraph,
		},
	}

	RadialMeta = RadialChart.Meta()
)

type RevenueSources struct {
	Sources []string  `csv:"Джерело доходу"`
	Sum     []float64 `csv:"Сума (грн)"`
}

func RenderRadialPlot(req *charting.RenderRequest) (res *charting.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)
	f, err := os.Open("./data/lab_4_var_12_revenue_sources.csv")
	if err != nil {
		return res.NewErrorf("encountered error while opening file: %v", err)
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.SetComma(';')
	rec := &RevenueSources{}
	if err := d.Decode(rec); err != nil {
		return res.NewErrorf("encountered error while decoding csv: %v", err)
	}

	chartCopy := charting.CopyChart(RadialChart)

	// Pie charts need simple data array, not point data
	err = chartCopy.UpdateDataForDataset(RadialGraphID, charting.ToAnySlice(rec.Sum))
	if err != nil {
		return res.NewErrorf("encountered error while updating data: %v", err)
	}

	chartCopy.Labels = rec.Sources

	res = charting.NewRenderResponse()
	res.AddChart(RadialChartID, &chartCopy)

	return res
}
