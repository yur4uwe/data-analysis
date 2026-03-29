package visualization

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	LinearChartID = "linear"

	LinearGraphID = "orig-data" // due to the default, baseline graph on frontend, one graph from chart should always be called orig-data
)

var (
	LinearGraph = charting.CategoricalDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Revenue $",
			BorderColor: charting.ColorAmber,
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: []charting.Color{"rgba(37, 99, 235, 0.1)"},
	}

	LinearChart = charting.Chart{
		ID:          LinearChartID,
		Title:       "Revenue Per Day",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Day №",
		YAxisLabel:  "Revenue $",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			LinearGraphID: &LinearGraph,
		},
	}
)

type DailyRevenue struct {
	Day     []string  `csv:"День"`
	Revenue []float64 `csv:"Прибуток (грн)"`
}

func RenderLinear(req *charting.RenderRequest) (res *charting.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)
	f, err := os.Open("./data/lab_4_var_12_revenue_per_day.csv")
	if err != nil {
		return res.NewErrorf("error while opening file: %v", err)
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.SetComma(',')
	rec := &DailyRevenue{}
	if err := d.Decode(rec); err != nil {
		return res.NewErrorf("error while decoding csv: %v", err)
	}

	chartCopy := charting.CopyChart(LinearChart)
	if err := chartCopy.UpdateDataForDataset(LinearGraphID, charting.ToAnySlice(rec.Revenue)); err != nil {
		return res.NewErrorf("error while updating chart dataset: %v", err)
	}

	chartCopy.Labels = rec.Day

	res = charting.NewRenderResponse()
	res.AddChart(LinearChartID, &chartCopy)
	return res
}
