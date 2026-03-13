package charting

import (
	"errors"
	"fmt"
)

// ==================== CHART CONFIGURATION ====================

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

// DataPoint represents a point with x and y coordinates for scatter/bubble charts
type DataPoint struct {
	X float64 `json:"x" csv:"x"`
	Y float64 `json:"y" csv:"y"`
}

// ChartDataset represents a single line/bar/scatter set
type ChartDataset struct {
	Label           string         `json:"label"`
	Data            []float64      `json:"data,omitempty"`      // For indexed data (line/bar with labels)
	PointData       []DataPoint    `json:"pointData,omitempty"` // For scatter/bubble charts with {x,y}
	GraphVariables  []MutableField `json:"fields,omitempty"`
	BorderColor     string         `json:"borderColor"`
	BackgroundColor []string       `json:"backgroundColor,omitempty"`
	Tension         float64        `json:"tension,omitempty"`
	Fill            bool           `json:"fill,omitempty"`
	Hidden          bool           `json:"hidden,omitempty"`
	PointRadius     int            `json:"pointRadius,omitempty"`
	BorderWidth     int            `json:"borderWidth,omitempty"`
	ShowLine        bool           `json:"showLine,omitempty"`
	Togglable       bool           `json:"togglable,omitempty"` // Can user toggle visibility?
	PointStyle      string         `json:"pointStyle,omitempty"`
	PointLabels     []string       `json:"pointLabels,omitempty"` // Per-point labels shown via datalabels plugin
}

func (cd *ChartDataset) UpdatePoints(x, y []float64) error {
	if len(x) == 0 || len(y) == 0 {
		return errors.New("empty x or y arrays")
	}
	if len(x) != len(y) {
		return errors.New("arrays of different length")
	}

	cd.PointData = make([]DataPoint, len(x))
	for i := range x {
		cd.PointData[i] = DataPoint{X: x[i], Y: y[i]}
	}

	return nil
}

// AxisConfig defines how an axis should be displayed

// Chart represents one visualization
type Chart struct {
	ID             string                               `json:"id"`
	Title          string                               `json:"title"`
	Type           ChartType                            `json:"type"` // "line", "bar", "scatter", "bubble", "pie"
	XAxisLabel     string                               `json:"xAxisLabel"`
	YAxisLabel     string                               `json:"yAxisLabel"`
	XAxisConfig    AxisConfig                           `json:"xAxisConfig,omitempty"` // X-axis scale configuration
	YAxisConfig    AxisConfig                           `json:"yAxisConfig,omitempty"` // Y-axis scale configuration
	Datasets       map[string]*ChartDataset             `json:"datasets"`
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
		if len(graph.GraphVariables) == 0 {
			graphVars[graphId] = []MutableField{}
		}
		graphVars[graphId] = graph.GraphVariables
	}

	meta.GraphVariables = graphVars

	return meta
}

func (c *Chart) UpdatePointsForDataset(datasetId string, x, y []float64) error {
	if _, ok := c.Datasets[datasetId]; !ok {
		return errors.New("dataset not found in chart")
	}
	return c.Datasets[datasetId].UpdatePoints(x, y)
}

func (c *Chart) UpdateDataPointsForDataset(datasetId string, points []DataPoint) error {
	dataset, ok := c.Datasets[datasetId]
	if !ok {
		return errors.New("dataset not found in chart")
	}
	dataset.PointData = points
	c.Datasets[datasetId] = dataset
	return nil
}

func (c *Chart) UpdateDataForDataset(datasetId string, data []float64) error {
	if _, ok := c.Datasets[datasetId]; !ok {
		return errors.New("dataset not found in chart")
	}
	c.Datasets[datasetId].Data = data
	return nil
}

// GenerateLabels derives chart Labels from the x-values of PointData across all datasets.
// It picks the dataset with the most points and formats each x-value with the given decimal precision.
// This should be called after all datasets have been populated with point data.
func (c *Chart) GenerateLabels(precision int) {
	var best []DataPoint
	for _, ds := range c.Datasets {
		if len(ds.PointData) > len(best) {
			best = ds.PointData
		}
	}
	if len(best) == 0 {
		return
	}
	labels := make([]string, len(best))
	format := fmt.Sprintf("%%.%df", precision)
	for i, p := range best {
		labels[i] = fmt.Sprintf(format, p.X)
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
	newChart.Datasets = make(map[string]*ChartDataset, len(original.Datasets))
	for key, dataset := range original.Datasets {
		if dataset == nil {
			continue
		}

		// Dereference to create a new struct with the same config values
		dsCopy := *dataset

		// Reset data slices for the copy
		dsCopy.Data = nil
		dsCopy.PointData = nil

		// Deep copy GraphVariables in the dataset
		if dataset.GraphVariables != nil {
			dsCopy.GraphVariables = make([]MutableField, len(dataset.GraphVariables))
			copy(dsCopy.GraphVariables, dataset.GraphVariables)
		}

		// Store the pointer to our new struct copy
		newChart.Datasets[key] = &dsCopy
	}

	return newChart
}
