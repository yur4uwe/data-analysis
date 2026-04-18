package polyapprox

import (
	"fmt"
	"labs/analysis"
	"labs/charting"
	"labs/uncsv"
	"os"
)

const (
	SampleMSEID     = "sample-mse"
	SampleMSEDataID = "mse-data"
)

var (
	SampleMSEChart = charting.Chart{
		ID:          SampleMSEID,
		Title:       "MSE vs Degree (CSV)",
		Type:        charting.ChartTypeLine,
		XAxisLabel:  "Polynomial Degree",
		YAxisLabel:  "Mean Squared Error",
		XAxisConfig: charting.LinearAxis,
		YAxisConfig: charting.LinearAxis,
		Datasets: map[string]charting.Dataset{
			OriginalDataID: &charting.GridDataset{
				BaseDataset: charting.BaseDataset{
					Label:          "MSE vs Degree",
					BorderColor:    charting.ColorAmber,
					BorderWidth:    2,
					GraphVariables: []charting.MutableField{BestDegreeField},
				},
			},
		},
		ChartVariables: ChartVariables,
	}

	SampleMSEMetadata = SampleMSEChart.Meta()
)

func RenderSamplePolynomialMSE(req *charting.RenderRequest) (res *charting.RenderResponse) {
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

	maxDegree := min(len(points.X)-1, 45)
	degrees := make([]float64, 0, maxDegree)
	errs := make([]float64, 0, maxDegree)

	bestDegree := 1
	minMSE := -1.0

	for degree := range maxDegree - 1 {
		degree += 1
		coeffs, err := analysis.SolvePolynomialFit(points.X, points.Y, degree)
		if err != nil {
			continue
		}
		predicted := make([]float64, len(points.X))
		for i := range points.X {
			predicted[i] = analysis.EvaluatePolynomial(coeffs, points.X[i])
		}
		mse := analysis.MSE(points.Y, predicted)

		if minMSE < 0 || mse < minMSE {
			minMSE = mse
			bestDegree = degree
		}

		degrees = append(degrees, float64(degree))
		errs = append(errs, mse)
	}

	chartCopy := charting.CopyChart(SampleMSEChart)
	chartCopy.UpdatePointsForDataset(OriginalDataID, degrees, errs)

	gvars := chartCopy.Datasets[OriginalDataID].GetFields()
	for i := range gvars {
		if gvars[i].ID == BestDegreeID {
			gvars[i].Label = fmt.Sprintf("Best Degree: %d (MSE: %.4e)", bestDegree, minMSE)
		}
	}

	res = charting.NewRenderResponse()
	res.AddChart(SampleMSEID, &chartCopy)
	return res
}
