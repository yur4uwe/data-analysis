package neuron

import (
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	LabID = "11"

	// Chart IDs
	DataChartID        = "data"
	BoundaryChartID    = "boundary"
	ConvergenceChartID = "convergence"
	TrajectoryChartID  = "trajectory"

	// Variable IDs
	VarEpochsID           = "epochs"
	VarLRID               = "lr"
	VarAlphaID            = "alpha"
	VarHeatmapPrecisionID = "heatmap_precision"
	VarWeightMarginID     = "weight_margin"
	VarActivationID       = "activation"
	VarFormulaID          = "formula"

	// Graph IDs (Datasets)
	GraphClass0ID    = "c0"
	GraphClass1ID    = "c1"
	GraphBoundaryID  = "boundary"
	GraphHeatmapID   = "heatmap"
	GraphLossID      = "loss" // Keep for backward compat or remove
	GraphPathID      = "path"
	GraphLandscapeID = "landscape"

	// Specific Loss Graph IDs
	GraphLossSigmoidID = "loss-sigmoid"
	GraphLossTanhID    = "loss-tanh"
	GraphLossReLUID    = "loss-relu"
)

type ClassificationData struct {
	X     []float64 `csv:"x"`
	Y     []float64 `csv:"y"`
	Class []bool    `csv:"class"`
}

var (
	SharedVariables = []charting.MutableField{
		{ID: VarEpochsID, Label: "Epochs", Default: 100, Min: 10, Max: 1000, Step: 10, Control: charting.ControlRange},
		{ID: VarLRID, Label: "Learning Rate", Default: 0.1, Min: 0.001, Max: 1.0, Step: 0.01, Control: charting.ControlRange},
		{ID: VarAlphaID, Label: "Alpha", Default: 1.0, Min: 0.1, Max: 5.0, Step: 0.1, Control: charting.ControlRange},
	}

	DisplayFormula = charting.MutableField{
		ID:      VarFormulaID,
		Label:   "Formula",
		Control: charting.ControlNoControl,
	}

	ActivationFuncField = charting.MutableField{
		ID:      VarActivationID,
		Label:   "Activation",
		Default: 0,
		Control: charting.ControlSelect,
		Options: []string{"Sigmoid", "Tanh", "ReLU"},
	}

	DataChart = charting.Chart{
		ID:          DataChartID,
		Title:       "Data",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "x",
		XAxisConfig: charting.LinearAxis,
		YAxisLabel:  "y",
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			GraphClass0ID: &Class0Dataset,
			GraphClass1ID: &Class1Dataset,
		},
		ChartVariables: SharedVariables,
	}

	Class0Dataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Class 0",
			BorderColor: charting.ColorRed,
		},
		HideLine:    true,
		PointRadius: 4,
	}

	Class1Dataset = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Class 1",
			BorderColor: charting.ColorBlue,
		},
		HideLine:    true,
		PointRadius: 4,
	}

	Config = charting.NewLabConfig(
		LabID,
		"1 Neuron Neural Network",
		map[string]*charting.Chart{
			DataChartID:        &DataChart,
			BoundaryChartID:    &BoundaryChart,
			ConvergenceChartID: &ConvergenceChart,
			TrajectoryChartID:  &TrajectoryChart,
		},
	)

	trainData        = (*ClassificationData)(nil)
	lastTrainResults = make(map[int]TrainingResult)
)

func init() {
	DataChart.RenderFunc = RenderData
	BoundaryChart.RenderFunc = RenderBoundary
	ConvergenceChart.RenderFunc = RenderConvergence
	TrajectoryChart.RenderFunc = RenderTrajectory
}

func loadData() error {
	if trainData != nil {
		return nil
	}
	f, err := os.Open("data/lab_11_var_12.csv")
	if err != nil {
		return err
	}
	defer f.Close()
	trainData = &ClassificationData{}
	return uncsv.NewDecoder(f).Decode(trainData)
}

func splitData(data *ClassificationData) (c0, c1 []charting.DataPoint) {
	for i := range data.X {
		y := data.Y[i]
		p := charting.DataPoint{X: data.X[i], Y: &y}
		if data.Class[i] {
			c1 = append(c1, p)
		} else {
			c0 = append(c0, p)
		}
	}
	return
}

func getActivation(actIdx int, alpha float64) (func(float64) float64, func(float64) float64) {
	switch actIdx {
	case 1:
		return newTanh(alpha)
	case 2:
		return newReLU(alpha)
	default:
		return newSigmoid(alpha)
	}
}

func ensureTrained(req *charting.RenderRequest) error {
	if err := loadData(); err != nil {
		return err
	}

	epochs := uint32(req.GetVariable(VarEpochsID))
	lr := req.GetVariable(VarLRID)
	alpha := req.GetVariable(VarAlphaID)

	// Train all three activation functions
	for i := range 3 {
		act, actDeriv := getActivation(i, alpha)
		lastTrainResults[i] = train(trainData, epochs, lr, act, actDeriv)
	}
	return nil
}

func RenderData(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if err := loadData(); err != nil {
		return res.NewErrorf("error loading data: %v", err)
	}

	res = charting.NewRenderResponse()
	dataChartCopy := charting.CopyChart(DataChart)

	c0, c1 := splitData(trainData)
	dataChartCopy.UpdateDataPointsForDataset(GraphClass0ID, c0)
	dataChartCopy.UpdateDataPointsForDataset(GraphClass1ID, c1)

	res.AddChart(DataChartID, &dataChartCopy)
	return res
}
