package visualization

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	BarChartID = "bar"

	BarGraphID = "orig-data"
)

var (
	BarGraph = charting.CategoricalDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Spending",
			BorderColor: "rgba(0, 0, 0, 0.1)",
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: []charting.Color{
			charting.ColorAmber,
			charting.ColorBlue,
			charting.ColorCyan,
			charting.ColorEmerald,
			charting.ColorFuchsia,
		},
	}

	BarChart = charting.Chart{
		ID:          BarChartID,
		Title:       "Spending By Category",
		Type:        charting.ChartTypeBar,
		XAxisLabel:  "Spending Type",
		YAxisLabel:  "Amount Spent",
		XAxisConfig: charting.CategoryAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			BarGraphID: &BarGraph,
		},
	}
)

type Spending struct {
	Category []string  `csv:"Категорія витрат"`
	Sum      []float64 `csv:"Сума (грн)"`
}

func RenderBarPlot(req *charting.RenderRequest) (res *charting.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)

	f, err := os.Open("./data/lab_4_var_12_spending.csv")
	if err != nil {
		return res.NewErrorf("error while opening file: %v", err)
	}
	defer f.Close()

	d := uncsv.NewDecoder(f)
	d.SetComma(',')
	spending := &Spending{}
	err = d.Decode(spending)
	if err != nil {
		return res.NewErrorf("error decoding csv: %v", err)
	}

	chartCopy := charting.CopyChart(BarChart)

	chartCopy.UpdateDataForDataset(BarGraphID, charting.F64ToAny(spending.Sum))

	chartCopy.Labels = spending.Category

	res = charting.NewRenderResponse()
	res.AddChart(BarChartID, &chartCopy)

	return res
}
