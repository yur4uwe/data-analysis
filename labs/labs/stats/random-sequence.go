package stats

import "labs/charting"

const (
	RandomSequenceChartID = "rand-sequence"

	RandomSequenceGraphID = "rand-sequence"
)

var (
	RandomSequenceChart = charting.Chart{
		ID:          RandomSequenceChartID,
		Title:       "Random Number Sequence",
		XAxisLabel:  "Index",
		YAxisLabel:  "Number",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Type:        charting.ChartTypeScatter,
		Datasets: map[string]*charting.ChartDataset{
			RandomSequenceGraphID: &RandomSequenceGraph,
		},
	}

	RandomSequenceMeta = RandomSequenceChart.Meta()

	RandomSequenceGraph = charting.ChartDataset{
		Label:           "Sequence of random numbers",
		BorderColor:     charting.Color6,
		BackgroundColor: []string{charting.Color5},
		PointRadius:     3,
		Togglable:       false,
	}
)

func RenderRandomSequence(req *charting.RenderRequest) (res *charting.RenderResponse) {
	sequence := GenerateNormalSamples(0, 0.2, 50_000)
	x := make([]float64, 0, len(sequence))
	for i := range sequence {
		x = append(x, float64(i))
	}

	chartCopy := charting.CopyChart(RandomSequenceChart)
	if err := chartCopy.UpdatePointsForDataset(RandomSequenceGraphID, x, sequence); err != nil {
		return res.NewErrorf("error updating dataset: %s", err.Error())
	}

	res = charting.NewRenderResponse()
	res.AddChart(chartCopy.ID, &chartCopy)
	return res
}
