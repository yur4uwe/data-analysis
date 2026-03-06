package stats

import (
	"fmt"
	"labs/charting"
	"labs/uncsv"
	"math"
	"math/rand/v2"
	"sort"
)

const (
	LabID = "5"
)

var (
	Config = charting.NewLabConfig(
		LabID,
		"Statistical Analysis",
		map[string]*charting.Chart{
			RandomSequenceChartID:        &RandomSequenceChart,
			CorrelationChartID:           &CorrelationChart,
			ProgrammerSalaryBarChartID:   &ProgrammerSalaryChart,
			TesterSalaryBarChartID:       &TesterSalaryChart,
			EmpiricalDistributionChartID: &EmpiricalDistributionChart,
		},
	)

	Metadata = Config.Lab
)

func init() {
	CorrelationChart.RenderFunc = RenderError
	Config.Charts[RandomSequenceChartID].RenderFunc = RenderRandomSequence
	Config.Charts[ProgrammerSalaryBarChartID].RenderFunc = RenderProgrammerSalary
	Config.Charts[TesterSalaryBarChartID].RenderFunc = RenderTesterSalary
	Config.Charts[EmpiricalDistributionChartID].RenderFunc = RenderEmpiricalDistribution
}

type PositionType string

const (
	Programmer PositionType = "Програміст"
	Tester     PositionType = "Тестувальник"
)

var _ uncsv.FieldDecoder = (*PositionType)(nil)

func (p *PositionType) DecodeCSV(field string) error {
	switch field {
	case string(Programmer):
		*p = Programmer
	case string(Tester):
		*p = Tester
	default:
		return fmt.Errorf("invalid position: %s", field)
	}
	return nil
}

func IsValidPosition(pos string) bool {
	switch pos {
	case string(Programmer), string(Tester):
		return true
	default:
		return false
	}
}

type SalaryRecord struct {
	ID              []string       `csv:"Працівник"`
	Age             []int          `csv:"Вік"`
	Position        []PositionType `csv:"Посада"`
	ExperienceYears []int          `csv:"Досвід (років)"`
	Projects        []int          `csv:"Проектів завершено"`
	Salary          []float64      `csv:"Зарплата (USD)"`
}

func GenerateNormalSamples(mean, stddev float64, sampleCount int) []float64 {
	samples := make([]float64, sampleCount)
	for i := range sampleCount {
		samples[i] = rand.NormFloat64()*stddev + mean
	}
	return samples
}

func CalculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func CalculateStdDev(data []float64, mean float64) float64 {
	return math.Sqrt(CalculateVariance(data, mean))
}

func CalculateVariance(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range data {
		diff := v - mean
		sumSquares += diff * diff
	}
	return sumSquares / float64(len(data))
}

func CalculateMedian(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

func salariesFor(position PositionType) []float64 {
	if salaryRecords == nil {
		return nil
	}
	var result []float64
	for i, pos := range salaryRecords.Position {
		if pos == position {
			result = append(result, salaryRecords.Salary[i])
		}
	}
	return result
}
