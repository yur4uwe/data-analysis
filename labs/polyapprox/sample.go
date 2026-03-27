package polyapprox

import (
	"encoding/csv"
	"fmt"
	"io"
	"labs/charting"
	"labs/uncsv"
	"math"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	SampleDataID = "sample-data"

	sampleDataGraphID          = "sample"
	sampleApproximationGraphID = "sample-approx"

	approximationDegree = "degree"
	coeffsDisplayID     = "coeffs"
)

var (
	approxDegreeVariable = charting.MutableField{
		ID:      approximationDegree,
		Label:   "Degree of Polynomial",
		Default: 1,
		Min:     0,
		Max:     10,
		Step:    1,
		Control: charting.ControlRange,
	}

	coeffsDisplayVariable = charting.MutableField{
		ID:      coeffsDisplayID,
		Label:   "Polynomial coefficients: ",
		Control: charting.ControlNoControl,
	}

	sampleDataGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Sample Data",
			Type:        charting.ChartTypeScatter,
			BorderColor: charting.ToColor(charting.ColorEmerald),
			BorderWidth: 2,
			Togglable:   true,
		},
		BackgroundColor: charting.ToColor("rgba(0, 0, 0, 0.1)"),
		PointRadius:     3,
		HideLine:        true,
	}

	sampleDataApproxGraph = charting.GridDataset{
		BaseDataset: charting.BaseDataset{
			Label:       "Sample Data Approximation",
			BorderColor: charting.ToColor(charting.ColorAmber),
			BorderWidth: 2,
			Togglable:   true,
			GraphVariables: []charting.MutableField{
				approxDegreeVariable,
				coeffsDisplayVariable,
			},
		},
		BackgroundColor: charting.ToColor("rgba(0, 0, 0, 0.1)"),
		PointRadius:     0,
	}

	SampleDataChart = charting.Chart{
		ID:          SampleDataID,
		Title:       "Sample Data (CSV)",
		Type:        charting.ChartTypeScatter,
		XAxisLabel:  "X",
		YAxisLabel:  "Y",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			OriginalDataID:             &sampleDataGraph,
			sampleApproximationGraphID: &sampleDataApproxGraph,
		},
		ChartVariables: []charting.MutableField{},
	}

	SampleDataMetadata = SampleDataChart.Meta()

	points = (*Points)(nil)
)

func sortXandY(x, y []float64) {
	slices.Sort(x)

	sort.SliceStable(y, func(i, j int) bool {
		return x[i] < x[j]
	})
}

type Points struct {
	X []float64 `csv:"x"`
	Y []float64 `csv:"y_noisy"`
}

func RenderSampleData(req *charting.RenderRequest) (res *charting.RenderResponse) {
	if points == nil {
		f, err := os.Open("./data/lab_3_var_12.csv")
		if err != nil {
			fmt.Println("failed to open file:", err)
			return res.NewError("failed to read sample data file")
		}
		defer f.Close()

		d := uncsv.NewDecoder(f)
		d.Comma = ','
		points = &Points{}
		if err := d.Decode(points); err != nil {
			fmt.Println("failed to decode csv:", err)
			return res.NewError("failed to decode sample data file")
		}
	}

	chartCopy := charting.CopyChart(SampleDataChart)
	chartCopy.UpdatePointsForDataset(OriginalDataID, points.X, points.Y)

	degree, ok := req.GetGraphVariable(SampleDataID, sampleApproximationGraphID, approximationDegree)
	if !ok {
		degree = 2.0
	}

	coeffs, err := SolvePolynomialFit(points.X, points.Y, int(degree))
	if err != nil {
		return res.NewErrorf("failed to solve polynomial fit: %v", err)
	}

	mse := CalculateMSE(points.X, points.Y, coeffs)
	fmt.Printf("Sample Data Fit (Degree %d) MSE: %.4e\n", int(degree), mse)

	minX, maxX := math.Inf(1), math.Inf(-1)
	for _, xi := range points.X {
		maxX = max(maxX, xi)
		minX = min(minX, xi)
	}

	step := (maxX - minX) / float64(len(points.X)-1)

	approx := make([]float64, 0, len(points.X))
	appeoxX := make([]float64, 0, len(points.X))
	for i := minX; i < maxX; i += step {
		approx = append(approx, EvaluatePolynomial(coeffs, i))
		appeoxX = append(appeoxX, i)
	}

	chartCopy.UpdatePointsForDataset(sampleApproximationGraphID, appeoxX, approx)

	var str strings.Builder
	str.WriteString(fmt.Sprintf("MSE: %.4e, Coefficients (", mse))
	for i, c := range coeffs {
		fmt.Fprintf(&str, "x%d=%.2f", i, c)
		if i != len(coeffs)-1 {
			fmt.Fprint(&str, ", ")
		}
	}
	str.WriteString(")")
	chartCopy.Datasets[sampleApproximationGraphID].UpdateVariableLabel(1, str.String())

	return &charting.RenderResponse{
		Charts: map[string]charting.Chart{
			SampleDataID: chartCopy,
		},
	}
}
func ReadSampleCSV(filename string) ([]float64, []float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Read() // Skip header
	reader.Comma = ','

	var xVals, yVals []float64
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		x, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, nil, err
		}
		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, nil, err
		}

		xVals = append(xVals, x)
		yVals = append(yVals, y)
	}

	return xVals, yVals, nil
}
