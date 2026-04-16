package charting

import (
	"fmt"
)

// MutableField defines an input field that affects chart rendering
type MutableField struct {
	ID      string       `json:"id"`
	Label   string       `json:"label"`
	Default float64      `json:"default"`
	Min     float64      `json:"min"`
	Max     float64      `json:"max"`
	Step    float64      `json:"step"`
	Control FieldControl `json:"control"`
	Hint    string       `json:"hint,omitempty"`
	Options []string     `json:"options,omitempty"` // For select control
}

// Chart represents one visualization
type Chart struct {
	ID             string                               `json:"id"`
	Title          string                               `json:"title"`
	Type           GraphType                            `json:"type"` // "line", "bar", "scatter", "bubble", "pie"
	XAxisLabel     string                               `json:"xAxisLabel"`
	YAxisLabel     string                               `json:"yAxisLabel"`
	XAxisConfig    AxisConfig                           `json:"xAxisConfig,omitempty"` // X-axis scale configuration
	YAxisConfig    AxisConfig                           `json:"yAxisConfig,omitempty"` // Y-axis scale configuration
	Datasets       map[string]Dataset                   `json:"datasets"`
	Labels         []string                             `json:"labels,omitempty"` // For indexed data (line/bar charts)
	ChartVariables []MutableField                       `json:"chartVariables"`   // Fields affecting this chart only
	Misc           string                               `json:"misc,omitempty"`   // Extra data/metadata
	RenderFunc     func(*RenderRequest) *RenderResponse `json:"-"`
}

func (c *Chart) Meta() ChartMetadata {
	meta := ChartMetadata{
		ID:             c.ID,
		Title:          c.Title,
		ChartVariables: c.ChartVariables,
	}

	graphVars := map[string][]MutableField{}
	for graphId, graph := range c.Datasets {
		vars := graph.GetFields()
		if len(vars) == 0 {
			graphVars[graphId] = []MutableField{}
		}
		graphVars[graphId] = vars
	}

	meta.GraphVariables = graphVars

	return meta
}

func (c *Chart) UpdateLabel(labelId string, newLabel string) {
	found := false
	for i := range c.ChartVariables {
		if c.ChartVariables[i].ID == labelId {
			c.ChartVariables[i].Label = newLabel
			found = true
			break
		}
	}
	// If not found, panic as it's a programmer error
	if !found {
		panic(fmt.Errorf("label %q not found", labelId))
	}
}

func (c *Chart) UpdatePointsForDataset(datasetId string, x, y []float64) {
	if _, ok := c.Datasets[datasetId]; !ok {
		panic(fmt.Errorf("dataset not found in chart"))
	}

	if len(x) != len(y) {
		panic(fmt.Errorf("x and y are of different length: %d and %d", len(x), len(y)))
	}

	p := make([]any, len(x))
	for i := range x {
		val := y[i]
		p[i] = &DataPoint{
			X: x[i],
			Y: &val,
		}
	}

	c.Datasets[datasetId].UpdateData(p)
}

func (c *Chart) UpdateDataPointsForDataset(datasetId string, points []DataPoint) {
	dataset, ok := c.Datasets[datasetId]
	if !ok {
		panic(fmt.Errorf("dataset not found in chart"))
	}

	p := make([]any, len(points))

	for i := range points {
		p[i] = points[i]
	}

	dataset.UpdateData(p)
}

func (c *Chart) UpdateDataForDataset(datasetId string, data []any) {
	dataset, ok := c.Datasets[datasetId]
	if !ok {
		panic(fmt.Errorf("dataset not found in chart"))
	}
	dataset.UpdateData(data)
}

// GenerateLabels derives chart Labels from the x-values of PointData across all datasets.
// It picks the dataset with the most points and formats each x-value with the given decimal precision.
// This should be called after all datasets have been populated with point data.
func (c *Chart) GenerateLabels(precision int) {
	var best []any
	for _, ds := range c.Datasets {
		p := ds.GetData()
		if len(p) > len(best) {
			// Check if it's a dataset with points
			if _, ok := ds.(*GridDataset); ok {
				best = p
			} else if _, ok := ds.(*HeatmapDataset); ok {
				best = p
			}
		}
	}
	if len(best) == 0 {
		return
	}
	points := AnyToPoints(best)
	labels := make([]string, len(points))
	format := fmt.Sprintf("%%.%df", precision)
	for i, p := range points {
		if p.Y != nil {
			labels[i] = fmt.Sprintf(format, p.X)
		} else {
			labels[i] = ""
		}
	}
	c.Labels = labels
}

type ChartMetadata struct {
	ID             string                    `json:"id"`
	Title          string                    `json:"title"`
	ChartVariables []MutableField            `json:"chartVariables"` // Fields affecting this chart only
	GraphVariables map[string][]MutableField `json:"graphVariables"` // Fields affecting charts graphs
}

// CopyChart creates a deep copy of a chart template without data
// Preserves all configuration but resets the Data fields to be populated
func CopyChart(original Chart) Chart {
	// Start with a value copy of the struct (copies non-reference types)
	newChart := original

	// Deep copy Labels slice
	if original.Labels != nil {
		newChart.Labels = make([]string, len(original.Labels))
		copy(newChart.Labels, original.Labels)
	}

	// Deep copy ChartVariables slice
	if original.ChartVariables != nil {
		newChart.ChartVariables = make([]MutableField, len(original.ChartVariables))
		copy(newChart.ChartVariables, original.ChartVariables)
	}

	// Deep copy Datasets map
	newChart.Datasets = make(map[string]Dataset, len(original.Datasets))
	for key, dataset := range original.Datasets {
		if dataset == nil {
			continue
		}

		// Store the pointer to our new struct copy
		newChart.Datasets[key] = dataset.Copy()
	}

	return newChart
}
