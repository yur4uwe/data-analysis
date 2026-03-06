package stats

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"math"
	"os"
)

const (
	TesterSalaryBarChartID = "tester-salary"
	TesterSalaryBarGraphID = "tester-salary"

	TesterSalaryStatsID = "tester-salary-stats"
)

var (
	TesterSalaryStats = charting.MutableField{
		ID:      TesterSalaryStatsID,
		Control: charting.ControlNoControl,
		Label:   "Tester salary statistics",
	}

	TesterSalaryGraph = charting.ChartDataset{
		Label:           "Tester Salary",
		BackgroundColor: []string{charting.Color1, charting.Color2, charting.Color3, charting.Color4, charting.Color5},
		PointRadius:     0,
		ShowLine:        true,
	}

	TesterSalaryChart = charting.Chart{
		ID:          TesterSalaryBarChartID,
		Title:       "Tester Salary Distribution",
		Type:        charting.ChartTypeBar,
		XAxisLabel:  "amount, $",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "people, n",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			TesterSalaryBarGraphID: &TesterSalaryGraph,
		},
		ChartVariables: []charting.MutableField{
			TesterSalaryStats,
		},
	}
)

func RenderTesterSalary(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if salaryRecords == nil {
		f, err := os.Open("../data/lab_5_var_12.csv")
		if err != nil {
			return res.NewErrorf("tester salary chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ';'
		salaryRecords = &SalaryRecord{}
		if err := d.Decode(salaryRecords); err != nil {
			return res.NewErrorf("tester salary chart: error while decoding csv: %s", err.Error())
		}
	}

	salaries := salariesFor(Tester)
	buckets := make([]float64, 5)
	min_salary := math.Inf(1)
	max_salary := math.Inf(-1)

	// Find min and max salary across all positions
	for _, salary := range salaries {
		min_salary = math.Min(min_salary, salary)
		max_salary = math.Max(max_salary, salary)
	}

	bucket_size := (max_salary - min_salary) / float64(len(buckets))
	for _, salary := range salaries {
		bucket_index := int((salary - min_salary) / bucket_size)
		if bucket_index >= len(buckets) {
			bucket_index = len(buckets) - 1
		}
		buckets[bucket_index]++
	}

	x := make([]float64, len(buckets))
	for i := range buckets {
		x[i] = min_salary + bucket_size*float64(i+1)
	}

	copyChart := charting.CopyChart(TesterSalaryChart)
	copyChart.UpdateDataForDataset(TesterSalaryBarGraphID, buckets)

	copyChart.Labels = make([]string, len(buckets))
	for i := range buckets {
		copyChart.Labels[i] = fmt.Sprintf("%.0f-%.0f", x[i], x[i]+bucket_size)
	}

	avg := CalculateMean(salaries)
	median := CalculateMedian(salaries)
	stddev := CalculateStdDev(salaries, avg)
	variance := CalculateVariance(salaries, avg)
	minSalary := math.Inf(1)
	for _, s := range salaries {
		minSalary = math.Min(minSalary, s)
	}

	copyChart.ChartVariables[0].Label = fmt.Sprintf(
		"Tester salary statistics:\nAverage = %.2f\nMedian = %.2f\nDeviation = %.2f\nVariance = %.2f\n",
		avg, median, stddev, variance,
	)

	res = charting.NewRenderResponse()
	res.AddChart(TesterSalaryBarChartID, &copyChart)
	return res
}
