package stats

import "labs/labs/common"

const (
	RandomSequenceChartID = "rand-sequence"

	RandomSequenceGraphID = "rand-sequence"
)

var (
	RandomSequenceChart = common.Chart{
		ID:          RandomSequenceChartID,
		Title:       "Random Number Sequence",
		XAxisLabel:  "Index",
		YAxisLabel:  "Number",
		XAxisConfig: common.LinearAxis,
		YAxisConfig: common.LinearAxis,
		Type:        common.ChartTypeScatter,
		Datasets: map[string]*common.ChartDataset{
			RandomSequenceGraphID: &RandomSequenceGraph,
		},
	}

	RandomSequenceMeta = RandomSequenceChart.Meta()

	RandomSequenceGraph = common.ChartDataset{
		Label:           "Sequence of random numbers",
		BorderColor:     common.Color6,
		BackgroundColor: []string{common.Color5},
		PointRadius:     0,
		ShowLine:        true,
		BorderWidth:     0,
		Togglable:       false,
	}
)

func RenderRandomSequence(req *common.RenderRequest) (res *common.RenderResponse) {
	sequence := GenerateNormalSamples(0, 0.2, 50_000)
	x := make([]float64, 0, len(sequence))
	for i := range sequence {
		x = append(x, float64(i))
	}

	chartCopy := common.CopyChart(RandomSequenceChart)
	if err := chartCopy.UpdatePointsForDataset(RandomSequenceGraphID, x, sequence); err != nil {
		return res.NewErrorf("error updating dataset: %s", err.Error())
	}

	res = common.NewRenderResponse()
	res.AddChart(chartCopy.ID, &chartCopy)
	return res
}
