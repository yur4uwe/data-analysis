package stats

import (
	"labs/labs/common"
	"math/rand"
	"os"
)

const (
	LabID = "5"
)

var (
	Config = common.LabConfig{
		Lab: Metadata,
		Charts: map[string]common.Chart{
			ErrorAnalysisChartID: ErrorAnalysisChart,
		},
	}

	Metadata = common.LabMetadata{
		ID:   LabID,
		Name: "Statistical Analysis",
		Charts: map[string]common.ChartMetadata{
			ErrorAnalysisChartID: ErrorAnalysisMeta,
		},
	}
)

// ...existing code...

func RandomSlice(stdev, mean float64, size int) []float64 {
	slice := make([]float64, size)

	for i := range size {
		slice[i] = stdev*rand.NormFloat64() + mean
	}

	return slice
}

func ReadSalaryCSV(filename string) ([]float64, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return nil, nil
}
