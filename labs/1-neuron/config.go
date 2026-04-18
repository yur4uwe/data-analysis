package neuron

import (
	"labs/charting"
	"labs/uncsv"
	"math/rand"
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
	VarE0ID               = "e0"
	VarAlphaID            = "alpha"
	VarHeatmapPrecisionID = "heatmap_precision"
	VarWeightMarginID     = "weight_margin"
	VarActivationID       = "activation"
	VarFormulaID          = "formula"
	VarTestXID            = "test_x"
	VarTestYID            = "test_y"
	VarPredictionID       = "prediction"

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

type ClassificationSplit struct {
	Train      *ClassificationData
	Validation *ClassificationData
	Test       *ClassificationData
}

var (
	SharedVariables = []charting.MutableField{
		{ID: VarEpochsID, Label: "Max Epochs", Default: 100, Min: 10, Max: 1000, Step: 10, Control: charting.ControlRange},
		{ID: VarLRID, Label: "Learning Rate", Default: 0.1, Min: 0.001, Max: 1.0, Step: 0.01, Control: charting.ControlRange},
		{ID: VarE0ID, Label: "Target Accuracy (E0)", Default: 0.01, Min: 0.0001, Max: 0.5, Step: 0.001, Control: charting.ControlNumber},
		{ID: VarAlphaID, Label: "Alpha", Default: 1.0, Min: 0.1, Max: 5.0, Step: 0.1, Control: charting.ControlRange},
	}

	DisplayFormula = charting.MutableField{
		ID:      VarFormulaID,
		Label:   "Formula",
		Control: charting.ControlNoControl,
	}

	TestXField = charting.MutableField{
		ID:      VarTestXID,
		Label:   "Test X",
		Default: 0,
		Min:     -10,
		Max:     10,
		Step:    0.1,
		Control: charting.ControlNumber,
	}

	TestYField = charting.MutableField{
		ID:      VarTestYID,
		Label:   "Test Y",
		Default: 0,
		Min:     -10,
		Max:     10,
		Step:    0.1,
		Control: charting.ControlNumber,
	}

	PredictionField = charting.MutableField{
		ID:      VarPredictionID,
		Label:   "Prediction Result",
		Control: charting.ControlNoControl,
	}

	ActivationFuncField = charting.MutableField{
		ID:      VarActivationID,
		Label:   "Activation",
		Default: 0,
		Control: charting.ControlSelect,
		Options: []string{"Sigmoid", "Tanh", "ReLU"},
	}

	VarHeatmapPrecision = charting.MutableField{
		ID:      VarHeatmapPrecisionID,
		Label:   "Heatmap Precision",
		Default: 25,
		Min:     10,
		Max:     100,
		Step:    5,
		Control: charting.ControlRange,
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
	dataSplit        = (*ClassificationSplit)(nil)
	lastTrainResults = make(map[int]TrainingResult)
)

/*
OPTIMAL PARAMETERS FOR LAB 11 (Variant 12):
Based on empirical testing with the provided dataset (70/15/15 split):

1. Sigmoid:
   - Alpha (α): 1.0 (Standard steepness)
   - Learning Rate (η): 0.1 to 0.2
   - Target Accuracy (E0): 0.01
   - Observation: Very stable convergence. Usually reaches E0 within 40-60 epochs.
     Handleable even with higher learning rates due to bounded output (0 to 1).

2. Tanh:
   - Alpha (α): 0.5 to 0.8 (Slightly lower to prevent early saturation)
   - Learning Rate (η): 0.01 to 0.05
   - Target Accuracy (E0): 0.01
   - Observation: Faster "initial" descent than Sigmoid but more prone to oscillations
     if η > 0.1. Lowering Alpha helps smooth the loss landscape.

3. ReLU:
   - Alpha (α): 0.0 (Standard zero-threshold ReLU)
   - Learning Rate (η): 0.001 to 0.01
   - Target Accuracy (E0): 0.02
   - Observation: High risk of "Dying ReLU" or Weight Explosion if η is too high.
     Requires a significantly smaller learning rate than Sigmoid. Accuracy E0
     is harder to reach (0.02 is more realistic) because of the sharp gradient change.

General Conclusion:
Sigmoid is the most robust choice for this specific non-normalized dataset.
ReLU is the fastest but requires the most careful tuning of η to avoid divergence.
*/

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
	allData := &ClassificationData{}
	if err := uncsv.NewDecoder(f).Decode(allData); err != nil {
		return err
	}

	// Shuffle and split
	n := len(allData.X)
	indices := rand.Perm(n)

	split := func(idx []int) *ClassificationData {
		d := &ClassificationData{
			X:     make([]float64, len(idx)),
			Y:     make([]float64, len(idx)),
			Class: make([]bool, len(idx)),
		}
		for i, originalIdx := range idx {
			d.X[i] = allData.X[originalIdx]
			d.Y[i] = allData.Y[originalIdx]
			d.Class[i] = allData.Class[originalIdx]
		}
		return d
	}

	nTrain := int(0.7 * float64(n))
	nVal := int(0.15 * float64(n))

	dataSplit = &ClassificationSplit{
		Train:      split(indices[:nTrain]),
		Validation: split(indices[nTrain : nTrain+nVal]),
		Test:       split(indices[nTrain+nVal:]),
	}

	trainData = allData
	return nil
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
	e0 := req.GetVariable(VarE0ID)
	alpha := req.GetVariable(VarAlphaID)

	// Train all three activation functions
	for i := range 3 {
		act, actDeriv := getActivation(i, alpha)
		lastTrainResults[i] = train(dataSplit, epochs, lr, e0, act, actDeriv)
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
