package visualization

import (
	"encoding/csv"
	"fmt"
	"labs/labs/common"
	"labs/uncsv"
	"os"
)

const (
	LinearChartID = "linear"

	LinearGraphID = "orig-data" // due to the default, baseline graph on frontend, one graph from chart should always be called orig-data
)

var (
	LinearGraph = common.ChartDataset{
		Label:           "Linear Data",
		BorderColor:     common.Color1,
		BackgroundColor: []string{"rgba(37, 99, 235, 0.1)"},
		PointRadius:     0,
		BorderWidth:     2,
		ShowLine:        true,
		Togglable:       true,
	}

	LinearChart = common.Chart{
		ID:          LinearChartID,
		Title:       "Linear Representation",
		Type:        common.ChartTypeLine,
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			LinearGraphID: &LinearGraph,
		},
	}

	LinearMeta = LinearChart.Meta()
)

type DailyRevenue struct {
	Day     []string  `csv:"День"`
	Revenue []float64 `csv:"Прибуток (грн)"`
}

func RenderLinear(req *common.RenderRequest) (res *common.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)
	f, err := os.Open("../data/lab_4_var_12_revenue_per_day.csv")
	if err != nil {
		return res.NewErrorf("error while opening file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ','
	d := uncsv.NewDecoder(r)
	rec := &DailyRevenue{}
	if err := d.Decode(rec); err != nil {
		return res.NewErrorf("error while decoding csv: %v", err)
	}

	chartCopy := common.CopyChart(LinearChart)
	if err := chartCopy.UpdateDataForDataset(LinearGraphID, rec.Revenue); err != nil {
		return res.NewErrorf("error while updating chart dataset: %v", err)
	}

	chartCopy.Labels = rec.Day

	res = common.NewRenderResponse()
	res.AddChart(LinearChartID, &chartCopy)
	return res
}
