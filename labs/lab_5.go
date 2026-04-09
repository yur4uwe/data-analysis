package labs

import (
	"labs/charting"
	statslab "labs/labs/stats"
)

type Lab5Provider struct{}

var _ charting.LabProvider = Lab5Provider{}

func NewLab5() *Lab5Provider {
	return &Lab5Provider{}
}

func (lp Lab5Provider) GetMetadata() charting.LabMetadata {
	return statslab.Metadata
}

func (lp Lab5Provider) GetConfig() charting.LabConfig {
	return statslab.Config
}

func (lp Lab5Provider) Render(req *charting.RenderRequest) *charting.RenderResponse {
	res := &charting.RenderResponse{}
	if req == nil {
		return res.NewError("request is nil")
	}

	chart, ok := statslab.Config.Charts[req.ChartID]
	if !ok {
		return res.NewErrorf("chart with id %s not found", req.ChartID)
	}

	return chart.RenderFunc(req)
}
