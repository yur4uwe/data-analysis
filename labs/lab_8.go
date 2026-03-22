package labs

import (
	"labs/charting"
	"labs/labs/holt"
)

type Lab8Provider struct{}

var _ charting.LabProvider = Lab8Provider{}

func NewLab8() *Lab8Provider {
	return &Lab8Provider{}
}

func (lp Lab8Provider) GetMetadata() charting.LabMetadata {
	return holt.Metadata
}

func (lp Lab8Provider) GetConfig() charting.LabConfig {
	return holt.Config
}

func (lp Lab8Provider) Render(req *charting.RenderRequest) *charting.RenderResponse {
	res := &charting.RenderResponse{}
	if req == nil {
		return res.NewError("request is nil")
	}

	chart, ok := holt.Config.Charts[req.ChartID]
	if !ok {
		return res.NewErrorf("chart with id %q not found", req.ChartID)
	}

	return chart.RenderFunc(req)
}
