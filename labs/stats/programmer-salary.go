package stats

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"math"
	"os"
)

const (
	ProgrammerSalaryBarChartID   = "programmer-salary"
	ProgrammerSalaryBarGraphID   = "programmer-salary"
	EmpyricalDestributionChartID = "distribution"

	DisplayProgrammerSalaryStatsID = "programmer-salary-stats"
)

var (
	ProgrammerSalaryStats = charting.MutableField{
		ID:      DisplayProgrammerSalaryStatsID,
		Control: charting.ControlNoControl,
		Label:   "Programmer salary statistics",
	}

	ProgrammerSalaryGraph = charting.ChartDataset{
		Label:           "Programmer Salary",
		BackgroundColor: []string{charting.ColorEmerald, charting.ColorLime, charting.ColorIndigo, charting.ColorSlate, charting.ColorFuchsia},
		PointRadius:     0,
		ShowLine:        true,
	}

	ProgrammerSalaryChart = charting.Chart{
		ID:          ProgrammerSalaryBarChartID,
		Title:       "Programmer Salary",
		Type:        charting.ChartTypeBar,
		XAxisLabel:  "amount, $",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "people, n",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]*charting.ChartDataset{
			ProgrammerSalaryBarGraphID: &ProgrammerSalaryGraph,
		},
		ChartVariables: []charting.MutableField{
			ProgrammerSalaryStats,
		},
	}

	ProgrammerSalaryMeta = ProgrammerSalaryChart.Meta()

	salaryRecords = (*SalaryRecord)(nil)
)

func RenderProgrammerSalary(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if salaryRecords == nil {
		f, err := os.Open("./data/lab_5_var_12.csv")
		if err != nil {
			return res.NewErrorf("programmer salary chart: error while reading file: %s", err.Error())
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ';'
		salaryRecords = &SalaryRecord{}
		if err := d.Decode(salaryRecords); err != nil {
			return res.NewErrorf("programmer salary chart: error while decoding csv: %s", err.Error())
		}
	}

	salaries := salariesFor(Programmer)
	numBuckets := int(math.Round(1 + 3.334*math.Log10(float64(len(salaries)))))
	buckets := make([]float64, numBuckets)
	min_salary := math.Inf(1)
	max_salary := math.Inf(-1)
	for i := range salaries {
		min_salary = math.Min(min_salary, salaries[i])
		max_salary = math.Max(max_salary, salaries[i])
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

	copyChart := charting.CopyChart(ProgrammerSalaryChart)
	copyChart.UpdateDataForDataset(ProgrammerSalaryBarGraphID, charting.ToAnySlice(buckets))

	copyChart.Labels = make([]string, len(buckets))
	for i := range buckets {
		copyChart.Labels[i] = fmt.Sprintf("%.0f-%.0f", x[i], x[i]+bucket_size)
	}

	avg := CalculateMean(salaries)
	median := CalculateMedian(salaries)
	stddev := CalculateStdDev(salaries, avg)
	variance := CalculateVariance(salaries, avg)

	copyChart.ChartVariables[0].Label = fmt.Sprintf(
		"Programmer salary statistics\nAverage: %.2f\nMedian: %.2f\nDeviation: %.2f\nVariance: %.2f\n",
		avg, median, stddev, variance,
	)

	res = charting.NewRenderResponse()
	res.AddChart(ProgrammerSalaryBarChartID, &copyChart)
	return res
}
