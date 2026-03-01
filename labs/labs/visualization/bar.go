package visualization

import (
	"encoding/csv"
	"fmt"
	"labs/labs/common"
	"labs/uncsv"
	"os"
)

const (
	BarChartID = "bar"

	BarGraphID = "orig-data"
)

var (
	BarGraph = common.ChartDataset{
		Label: "Bar Representation",
		BackgroundColor: []string{
			common.Color1,
			common.Color2,
			common.Color3,
			common.Color4,
			common.Color5,
		},
		BorderColor: "rgba(0, 0, 0, 0.1)",
		BorderWidth: 2,
		PointRadius: 0,
		ShowLine:    true,
		Togglable:   true,
	}

	BarChart = common.Chart{
		ID:          BarChartID,
		Title:       "Bar Plot",
		Type:        common.ChartTypeBar,
		XAxisLabel:  "Spending Type",
		YAxisLabel:  "Amount Spent",
		XAxisConfig: common.CategoryAxis,
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			BarGraphID: &BarGraph,
		},
	}

	BarMeta = BarChart.Meta()
)

type Spending struct {
	Category []string  `csv:"Категорія витрат"`
	Sum      []float64 `csv:"Сума (грн)"`
}

func RenderBarPlot(req *common.RenderRequest) (res *common.RenderResponse) {
	fmt.Printf("Rendering %s\n", req.ChartID)

	f, err := os.Open("../data/lab_4_var_12_spending.csv")
	if err != nil {
		return res.NewErrorf("error while opening file: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	d := uncsv.NewDecoder(r)

	spending := &Spending{}
	err = d.Decode(spending)
	if err != nil {
		return res.NewErrorf("error decoding csv: %v", err)
	}

	chartCopy := common.CopyChart(BarChart)

	err = chartCopy.UpdateDataForDataset(BarGraphID, spending.Sum)
	if err != nil {
		return res.NewErrorf("encountered error while updating points: %v", err)
	}

	chartCopy.Labels = spending.Category

	res = common.NewRenderResponse()
	res.AddChart(BarChartID, &chartCopy)

	return res
}
