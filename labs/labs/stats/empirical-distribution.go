package stats

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
	EmpiricalDistributionProgrammerGraph = charting.ChartDataset{
		Label:           "Programmer F(x)",
		BorderColor:     charting.Color2,
		BackgroundColor: []string{charting.ColorTransparent},
		ShowLine:        true,
		PointRadius:     3,
		BorderWidth:     2,
		Togglable:       true,
	}

	EmpiricalDistributionTesterGraph = charting.ChartDataset{
		Label:           "Tester F(x)",
		BorderColor:     charting.Color4,
		BackgroundColor: []string{charting.ColorTransparent},
		ShowLine:        true,
		PointRadius:     3,
		BorderWidth:     2,
		Togglable:       true,
	}

	EmpiricalDistributionChart = charting.Chart{
		ID:          EmpiricalDistributionChartID,
		Title:       "Empirical Distribution Function of Salaries",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "Salary (USD)",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "F(x) - Cumulative Probability",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
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
	if salaryRecords == nil {
		f, err := os.Open("../data/lab_5_var_12.csv")
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
	if err := copyChart.UpdatePointsForDataset(EmpiricalDistributionProgrammerID, px, py); err != nil {
		return res.NewErrorf("error updating programmer dataset: %s", err.Error())
	}

	tx, ty := buildEDF(salariesFor(Tester))
	if err := copyChart.UpdatePointsForDataset(EmpiricalDistributionTesterID, tx, ty); err != nil {
		return res.NewErrorf("error updating tester dataset: %s", err.Error())
	}

	res = charting.NewRenderResponse()
	res.AddChart(EmpiricalDistributionChartID, &copyChart)
	return res
}
