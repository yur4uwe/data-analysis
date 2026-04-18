package network

import (
	"labs/charting"
	"labs/uncsv"
	"os"
)

type ClusterizationPoint struct {
	X float64 `csv:"x"`
	Y float64 `csv:"y"`
}

const (
	LabID = "12"

	DataChartID     = "data"
	WeightsChartID  = "weights"
	ClustersChartID = "clusters"

	DataGraphID    = "data"
	WeightsGraphID = "weights"
	UnitCircleID   = "unit-circle"

	NumClustersID  = "numClusters"
	NumEpochsID    = "numEpochs"
	LearningRateID = "learningRate"
)

var (
	Config = charting.NewLabConfig(
		LabID,
		"Neural Network",
		map[string]*charting.Chart{
			DataChartID:     &DataChart,
			WeightsChartID:  &WeightsChart,
			ClustersChartID: &ClustersChart,
		},
	)

	VarNumClusters = charting.MutableField{
		ID:      NumClustersID,
		Label:   "Number of Clusters",
		Default: 3,
		Min:     1,
		Max:     10,
		Step:    1,
		Control: charting.ControlRange,
	}

	VarNumEpochs = charting.MutableField{
		ID:      NumEpochsID,
		Label:   "Number of Epochs",
		Default: 100,
		Min:     10,
		Max:     1000,
		Step:    10,
		Control: charting.ControlRange,
	}

	VarLearningRate = charting.MutableField{
		ID:      LearningRateID,
		Label:   "Learning Rate",
		Default: 0.1,
		Min:     0.001,
		Max:     1.0,
		Step:    0.01,
		Control: charting.ControlRange,
	}

	UnitCircleGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Unit Circle",
			BorderColor: charting.ColorAmber,
			BorderWidth: 2,
		},
		PointRadius:     2,
		BackgroundColor: charting.ColorAmber,
	}

	DataGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label: "Data",
		},
		PointRadius:     3,
		BackgroundColor: charting.ColorAmber,
	}

	DataChart = charting.Chart{
		ID:          "data",
		Title:       "Data",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "x",
		YAxisLabel:  "y",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			DataGraphID: &DataGraph,
		},
	}

	WeightsChart = charting.Chart{
		ID:           WeightsChartID,
		Title:        "Weights Movement",
		Type:         charting.ChartTypeScatter,
		XAxisLabel:   "w0",
		YAxisLabel:   "w1",
		XAxisConfig:  charting.LinearAxis,
		YAxisConfig:  charting.LinearAxis,
		SquareLayout: true,
		ChartVariables: []charting.MutableField{
			VarNumClusters,
			VarNumEpochs,
			VarLearningRate,
		},
		Datasets: map[string]charting.Dataset{
			UnitCircleID: &UnitCircleGraph,
		},
	}

	ClustersChart = charting.Chart{
		ID:          ClustersChartID,
		Title:       "Clusterized Data",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "x",
		YAxisLabel:  "y",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		ChartVariables: []charting.MutableField{
			VarNumClusters,
			VarNumEpochs,
			VarLearningRate,
		},
		Datasets: map[string]charting.Dataset{},
	}

	data = ([]ClusterizationPoint)(nil)
)

func init() {
	DataChart.RenderFunc = RenderData
	WeightsChart.RenderFunc = RenderWeights
	ClustersChart.RenderFunc = RenderClusters
}

func loadData() error {
	if data != nil {
		return nil
	}
	f, err := os.Open("data/lab_12_var_12.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	data = []ClusterizationPoint{}
	if err := uncsv.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	return nil
}

func RenderData(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadData(); err != nil {
		return res.NewErrorf("error loading data: %v", err)
	}
	chartCopy := charting.CopyChart(DataChart)
	points := make([]charting.DataPoint, len(data))
	for i, x := range data {
		points[i] = charting.DataPoint{
			X: x.X,
			Y: &x.Y,
		}
	}

	chartCopy.UpdateDataPointsForDataset(DataGraphID, points)

	res = charting.NewRenderResponse()
	res.AddChart(DataChartID, &chartCopy)
	return res
}
