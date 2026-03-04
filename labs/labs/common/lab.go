package common

// LabMetadata contains information about a single lab
type LabMetadata struct {
	ID     string
	Name   string
	Charts map[string]ChartMetadata
}

// GetLabsResponse returns all available labs for the UI
type GetLabsResponse struct {
	Labs []LabMetadata `json:"labs"`
}

// LabConfig is the complete configuration for a lab
type LabConfig struct {
	Lab    LabMetadata       `json:"lab"`
	Charts map[string]*Chart `json:"charts"`
}

func NewLabConfig(labID, labName string, charts map[string]*Chart) LabConfig {
	chartsMeta := make(map[string]ChartMetadata, len(charts))
	for id, chart := range charts {
		chartsMeta[id] = chart.Meta()
	}

	return LabConfig{
		Lab: LabMetadata{
			ID:     labID,
			Name:   labName,
			Charts: chartsMeta,
		},
		Charts: charts,
	}
}
