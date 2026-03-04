package stats

import (
	"fmt"
	"labs/labs/common"
	"labs/uncsv"
	"math"
	"math/rand/v2"
)

const (
	LabID = "5"
)

var (
	Config = common.NewLabConfig(
		LabID,
		"Statistical Analysis",
		map[string]*common.Chart{
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

// GenerateNormalSamples generates samples from normal distribution
func GenerateNormalSamples(mean, stddev float64, sampleCount int) []float64 {
	samples := make([]float64, sampleCount)
	for i := range sampleCount {
		samples[i] = rand.NormFloat64()*stddev + mean
	}
	return samples
}

// CalculateMean computes the mean of a slice
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

// CalculateStdDev computes the standard deviation of a slice
func CalculateStdDev(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range data {
		diff := v - mean
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(data))
	return math.Sqrt(variance)
}
