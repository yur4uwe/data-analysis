package stats

import (
	"encoding/csv"
	"fmt"
	"labs/labs/common"
	"labs/uncsv"
	"math"
	"os"
)

const (
	TesterSalaryBarChartID = "tester-salary"
	TesterSalaryBarGraphID = "tester-salary"
)

var (
	TesterSalaryGraph = common.ChartDataset{
		Label:           "Tester Salary",
		BackgroundColor: []string{common.Color1, common.Color2, common.Color3, common.Color4, common.Color5},
		PointRadius:     0,
		ShowLine:        true,
	}

	TesterSalaryChart = common.Chart{
		ID:          TesterSalaryBarChartID,
		Title:       "Tester Salary Distribution",
		Type:        common.ChartTypeBar,
		XAxisLabel:  "amount, $",
		XAxisConfig: common.LinearAxis,
		YAxisLabel:  "people, n",
		YAxisConfig: common.LinearAxis,
		Datasets: map[string]*common.ChartDataset{
			TesterSalaryBarGraphID: &TesterSalaryGraph,
		},
	}
)

func RenderTesterSalary(req *common.RenderRequest) (res *common.RenderResponse) {
	if salaryRecords == nil {
		f, err := os.Open("../data/lab_5_var_12.csv")
		if err != nil {
			return res.NewErrorf("tester salary chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		r := csv.NewReader(f)
		r.Comma = ';'
		d := uncsv.NewDecoder(r)
		salaryRecords = &SalaryRecord{}
		if err := d.Decode(salaryRecords); err != nil {
			return res.NewErrorf("tester salary chart: error while decoding csv: %s", err.Error())
		}
	}

	buckets := make([]float64, 5)
	min_salary := math.Inf(1)
	max_salary := math.Inf(-1)

	// Find min and max salary across all positions
	for i := range salaryRecords.ID {
		min_salary = math.Min(min_salary, salaryRecords.Salary[i])
		max_salary = math.Max(max_salary, salaryRecords.Salary[i])
	}

	bucket_size := (max_salary - min_salary) / float64(len(buckets))

	// Fill buckets only for testers
	for i := range salaryRecords.ID {
		if salaryRecords.Position[i] != Tester {
			continue
		}

		bucket_index := int((salaryRecords.Salary[i] - min_salary) / bucket_size)
		if bucket_index >= len(buckets) {
			bucket_index = len(buckets) - 1
		}
		buckets[bucket_index]++
	}

	x := make([]float64, len(buckets))
	for i := range buckets {
		x[i] = min_salary + bucket_size*float64(i+1)
	}

	copyChart := common.CopyChart(TesterSalaryChart)
	copyChart.UpdateDataForDataset(TesterSalaryBarGraphID, buckets)

	copyChart.Labels = make([]string, len(buckets))
	for i := range buckets {
		copyChart.Labels[i] = fmt.Sprintf("%.0f-%.0f", x[i], x[i]+bucket_size)
	}

	res = common.NewRenderResponse()
	res.AddChart(TesterSalaryBarChartID, &copyChart)
	return res
}
