package visualization

import (
	"labs/labs/common"
)

const (
	LabID = "4"
)

var (
	Config = common.LabConfig{
		Lab: Metadata,
		Charts: map[string]*common.Chart{
			BarChartID:      &BarChart,
			FunctionChartID: &FunctionChart,
			LinearChartID:   &LinearChart,
			RadialChartID:   &RadialChart,
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
