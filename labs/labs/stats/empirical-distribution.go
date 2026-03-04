package stats

import (
	"encoding/csv"
	"labs/labs/common"
	"labs/uncsv"
	"os"
	"sort"
)

const (
	EmpiricalDistributionChartID = "empirical-distribution"
	EmpiricalDistributionGraphID = "empirical-distribution"
)

var (
	EmpiricalDistributionGraph = common.ChartDataset{
		Label:           "Empirical Distribution Function F(x)",
		BorderColor:     common.Color2,
		BackgroundColor: []string{common.ColorTransparent},
		ShowLine:        true,
		PointRadius:     3,
		BorderWidth:     2,
		Togglable:       false,
	}

	EmpiricalDistributionChart = common.Chart{
		ID:          EmpiricalDistributionChartID,
		Title:       "Empirical Distribution Function of Salaries",
		Type:        common.ChartTypeLine,
		XAxisLabel:  "Salary (USD)",
		XAxisConfig: common.LinearAxis,
		YAxisLabel:  "F(x) - Cumulative Probability",
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			EmpiricalDistributionGraphID: &EmpiricalDistributionGraph,
		},
	}
)

func RenderEmpiricalDistribution(req *common.RenderRequest) (res *common.RenderResponse) {
	// Load data if not already loaded
	if salaryRecords == nil {
		f, err := os.Open("../data/lab_5_var_12.csv")
		if err != nil {
			return res.NewErrorf("empirical distribution chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		r := csv.NewReader(f)
		r.Comma = ';'
		d := uncsv.NewDecoder(r)
		salaryRecords = &SalaryRecord{}
		if err := d.Decode(salaryRecords); err != nil {
			return res.NewErrorf("empirical distribution chart: error while decoding csv: %s", err.Error())
		}
	}

	// Get all salary values
	salaries := make([]float64, len(salaryRecords.Salary))
	copy(salaries, salaryRecords.Salary)

	// Sort salaries for EDF calculation
	sort.Float64s(salaries)

	n := float64(len(salaries))

	// Create points for the empirical distribution function
	// F(x) = n_x / n, where n_x is the number of values less than x
	x := make([]float64, 0)
	y := make([]float64, 0)

	// Add starting point (minimum value has F(x) = 0 just before it)
	if len(salaries) > 0 {
		x = append(x, salaries[0]-1)
		y = append(y, 0)
	}

	// For each salary value, calculate F(x)
	// Since salaries are sorted, the index+1 gives us the count of values <= current value
	for i, salary := range salaries {
		// F(x) = (number of values <= x) / total count
		fx := float64(i+1) / n
		x = append(x, salary)
		y = append(y, fx)
	}

	copyChart := common.CopyChart(EmpiricalDistributionChart)
	if err := copyChart.UpdatePointsForDataset(EmpiricalDistributionGraphID, x, y); err != nil {
		return res.NewErrorf("error updating dataset: %s", err.Error())
	}

	res = common.NewRenderResponse()
	res.AddChart(EmpiricalDistributionChartID, &copyChart)
	return res
}
