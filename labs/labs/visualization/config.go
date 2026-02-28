package visualization

import (
	"encoding/csv"
	"errors"
	"io"
	"labs/labs/common"
	"os"
	"strconv"
)

const (
	LabID = "4"
)

var (
	Config = common.LabConfig{
		Lab: Metadata,
		Charts: map[string]common.Chart{
			BarChartID:      BarChart,
			FunctionChartID: FunctionChart,
			LinearChartID:   LinearChart,
			RadialChartID:   RadialChart,
		},
	}

	Metadata = common.LabMetadata{
		ID:   LabID,
		Name: "Visualization",
		Charts: map[string]common.ChartMetadata{
			BarChartID:      BarMeta,
			FunctionChartID: FunctionMeta,
			LinearChartID:   LinearMeta,
			RadialChartID:   RadialMeta,
		},
	}
)

func ReadCategoricalCSV(path string) (map[string]float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = ','

	reader.Read() // Skip Header

	categories := make(map[string]float64)

	for {
		records, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if len(records) < 2 {
			return nil, errors.New("Incorrect csv schema")
		}

		sum, err := strconv.ParseFloat(records[1], 64)
		if err != nil {
			return nil, err
		}

		categories[records[0]] = sum
	}

	return categories, nil
}
