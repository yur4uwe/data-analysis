package statslab

import (
	"labs/charting"
	"labs/uncsv"
	"os"
	"sort"
)

const (
	EmpiricalDistributionChartID      = "empirical-distribution"
	EmpiricalDistributionProgrammerID = "empirical-distribution-programmer"
	EmpiricalDistributionTesterID     = "empirical-distribution-tester"
)

var (
	EmpiricalDistributionProgrammerGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Programmer F(x)",
			BorderColor: charting.ColorEmerald,
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     3,
	}

	EmpiricalDistributionTesterGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Tester F(x)",
			BorderColor: charting.ColorLime,
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ColorTransparent,
		PointRadius:     3,
	}

	EmpiricalDistributionChart = charting.Chart{
		ID:          EmpiricalDistributionChartID,
		Title:       "Empirical Distribution Function of Salaries",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "Salary (USD)",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "F(x) - Cumulative Probability",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			EmpiricalDistributionProgrammerID: &EmpiricalDistributionProgrammerGraph,
			EmpiricalDistributionTesterID:     &EmpiricalDistributionTesterGraph,
		},
	}
)

func buildEDF(salaries []float64) (x, y []float64) {
	sorted := make([]float64, len(salaries))
	copy(sorted, salaries)
	sort.Float64s(sorted)
	n := float64(len(sorted))

	if len(sorted) > 0 {
		x = append(x, sorted[0]-1)
		y = append(y, 0)
	}
	for i, v := range sorted {
		x = append(x, v)
		y = append(y, float64(i+1)/n)
	}
	return
}

func RenderEmpiricalDistribution(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if len(salaryRecords.ID) == 0 {
		f, err := os.Open("./data/lab_5_var_12.csv")
		if err != nil {
			return res.NewErrorf("empirical distribution chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ';'
		salaryRecords = &SalaryRecord{}
		if err := d.Decode(salaryRecords); err != nil {
			return res.NewErrorf("empirical distribution chart: error while decoding csv: %s", err.Error())
		}
	}

	copyChart := charting.CopyChart(EmpiricalDistributionChart)

	px, py := buildEDF(salariesFor(Programmer))
	copyChart.UpdatePointsForDataset(EmpiricalDistributionProgrammerID, px, py)

	tx, ty := buildEDF(salariesFor(Tester))
	copyChart.UpdatePointsForDataset(EmpiricalDistributionTesterID, tx, ty)

	res = charting.NewRenderResponse()
	res.AddChart(EmpiricalDistributionChartID, &copyChart)
	return res
}
